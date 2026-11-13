package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	_ "github.com/dutra/recodationsApp/docs"
	auditCtrl "github.com/dutra/recodationsApp/internal/modules/audit/controller"
	auditRepo "github.com/dutra/recodationsApp/internal/modules/audit/repository"
	authCtrl "github.com/dutra/recodationsApp/internal/modules/auth/controller"
	authUC "github.com/dutra/recodationsApp/internal/modules/auth/usecase"
	inboxDomain "github.com/dutra/recodationsApp/internal/modules/inbox/domain"
	inboxRepo "github.com/dutra/recodationsApp/internal/modules/inbox/repository"
	inboxUC "github.com/dutra/recodationsApp/internal/modules/inbox/usecase"
	messagesCtrl "github.com/dutra/recodationsApp/internal/modules/messages/controller"
	messagesRepo "github.com/dutra/recodationsApp/internal/modules/messages/repository"
	messagesUC "github.com/dutra/recodationsApp/internal/modules/messages/usecase"
	rulesCtrl "github.com/dutra/recodationsApp/internal/modules/rules/controller"
	rulesRepo "github.com/dutra/recodationsApp/internal/modules/rules/repository"
	rulesUC "github.com/dutra/recodationsApp/internal/modules/rules/usecase"
	usersCtrl "github.com/dutra/recodationsApp/internal/modules/users/controller"
	usersRepo "github.com/dutra/recodationsApp/internal/modules/users/repository"
	usersUC "github.com/dutra/recodationsApp/internal/modules/users/usecase"
	"github.com/dutra/recodationsApp/internal/shared/infrastructure/cache"
	"github.com/dutra/recodationsApp/internal/shared/infrastructure/database"
	"github.com/dutra/recodationsApp/internal/shared/infrastructure/logger"
	"github.com/dutra/recodationsApp/internal/shared/infrastructure/middleware"
	"github.com/dutra/recodationsApp/internal/shared/infrastructure/storage"
	"github.com/dutra/recodationsApp/pkg/snowflake"
	"github.com/gofiber/fiber/v3"
	"github.com/redis/go-redis/v9"
)

// @title Recodations Modular API
// @version 1.0
// @description Modular Monolith API in Go using Fiber v3, SOLID, SnowflakeID, Optimistic Locking, State Transition Pattern, Inbox Pattern, Two-Level L1/L2 Cache, Resilient Google Cloud Storage (floci-gcp), RBAC with JWT, and Audit Logging.
// @contact.name API Team
// @contact.email support@recodations.com
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
// @host localhost:3000
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Enter Bearer token as: Bearer <token>
func main() {
	env := getEnv("APP_ENV", "development")
	logger.InitLogger(env)

	// 1. Initialize Snowflake ID Generator
	nodeID, _ := strconv.ParseInt(getEnv("SNOWFLAKE_NODE_ID", "1"), 10, 64)
	if err := snowflake.Init(nodeID); err != nil {
		slog.Error("Failed to initialize Snowflake generator", "error", err)
		os.Exit(1)
	}
	slog.Info("Snowflake generator initialized", "nodeId", nodeID)

	// 2. Connect to PostgreSQL
	dbCfg := database.Config{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASS", "postgres"),
		DBName:   getEnv("DB_NAME", "recodations"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}

	db, err := database.Connect(dbCfg)
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}
	slog.Info("Connected to PostgreSQL successfully")

	// 3. Initialize Two-Level Cache (L1 Memory + L2 Redis)
	l1Cache := cache.NewMemoryCache()
	var l2Cache cache.Cache
	redisAddr := getEnv("REDIS_ADDR", "localhost:6379")
	redisPass := getEnv("REDIS_PASSWORD", "")
	rClient := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPass,
		DB:       0,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := rClient.Ping(ctx).Err(); err == nil {
		l2Cache = cache.NewRedisCache(rClient)
		slog.Info("Connected to Redis L2 cache successfully", "addr", redisAddr)
	} else {
		slog.Warn("Redis unavailable, proceeding with L1 Memory Cache only", "addr", redisAddr, "error", err)
	}
	tieredCache := cache.NewTieredCache(l1Cache, l2Cache)

	// 4. Initialize Resilient Google Cloud Storage
	gcpProjectID := getEnv("GOOGLE_CLOUD_PROJECT", "floci-local")
	bucketName := getEnv("GCS_BUCKET_NAME", "photos")
	gcs, err := storage.NewGCSStorage(context.Background(), gcpProjectID)
	var resilientStorage storage.StorageService
	if err != nil {
		slog.Warn("Failed initializing GCS client, storage disabled", "error", err)
	} else {
		resilientStorage = storage.NewResilientStorage(gcs)
		slog.Info("Google Cloud Storage client initialized with Circuit Breaker and Retry", "bucket", bucketName)
	}

	// 5. Initialize Repositories
	ruleRepository := rulesRepo.NewRuleRepository(db)
	userRepository := usersRepo.NewUserRepository(db)
	messageRepository := messagesRepo.NewMessageRepository(db)
	inboxRepository := inboxRepo.NewInboxRepository(db)
	auditRepository := auditRepo.NewAuditRepository(db)

	// 6. Initialize Inbox Service & Background Worker
	inboxService := inboxUC.NewInboxService(inboxRepository)
	go startInboxWorker(context.Background(), inboxService)

	// 7. Initialize Use Cases
	jwtSecret := getEnv("JWT_SECRET", "super-secret-jwt-token-key-32-bytes!")
	loginUC := authUC.NewLoginUseCase(userRepository, jwtSecret, 24*time.Hour)

	createRuleUC := rulesUC.NewCreateRuleUseCase(ruleRepository)
	listRulesUC := rulesUC.NewListRulesUseCase(ruleRepository)

	createUserUC := usersUC.NewCreateUserUseCase(userRepository)
	getUserUC := usersUC.NewGetUserUseCase(userRepository, tieredCache)
	listUsersUC := usersUC.NewListUsersUseCase(userRepository)
	assignRulesUC := usersUC.NewAssignRulesUseCase(userRepository, tieredCache)

	createMsgUC := messagesUC.NewCreateMessageUseCase(messageRepository, resilientStorage, inboxService, tieredCache, bucketName)
	getMsgUC := messagesUC.NewGetMessageUseCase(messageRepository, tieredCache)
	listMsgUC := messagesUC.NewListMessagesUseCase(messageRepository)
	updateMsgUC := messagesUC.NewUpdateMessageUseCase(messageRepository, resilientStorage, tieredCache, bucketName)
	deleteMsgUC := messagesUC.NewDeleteMessageUseCase(messageRepository, resilientStorage, tieredCache, bucketName)
	transitionMsgUC := messagesUC.NewTransitionMessageUseCase(messageRepository, tieredCache)

	// 8. Initialize Controllers
	authController := authCtrl.NewAuthController(loginUC, createUserUC)
	ruleController := rulesCtrl.NewRuleController(createRuleUC, listRulesUC)
	userController := usersCtrl.NewUserController(createUserUC, getUserUC, listUsersUC, assignRulesUC)
	messageController := messagesCtrl.NewMessageController(createMsgUC, getMsgUC, listMsgUC, updateMsgUC, deleteMsgUC, transitionMsgUC)
	auditController := auditCtrl.NewAuditController(auditRepository)

	// 9. Build Fiber v3 Application
	app := fiber.New(fiber.Config{
		AppName: "Recodations Modular API",
	})

	// Middlewares
	app.Use(middleware.RecoverMiddleware())
	app.Use(middleware.TraceMiddleware())
	app.Use(middleware.RateLimitMiddleware(100, 1*time.Minute))
	app.Use(middleware.AuditMiddleware(auditRepository))

	// Health Checks
	app.Get("/healthz", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok", "time": time.Now().UTC()})
	})
	app.Get("/livez", func(c fiber.Ctx) error {
		return c.SendString("alive")
	})
	app.Get("/readyz", func(c fiber.Ctx) error {
		sqlDB, err := db.DB()
		if err != nil || sqlDB.Ping() != nil {
			return c.Status(fiber.StatusServiceUnavailable).SendString("database not ready")
		}
		return c.SendString("ready")
	})

	// Documentation Routes (Swagger UI + Swagger JSON + Scalar Modern Interactive UI)
	app.Get("/swagger/doc.json", func(c fiber.Ctx) error {
		c.Set("Content-Type", "application/json")
		return c.SendFile("./docs/swagger.json")
	})
	app.Get("/swagger", func(c fiber.Ctx) error {
		c.Set("Content-Type", "text/html")
		return c.SendString(GetSwaggerUIHTML())
	})
	app.Get("/swagger/*", func(c fiber.Ctx) error {
		c.Set("Content-Type", "text/html")
		return c.SendString(GetSwaggerUIHTML())
	})
	app.Get("/docs", func(c fiber.Ctx) error {
		c.Set("Content-Type", "text/html")
		return c.SendString(`<!doctype html>
<html>
  <head>
    <title>Recodations API Reference</title>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
  </head>
  <body>
    <script id="api-reference" data-url="/swagger/doc.json"></script>
    <script src="https://cdn.jsdelivr.net/npm/@scalar/api-reference"></script>
  </body>
</html>`)
	})

	// API Routes
	api := app.Group("/api/v1")

	// Public Auth endpoints
	authGroup := api.Group("/auth")
	authGroup.Post("/login", authController.Login)
	authGroup.Post("/register", authController.Register)

	// Protected Routes
	protected := api.Group("", middleware.JWTMiddleware(jwtSecret))

	// Rules routes
	rulesGroup := protected.Group("/rules")
	rulesGroup.Post("", middleware.RequireRule("users:manage"), ruleController.CreateRule)
	rulesGroup.Get("", ruleController.ListRules)

	// Users routes
	usersGroup := protected.Group("/users")
	usersGroup.Post("", middleware.RequireRule("users:manage"), userController.CreateUser)
	usersGroup.Get("", middleware.RequireRule("users:manage"), userController.ListUsers)
	usersGroup.Get("/:id", userController.GetUser)
	usersGroup.Put("/:id/rules", middleware.RequireRule("users:manage"), userController.AssignRules)

	// Messages routes
	messagesGroup := protected.Group("/messages")
	messagesGroup.Post("", middleware.RequireRule("messages:write"), messageController.CreateMessage)
	messagesGroup.Get("", middleware.RequireRule("messages:read"), messageController.ListMessages)
	messagesGroup.Get("/:id", middleware.RequireRule("messages:read"), messageController.GetMessage)
	messagesGroup.Put("/:id", middleware.RequireRule("messages:write"), messageController.UpdateMessage)
	messagesGroup.Patch("/:id/transition", middleware.RequireRule("messages:write"), messageController.TransitionMessage)
	messagesGroup.Delete("/:id", middleware.RequireRule("messages:delete"), messageController.DeleteMessage)

	// Audit Logs routes
	auditGroup := protected.Group("/audit-logs")
	auditGroup.Get("", middleware.RequireRule("audit:read"), auditController.ListAuditLogs)

	// 10. Start Server with Graceful Shutdown
	port := getEnv("PORT", "3000")
	go func() {
		slog.Info("Starting Recodations API", "port", port)
		if err := app.Listen(":" + port); err != nil {
			slog.Info("Server stopped listening", "error", err)
		}
	}()

	// Wait for termination signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	slog.Info("Shutting down server gracefully...")
	_ = app.Shutdown()
	slog.Info("Server shutdown complete.")
}

func startInboxWorker(ctx context.Context, inboxSvc *inboxUC.InboxService) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			_, _ = inboxSvc.ProcessPendingEvents(ctx, 10, func(c context.Context, event inboxDomain.InboxMessage) error {
				slog.Info("Asynchronously handled inbox event", "eventId", event.EventID, "type", event.EventType)
				return nil
			})
		}
	}
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}
