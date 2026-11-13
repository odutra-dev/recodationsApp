package usecase

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	inboxUC "github.com/dutra/recodationsApp/internal/modules/inbox/usecase"
	"github.com/dutra/recodationsApp/internal/modules/messages/domain"
	"github.com/dutra/recodationsApp/internal/modules/messages/dto"
	"github.com/dutra/recodationsApp/internal/modules/messages/repository"
	sharedDomain "github.com/dutra/recodationsApp/internal/shared/domain"
	"github.com/dutra/recodationsApp/internal/shared/infrastructure/cache"
	"github.com/dutra/recodationsApp/internal/shared/infrastructure/storage"
	"github.com/dutra/recodationsApp/pkg/snowflake"
)

type CreateMessageUseCase struct {
	messageRepo  repository.MessageRepository
	storageSvc   storage.StorageService
	inboxSvc     *inboxUC.InboxService
	cache        cache.Cache
	bucketName   string
}

func NewCreateMessageUseCase(
	messageRepo repository.MessageRepository,
	storageSvc storage.StorageService,
	inboxSvc *inboxUC.InboxService,
	cache cache.Cache,
	bucketName string,
) *CreateMessageUseCase {
	return &CreateMessageUseCase{
		messageRepo: messageRepo,
		storageSvc:  storageSvc,
		inboxSvc:    inboxSvc,
		cache:       cache,
		bucketName:  bucketName,
	}
}

func (uc *CreateMessageUseCase) Execute(
	ctx context.Context,
	input dto.CreateMessageInput,
	fileReader io.Reader,
	filename string,
	contentType string,
) sharedDomain.Result[dto.MessageResponse] {
	valErrors := make(map[string][]string)

	name := strings.TrimSpace(input.Name)
	if name == "" {
		valErrors["name"] = append(valErrors["name"], "name is required")
	}

	content := strings.TrimSpace(input.Content)
	if content == "" {
		valErrors["content"] = append(valErrors["content"], "content is required")
	}

	if len(valErrors) > 0 {
		return sharedDomain.FailWithValidation[dto.MessageResponse]("validation failed", valErrors)
	}

	var imageURL *string
	if fileReader != nil && uc.storageSvc != nil && filename != "" {
		uniqueName := fmt.Sprintf("%d-%s", time.Now().UnixNano(), filename)
		uploadRes, err := uc.storageSvc.UploadFile(ctx, uc.bucketName, uniqueName, contentType, fileReader)
		if err != nil {
			return sharedDomain.Fail[dto.MessageResponse](err, "failed uploading image to storage")
		}
		imageURL = &uploadRes.URL
	} else if input.Image != nil {
		imageURL = input.Image
	}

	now := time.Now().UTC()
	message := &domain.Message{
		ID:        snowflake.Generate(),
		Name:      name,
		Content:   content,
		Image:     imageURL,
		Status:    domain.MessageStatusDraft,
		Version:   1,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := uc.messageRepo.Create(ctx, message); err != nil {
		return sharedDomain.Fail[dto.MessageResponse](err, "failed saving message")
	}

	// Emit Inbox Event (Inbox Pattern for idempotency and asynchronous handlers)
	if uc.inboxSvc != nil {
		eventID := fmt.Sprintf("msg-created-%d", message.ID)
		_ = uc.inboxSvc.SaveEvent(ctx, eventID, "message.created", map[string]any{
			"messageId": message.ID,
			"name":      message.Name,
			"hasImage":  imageURL != nil,
		})
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

	// Cache entry
	if uc.cache != nil {
		_ = uc.cache.Set(ctx, fmt.Sprintf("message:%d", message.ID), response, 10*time.Minute)
	}

	return sharedDomain.Ok(response, "Message created successfully")
}
