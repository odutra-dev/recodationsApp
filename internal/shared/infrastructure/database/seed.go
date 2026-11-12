package database

import (
	"context"
	"log/slog"
	"time"

	rulesDomain "github.com/dutra/recodationsApp/internal/modules/rules/domain"
	usersDomain "github.com/dutra/recodationsApp/internal/modules/users/domain"
	"github.com/dutra/recodationsApp/pkg/snowflake"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// SeedInitialData populates default rules and a root administrator user if they do not exist.
func SeedInitialData(ctx context.Context, db *gorm.DB) error {
	defaultRules := []rulesDomain.Rule{
		{
			Name:        "Read Messages",
			Code:        "messages:read",
			Description: "Permission to view messages",
		},
		{
			Name:        "Write Messages",
			Code:        "messages:write",
			Description: "Permission to create and update messages",
		},
		{
			Name:        "Delete Messages",
			Code:        "messages:delete",
			Description: "Permission to delete messages",
		},
		{
			Name:        "Manage Users",
			Code:        "users:manage",
			Description: "Permission to manage users and assign rules",
		},
		{
			Name:        "Read Audit Logs",
			Code:        "audit:read",
			Description: "Permission to view system audit logs",
		},
	}

	var seededRules []rulesDomain.Rule
	for _, r := range defaultRules {
		var existing rulesDomain.Rule
		err := db.WithContext(ctx).Where("code = ?", r.Code).First(&existing).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				r.ID = snowflake.Generate()
				r.CreatedAt = time.Now().UTC()
				r.UpdatedAt = time.Now().UTC()
				if err := db.WithContext(ctx).Create(&r).Error; err != nil {
					return err
				}
				slog.Info("Seeded rule", "code", r.Code)
				seededRules = append(seededRules, r)
			} else {
				return err
			}
		} else {
			seededRules = append(seededRules, existing)
		}
	}

	// Seed Admin User
	adminEmail := "admin@recodations.com"
	var existingAdmin usersDomain.User
	err := db.WithContext(ctx).Where("email = ?", adminEmail).First(&existingAdmin).Error
	if err == gorm.ErrRecordNotFound {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("Admin@123"), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		adminUser := usersDomain.User{
			ID:           snowflake.Generate(),
			Name:         "System Administrator",
			Email:        adminEmail,
			PasswordHash: string(hashedPassword),
			Version:      1,
			CreatedAt:    time.Now().UTC(),
			UpdatedAt:    time.Now().UTC(),
		}

		if err := db.WithContext(ctx).Create(&adminUser).Error; err != nil {
			return err
		}

		// Assign all rules to admin user
		for _, r := range seededRules {
			userRule := rulesDomain.UserRule{
				UserID: adminUser.ID,
				RuleID: r.ID,
			}
			if err := db.WithContext(ctx).Create(&userRule).Error; err != nil {
				slog.Error("Failed to assign rule to admin", "error", err, "rule", r.Code)
			}
		}

		slog.Info("Seeded default admin user", "email", adminEmail)
	}

	return nil
}
