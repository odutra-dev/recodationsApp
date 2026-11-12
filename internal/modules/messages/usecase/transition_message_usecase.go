package usecase

import (
	"context"
	"fmt"

	"github.com/dutra/recodationsApp/internal/modules/messages/domain"
	"github.com/dutra/recodationsApp/internal/modules/messages/dto"
	"github.com/dutra/recodationsApp/internal/modules/messages/repository"
	sharedDomain "github.com/dutra/recodationsApp/internal/shared/domain"
	"github.com/dutra/recodationsApp/internal/shared/infrastructure/cache"
)

type TransitionMessageUseCase struct {
	messageRepo repository.MessageRepository
	cache       cache.Cache
}

func NewTransitionMessageUseCase(messageRepo repository.MessageRepository, cache cache.Cache) *TransitionMessageUseCase {
	return &TransitionMessageUseCase{
		messageRepo: messageRepo,
		cache:       cache,
	}
}

func (uc *TransitionMessageUseCase) Execute(ctx context.Context, id int64, input dto.TransitionMessageInput) sharedDomain.Result[dto.MessageResponse] {
	if id <= 0 {
		return sharedDomain.Fail[dto.MessageResponse](sharedDomain.ErrValidation, "invalid message id")
	}

	message, err := uc.messageRepo.FindByID(ctx, id)
	if err != nil {
		return sharedDomain.Fail[dto.MessageResponse](err, "message not found")
	}

	// Verify optimistic lock version
	if input.Version > 0 && input.Version != message.Version {
		return sharedDomain.Fail[dto.MessageResponse](sharedDomain.ErrOptimisticLockConflict, "record was modified by another transaction")
	}

	targetStatus := domain.MessageStatus(input.Status)

	// Validate state transition
	if err := message.TransitionTo(targetStatus); err != nil {
		return sharedDomain.Fail[dto.MessageResponse](err, "invalid state transition")
	}

	if err := uc.messageRepo.Update(ctx, message); err != nil {
		return sharedDomain.Fail[dto.MessageResponse](err, "failed updating message status")
	}

	// Invalidate L1/L2 Cache
	if uc.cache != nil {
		_ = uc.cache.Delete(ctx, fmt.Sprintf("message:%d", id))
	}

	response := dto.MessageResponse{
		ID:        message.ID,
		Name:      message.Name,
		Content:   message.Content,
		Image:     message.Image,
		Status:    string(message.Status),
		Version:   message.Version,
		CreatedAt: message.CreatedAt,
		UpdatedAt: message.UpdatedAt,
	}

	return sharedDomain.Ok(response, "Message state transitioned successfully")
}
