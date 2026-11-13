package controller

import (
	authDto "github.com/dutra/recodationsApp/internal/modules/auth/dto"
	authUC "github.com/dutra/recodationsApp/internal/modules/auth/usecase"
	usersDto "github.com/dutra/recodationsApp/internal/modules/users/dto"
	usersUC "github.com/dutra/recodationsApp/internal/modules/users/usecase"
	"github.com/dutra/recodationsApp/internal/shared/presenter"
	"github.com/gofiber/fiber/v3"
)

type AuthController struct {
	loginUseCase    *authUC.LoginUseCase
	registerUseCase *usersUC.CreateUserUseCase
}

func NewAuthController(loginUseCase *authUC.LoginUseCase, registerUseCase *usersUC.CreateUserUseCase) *AuthController {
	return &AuthController{
		loginUseCase:    loginUseCase,
		registerUseCase: registerUseCase,
	}
}

// Login godoc
// @Summary User login
// @Description Authenticates user and returns JWT bearer token
// @Tags Auth
// @Accept json
// @Produce json
// @Param credentials body authDto.LoginInput true "Credentials"
// @Success 200 {object} presenter.ResponseDoc
// @Failure 401 {object} presenter.ResponseDoc
// @Router /api/v1/auth/login [post]
func (ctrl *AuthController) Login(c fiber.Ctx) error {
	var input authDto.LoginInput
	if err := c.Bind().JSON(&input); err != nil {
		return presenter.RenderValidationError(c, "Invalid JSON payload", map[string][]string{
			"body": {"failed parsing input"},
		})
	}

	result := ctrl.loginUseCase.Execute(c.Context(), input)
	if !result.Success {
		return presenter.RenderResponse(c, presenter.MapErrorToStatusCode(result.Error), result)
	}

	return presenter.RenderResponse(c, fiber.StatusOK, result)
}

// Register godoc
// @Summary Public user registration
// @Description Registers a new standard user
// @Tags Auth
// @Accept json
// @Produce json
// @Param input body usersDto.CreateUserInput true "User Registration Data"
// @Success 201 {object} presenter.ResponseDoc
// @Failure 400 {object} presenter.ResponseDoc
// @Failure 409 {object} presenter.ResponseDoc
// @Router /api/v1/auth/register [post]
func (ctrl *AuthController) Register(c fiber.Ctx) error {
	var input usersDto.CreateUserInput
	if err := c.Bind().JSON(&input); err != nil {
		return presenter.RenderValidationError(c, "Invalid JSON payload", map[string][]string{
			"body": {"failed parsing input"},
		})
	}

	result := ctrl.registerUseCase.Execute(c.Context(), input)
	if !result.Success {
		return presenter.RenderResponse(c, presenter.MapErrorToStatusCode(result.Error), result)
	}

	return presenter.RenderResponse(c, fiber.StatusCreated, result)
}
