package controller

import (
	"io"
	"strconv"
	"strings"

	"github.com/dutra/recodationsApp/internal/modules/messages/dto"
	"github.com/dutra/recodationsApp/internal/modules/messages/usecase"
	"github.com/dutra/recodationsApp/internal/shared/domain"
	"github.com/dutra/recodationsApp/internal/shared/presenter"
	"github.com/gofiber/fiber/v3"
)

type MessageController struct {
	createMsgUC     *usecase.CreateMessageUseCase
	getMsgUC        *usecase.GetMessageUseCase
	listMsgUC       *usecase.ListMessagesUseCase
	updateMsgUC     *usecase.UpdateMessageUseCase
	deleteMsgUC     *usecase.DeleteMessageUseCase
	transitionMsgUC *usecase.TransitionMessageUseCase
}

func NewMessageController(
	createMsgUC *usecase.CreateMessageUseCase,
	getMsgUC *usecase.GetMessageUseCase,
	listMsgUC *usecase.ListMessagesUseCase,
	updateMsgUC *usecase.UpdateMessageUseCase,
	deleteMsgUC *usecase.DeleteMessageUseCase,
	transitionMsgUC *usecase.TransitionMessageUseCase,
) *MessageController {
	return &MessageController{
		createMsgUC:     createMsgUC,
		getMsgUC:        getMsgUC,
		listMsgUC:       listMsgUC,
		updateMsgUC:     updateMsgUC,
		deleteMsgUC:     deleteMsgUC,
		transitionMsgUC: transitionMsgUC,
	}
}

// CreateMessage godoc
// @Summary Create a new message
// @Description Creates a new message with optional image upload to GCS
// @Tags Messages
// @Accept multipart/form-data,json
// @Produce json
// @Security BearerAuth
// @Param name formData string false "Message Name"
// @Param content formData string false "Message Content"
// @Param file formData file false "Optional Image File"
// @Param input body dto.CreateMessageInput false "Or JSON payload"
// @Success 201 {object} presenter.ResponseDoc
// @Failure 400 {object} presenter.ResponseDoc
// @Failure 500 {object} presenter.ResponseDoc
// @Router /api/v1/messages [post]
func (ctrl *MessageController) CreateMessage(c fiber.Ctx) error {
	var input dto.CreateMessageInput
	var fileReader io.Reader
	var filename, contentType string

	contentTypeHeader := c.Get("Content-Type")
	if strings.Contains(contentTypeHeader, "multipart/form-data") {
		input.Name = c.FormValue("name")
		input.Content = c.FormValue("content")

		file, err := c.FormFile("file")
		if err == nil && file != nil {
			f, fErr := file.Open()
			if fErr == nil {
				defer f.Close()
				fileReader = f
				filename = file.Filename
				contentType = file.Header.Get("Content-Type")
			}
		}
	} else {
		if err := c.Bind().JSON(&input); err != nil {
			return presenter.RenderValidationError(c, "Invalid JSON payload", map[string][]string{
				"body": {"failed parsing input"},
			})
		}
	}

	result := ctrl.createMsgUC.Execute(c.Context(), input, fileReader, filename, contentType)
	if !result.Success {
		return presenter.RenderResponse(c, presenter.MapErrorToStatusCode(result.Error), result)
	}

	return presenter.RenderResponse(c, fiber.StatusCreated, result)
}

// GetMessage godoc
// @Summary Get message by ID
// @Description Retrieves message details (cached via L1/L2 cache)
// @Tags Messages
// @Produce json
// @Security BearerAuth
// @Param id path int true "Message Snowflake ID"
// @Success 200 {object} presenter.ResponseDoc
// @Failure 404 {object} presenter.ResponseDoc
// @Router /api/v1/messages/{id} [get]
func (ctrl *MessageController) GetMessage(c fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return presenter.RenderError(c, domain.ErrValidation, "Invalid message ID format")
	}

	result := ctrl.getMsgUC.Execute(c.Context(), id)
	if !result.Success {
		return presenter.RenderResponse(c, presenter.MapErrorToStatusCode(result.Error), result)
	}

	return presenter.RenderResponse(c, fiber.StatusOK, result)
}

// ListMessages godoc
// @Summary List messages
// @Description Retrieves paginated messages
// @Tags Messages
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Page size" default(20)
// @Success 200 {object} presenter.ResponseDoc
// @Failure 500 {object} presenter.ResponseDoc
// @Router /api/v1/messages [get]
func (ctrl *MessageController) ListMessages(c fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))

	result := ctrl.listMsgUC.Execute(c.Context(), page, limit)
	if !result.Success {
		return presenter.RenderResponse(c, presenter.MapErrorToStatusCode(result.Error), result)
	}

	return presenter.RenderResponse(c, fiber.StatusOK, result)
}

// UpdateMessage godoc
// @Summary Update message
// @Description Updates message fields with optimistic concurrency locking
// @Tags Messages
// @Accept multipart/form-data,json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Message Snowflake ID"
// @Param name formData string false "Updated Name"
// @Param content formData string false "Updated Content"
// @Param version formData int true "Current Version (Optimistic Lock)"
// @Param file formData file false "Updated Image File"
// @Param input body dto.UpdateMessageInput false "Or JSON payload"
// @Success 200 {object} presenter.ResponseDoc
// @Failure 409 {object} presenter.ResponseDoc
// @Router /api/v1/messages/{id} [put]
func (ctrl *MessageController) UpdateMessage(c fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return presenter.RenderError(c, domain.ErrValidation, "Invalid message ID format")
	}

	var input dto.UpdateMessageInput
	var fileReader io.Reader
	var filename, contentType string

	contentTypeHeader := c.Get("Content-Type")
	if strings.Contains(contentTypeHeader, "multipart/form-data") {
		if nameVal := c.FormValue("name"); nameVal != "" {
			input.Name = &nameVal
		}
		if contentVal := c.FormValue("content"); contentVal != "" {
			input.Content = &contentVal
		}
		if vVal := c.FormValue("version"); vVal != "" {
			v, _ := strconv.ParseInt(vVal, 10, 64)
			input.Version = v
		}

		file, err := c.FormFile("file")
		if err == nil && file != nil {
			f, fErr := file.Open()
			if fErr == nil {
				defer f.Close()
				fileReader = f
				filename = file.Filename
				contentType = file.Header.Get("Content-Type")
			}
		}
	} else {
		if err := c.Bind().JSON(&input); err != nil {
			return presenter.RenderValidationError(c, "Invalid JSON payload", map[string][]string{
				"body": {"failed parsing input"},
			})
		}
	}

	result := ctrl.updateMsgUC.Execute(c.Context(), id, input, fileReader, filename, contentType)
	if !result.Success {
		return presenter.RenderResponse(c, presenter.MapErrorToStatusCode(result.Error), result)
	}

	return presenter.RenderResponse(c, fiber.StatusOK, result)
}

// TransitionMessage godoc
// @Summary Transition message state
// @Description Changes message state following the State Transition Pattern (DRAFT -> PUBLISHED -> ARCHIVED)
// @Tags Messages
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Message Snowflake ID"
// @Param input body dto.TransitionMessageInput true "Transition Input"
// @Success 200 {object} presenter.ResponseDoc
// @Failure 422 {object} presenter.ResponseDoc
// @Router /api/v1/messages/{id}/transition [patch]
func (ctrl *MessageController) TransitionMessage(c fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return presenter.RenderError(c, domain.ErrValidation, "Invalid message ID format")
	}

	var input dto.TransitionMessageInput
	if err := c.Bind().JSON(&input); err != nil {
		return presenter.RenderValidationError(c, "Invalid JSON payload", map[string][]string{
			"body": {"failed parsing input"},
		})
	}

	result := ctrl.transitionMsgUC.Execute(c.Context(), id, input)
	if !result.Success {
		return presenter.RenderResponse(c, presenter.MapErrorToStatusCode(result.Error), result)
	}

	return presenter.RenderResponse(c, fiber.StatusOK, result)
}

// DeleteMessage godoc
// @Summary Delete message
// @Description Deletes a message and removes associated image from GCS
// @Tags Messages
// @Produce json
// @Security BearerAuth
// @Param id path int true "Message Snowflake ID"
// @Success 200 {object} presenter.ResponseDoc
// @Failure 404 {object} presenter.ResponseDoc
// @Router /api/v1/messages/{id} [delete]
func (ctrl *MessageController) DeleteMessage(c fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return presenter.RenderError(c, domain.ErrValidation, "Invalid message ID format")
	}

	result := ctrl.deleteMsgUC.Execute(c.Context(), id)
	if !result.Success {
		return presenter.RenderResponse(c, presenter.MapErrorToStatusCode(result.Error), result)
	}

	return presenter.RenderResponse(c, fiber.StatusOK, result)
}
