package usecase

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/dutra/recodationsApp/internal/modules/messages/dto"
	"github.com/dutra/recodationsApp/internal/modules/messages/repository"
	sharedDomain "github.com/dutra/recodationsApp/internal/shared/domain"
	"github.com/dutra/recodationsApp/internal/shared/infrastructure/cache"
	"github.com/dutra/recodationsApp/internal/shared/infrastructure/storage"
)

type UpdateMessageUseCase struct {
	messageRepo repository.MessageRepository
	storageSvc  storage.StorageService
	cache       cache.Cache
	bucketName  string
}

func NewUpdateMessageUseCase(
	messageRepo repository.MessageRepository,
	storageSvc storage.StorageService,
	cache cache.Cache,
	bucketName string,
) *UpdateMessageUseCase {
	return &UpdateMessageUseCase{
		messageRepo: messageRepo,
		storageSvc:  storageSvc,
		cache:       cache,
		bucketName:  bucketName,
	}
}

func (uc *UpdateMessageUseCase) Execute(
	ctx context.Context,
	id int64,
	input dto.UpdateMessageInput,
	fileReader io.Reader,
	filename string,
	contentType string,
) sharedDomain.Result[dto.MessageResponse] {
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

	if input.Name != nil && *input.Name != "" {
		message.Name = *input.Name
	}
	if input.Content != nil && *input.Content != "" {
		message.Content = *input.Content
	}

	if fileReader != nil && uc.storageSvc != nil && filename != "" {
		uniqueName := fmt.Sprintf("%d-%s", time.Now().UnixNano(), filename)
		uploadRes, err := uc.storageSvc.UploadFile(ctx, uc.bucketName, uniqueName, contentType, fileReader)
		if err != nil {
			return sharedDomain.Fail[dto.MessageResponse](err, "failed uploading new image")
		}
		message.Image = &uploadRes.URL
	} else if input.Image != nil {
		message.Image = input.Image
	}

	if err := uc.messageRepo.Update(ctx, message); err != nil {
		return sharedDomain.Fail[dto.MessageResponse](err, "failed updating message")
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

	return sharedDomain.Ok(response, "Message updated successfully")
}
