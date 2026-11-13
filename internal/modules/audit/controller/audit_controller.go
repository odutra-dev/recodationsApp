package controller

import (
	"strconv"

	"github.com/dutra/recodationsApp/internal/modules/audit/domain"
	"github.com/dutra/recodationsApp/internal/modules/audit/repository"
	sharedDomain "github.com/dutra/recodationsApp/internal/shared/domain"
	"github.com/dutra/recodationsApp/internal/shared/presenter"
	"github.com/gofiber/fiber/v3"
)

type AuditController struct {
	repo repository.AuditRepository
}

func NewAuditController(repo repository.AuditRepository) *AuditController {
	return &AuditController{repo: repo}
}

// ListAuditLogs godoc
// @Summary List system audit logs
// @Description Retrieves system audit logs for administrative inspection
// @Tags Audit
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Page size" default(20)
// @Success 200 {object} presenter.ResponseDoc
// @Failure 403 {object} presenter.ResponseDoc
// @Router /api/v1/audit-logs [get]
func (ctrl *AuditController) ListAuditLogs(c fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))

	logs, _, err := ctrl.repo.List(c.Context(), page, limit)
	if err != nil {
		return presenter.RenderResponse(c, fiber.StatusInternalServerError, sharedDomain.Fail[[]domain.AuditLog](err, "failed fetching audit logs"))
	}

	return presenter.RenderSuccess(c, fiber.StatusOK, logs, "Audit logs retrieved successfully")
}
