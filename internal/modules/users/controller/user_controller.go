package controller

import (
	"strconv"

	"github.com/dutra/recodationsApp/internal/modules/users/dto"
	"github.com/dutra/recodationsApp/internal/modules/users/usecase"
	"github.com/dutra/recodationsApp/internal/shared/domain"
	"github.com/dutra/recodationsApp/internal/shared/presenter"
	"github.com/gofiber/fiber/v3"
)

type UserController struct {
	createUserUC  *usecase.CreateUserUseCase
	getUserUC     *usecase.GetUserUseCase
	listUsersUC   *usecase.ListUsersUseCase
	assignRulesUC *usecase.AssignRulesUseCase
}

func NewUserController(
	createUserUC *usecase.CreateUserUseCase,
	getUserUC *usecase.GetUserUseCase,
	listUsersUC *usecase.ListUsersUseCase,
	assignRulesUC *usecase.AssignRulesUseCase,
) *UserController {
	return &UserController{
		createUserUC:  createUserUC,
		getUserUC:     getUserUC,
		listUsersUC:   listUsersUC,
		assignRulesUC: assignRulesUC,
	}
}

// CreateUser godoc
// @Summary Create a new user
// @Description Registers a new user with optional rule associations
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param input body dto.CreateUserInput true "User Creation Data"
// @Success 201 {object} presenter.ResponseDoc
// @Failure 400 {object} presenter.ResponseDoc
// @Failure 409 {object} presenter.ResponseDoc
// @Router /api/v1/users [post]
func (ctrl *UserController) CreateUser(c fiber.Ctx) error {
	var input dto.CreateUserInput
	if err := c.Bind().JSON(&input); err != nil {
		return presenter.RenderValidationError(c, "Invalid JSON payload", map[string][]string{
			"body": {"failed parsing input"},
		})
	}

	result := ctrl.createUserUC.Execute(c.Context(), input)
	if !result.Success {
		return presenter.RenderResponse(c, presenter.MapErrorToStatusCode(result.Error), result)
	}

	return presenter.RenderResponse(c, fiber.StatusCreated, result)
}

// GetUser godoc
// @Summary Get user by ID
// @Description Fetches user profile with rules (cached via L1/L2 cache)
// @Tags Users
// @Produce json
// @Security BearerAuth
// @Param id path int true "User Snowflake ID"
// @Success 200 {object} presenter.ResponseDoc
// @Failure 404 {object} presenter.ResponseDoc
// @Router /api/v1/users/{id} [get]
func (ctrl *UserController) GetUser(c fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return presenter.RenderError(c, domain.ErrValidation, "Invalid user ID format")
	}

	result := ctrl.getUserUC.Execute(c.Context(), id)
	if !result.Success {
		return presenter.RenderResponse(c, presenter.MapErrorToStatusCode(result.Error), result)
	}

	return presenter.RenderResponse(c, fiber.StatusOK, result)
}

// ListUsers godoc
// @Summary List all users
// @Description Retrieves paginated list of users
// @Tags Users
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Page size" default(20)
// @Success 200 {object} presenter.ResponseDoc
// @Failure 500 {object} presenter.ResponseDoc
// @Router /api/v1/users [get]
func (ctrl *UserController) ListUsers(c fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))

	result := ctrl.listUsersUC.Execute(c.Context(), page, limit)
	if !result.Success {
		return presenter.RenderResponse(c, presenter.MapErrorToStatusCode(result.Error), result)
	}

	return presenter.RenderResponse(c, fiber.StatusOK, result)
}

// AssignRules godoc
// @Summary Assign rules to user
// @Description Updates user permissions and invalidates L1/L2 cache
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "User Snowflake ID"
// @Param input body dto.AssignRulesInput true "Rule IDs"
// @Success 200 {object} presenter.ResponseDoc
// @Failure 404 {object} presenter.ResponseDoc
// @Router /api/v1/users/{id}/rules [put]
func (ctrl *UserController) AssignRules(c fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return presenter.RenderError(c, domain.ErrValidation, "Invalid user ID format")
	}

	var input dto.AssignRulesInput
	if err := c.Bind().JSON(&input); err != nil {
		return presenter.RenderValidationError(c, "Invalid JSON payload", map[string][]string{
			"body": {"failed parsing input"},
		})
	}

	result := ctrl.assignRulesUC.Execute(c.Context(), id, input)
	if !result.Success {
		return presenter.RenderResponse(c, presenter.MapErrorToStatusCode(result.Error), result)
	}

	return presenter.RenderResponse(c, fiber.StatusOK, result)
}
