package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/dutra/recodationsApp/internal/modules/messages/dto"
	"github.com/dutra/recodationsApp/internal/modules/messages/repository"
	sharedDomain "github.com/dutra/recodationsApp/internal/shared/domain"
	"github.com/dutra/recodationsApp/internal/shared/infrastructure/cache"
)

type GetMessageUseCase struct {
	messageRepo repository.MessageRepository
	cache       cache.Cache
}

func NewGetMessageUseCase(messageRepo repository.MessageRepository, cache cache.Cache) *GetMessageUseCase {
	return &GetMessageUseCase{
		messageRepo: messageRepo,
		cache:       cache,
	}
}

func (uc *GetMessageUseCase) Execute(ctx context.Context, id int64) sharedDomain.Result[dto.MessageResponse] {
	if id <= 0 {
		return sharedDomain.Fail[dto.MessageResponse](sharedDomain.ErrValidation, "invalid message id")
	}

	cacheKey := fmt.Sprintf("message:%d", id)

	// 1. Check L1/L2 Cache
	if uc.cache != nil {
		var cached dto.MessageResponse
		if err := uc.cache.Get(ctx, cacheKey, &cached); err == nil && cached.ID != 0 {
			return sharedDomain.Ok(cached, "Message retrieved from cache")
		}
	}

	// 2. Query Repository
	message, err := uc.messageRepo.FindByID(ctx, id)
	if err != nil {
		return sharedDomain.Fail[dto.MessageResponse](err, "message not found")
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

	// 3. Populate L1/L2 Cache
	if uc.cache != nil {
		_ = uc.cache.Set(ctx, cacheKey, response, 10*time.Minute)
	}

	return sharedDomain.Ok(response, "Message retrieved successfully")
}
