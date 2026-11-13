package controller

import (
	"github.com/dutra/recodationsApp/internal/modules/rules/dto"
	"github.com/dutra/recodationsApp/internal/modules/rules/usecase"
	"github.com/dutra/recodationsApp/internal/shared/presenter"
	"github.com/gofiber/fiber/v3"
)

type RuleController struct {
	createRuleUC *usecase.CreateRuleUseCase
	listRulesUC  *usecase.ListRulesUseCase
}

func NewRuleController(createRuleUC *usecase.CreateRuleUseCase, listRulesUC *usecase.ListRulesUseCase) *RuleController {
	return &RuleController{
		createRuleUC: createRuleUC,
		listRulesUC:  listRulesUC,
	}
}

// CreateRule godoc
// @Summary Create a new RBAC rule
// @Description Creates a new granular permission rule
// @Tags Rules
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param input body dto.CreateRuleInput true "Rule Input"
// @Success 201 {object} presenter.ResponseDoc
// @Failure 400 {object} presenter.ResponseDoc
// @Failure 409 {object} presenter.ResponseDoc
// @Router /api/v1/rules [post]
func (ctrl *RuleController) CreateRule(c fiber.Ctx) error {
	var input dto.CreateRuleInput
	if err := c.Bind().JSON(&input); err != nil {
		return presenter.RenderValidationError(c, "Invalid request payload", map[string][]string{
			"body": {"failed parsing JSON payload"},
		})
	}

	result := ctrl.createRuleUC.Execute(c.Context(), input)
	if !result.Success {
		return presenter.RenderResponse(c, presenter.MapErrorToStatusCode(result.Error), result)
	}

	return presenter.RenderResponse(c, fiber.StatusCreated, result)
}

// ListRules godoc
// @Summary List all registered rules
// @Description Retrieves all available permissions in the system
// @Tags Rules
// @Produce json
// @Security BearerAuth
// @Success 200 {object} presenter.ResponseDoc
// @Failure 500 {object} presenter.ResponseDoc
// @Router /api/v1/rules [get]
func (ctrl *RuleController) ListRules(c fiber.Ctx) error {
	result := ctrl.listRulesUC.Execute(c.Context())
	if !result.Success {
		return presenter.RenderResponse(c, presenter.MapErrorToStatusCode(result.Error), result)
	}

	return presenter.RenderResponse(c, fiber.StatusOK, result)
}
