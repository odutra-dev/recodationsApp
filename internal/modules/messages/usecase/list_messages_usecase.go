package usecase

import (
	"context"

	"github.com/dutra/recodationsApp/internal/modules/messages/dto"
	"github.com/dutra/recodationsApp/internal/modules/messages/repository"
	sharedDomain "github.com/dutra/recodationsApp/internal/shared/domain"
)

type ListMessagesUseCase struct {
	messageRepo repository.MessageRepository
}

func NewListMessagesUseCase(messageRepo repository.MessageRepository) *ListMessagesUseCase {
	return &ListMessagesUseCase{messageRepo: messageRepo}
}

func (uc *ListMessagesUseCase) Execute(ctx context.Context, page, limit int) sharedDomain.Result[[]dto.MessageResponse] {
	messages, _, err := uc.messageRepo.List(ctx, page, limit)
	if err != nil {
		return sharedDomain.Fail[[]dto.MessageResponse](err, "failed fetching messages")
	}

	response := make([]dto.MessageResponse, len(messages))
	for i, m := range messages {
		response[i] = dto.MessageResponse{
			ID:        m.ID,
			Name:      m.Name,
			Content:   m.Content,
			Image:     m.Image,
			Status:    string(m.Status),
			Version:   m.Version,
			CreatedAt: m.CreatedAt,
			UpdatedAt: m.UpdatedAt,
		}
	}

	return sharedDomain.Ok(response, "Messages retrieved successfully")
}
