package usecase

import (
	"context"
	"fmt"
	"strings"

	"github.com/dutra/recodationsApp/internal/modules/messages/repository"
	sharedDomain "github.com/dutra/recodationsApp/internal/shared/domain"
	"github.com/dutra/recodationsApp/internal/shared/infrastructure/cache"
	"github.com/dutra/recodationsApp/internal/shared/infrastructure/storage"
)

type DeleteMessageUseCase struct {
	messageRepo repository.MessageRepository
	storageSvc  storage.StorageService
	cache       cache.Cache
	bucketName  string
}

func NewDeleteMessageUseCase(
	messageRepo repository.MessageRepository,
	storageSvc storage.StorageService,
	cache cache.Cache,
	bucketName string,
) *DeleteMessageUseCase {
	return &DeleteMessageUseCase{
		messageRepo: messageRepo,
		storageSvc:  storageSvc,
		cache:       cache,
		bucketName:  bucketName,
	}
}

func (uc *DeleteMessageUseCase) Execute(ctx context.Context, id int64) sharedDomain.Result[string] {
	if id <= 0 {
		return sharedDomain.Fail[string](sharedDomain.ErrValidation, "invalid message id")
	}

	message, err := uc.messageRepo.FindByID(ctx, id)
	if err != nil {
		return sharedDomain.Fail[string](err, "message not found")
	}

	// Delete image from GCS if exists
	if message.Image != nil && uc.storageSvc != nil {
		parts := strings.Split(*message.Image, "/")
		filename := parts[len(parts)-1]
		_ = uc.storageSvc.DeleteFile(ctx, uc.bucketName, filename)
	}

	if err := uc.messageRepo.Delete(ctx, id); err != nil {
		return sharedDomain.Fail[string](err, "failed deleting message")
	}

	// Invalidate L1/L2 Cache
	if uc.cache != nil {
		_ = uc.cache.Delete(ctx, fmt.Sprintf("message:%d", id))
	}

	return sharedDomain.Ok("Message deleted successfully", "Success")
}
