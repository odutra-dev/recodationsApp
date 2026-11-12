package usecase

import (
	"bytes"
	"context"
	"io"
	"testing"
	"time"

	"github.com/dutra/recodationsApp/internal/modules/messages/domain"
	"github.com/dutra/recodationsApp/internal/modules/messages/dto"
	sharedDomain "github.com/dutra/recodationsApp/internal/shared/domain"
	"github.com/dutra/recodationsApp/internal/shared/infrastructure/cache"
	"github.com/dutra/recodationsApp/internal/shared/infrastructure/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockMessageRepository struct {
	mock.Mock
}

func (m *mockMessageRepository) Create(ctx context.Context, msg *domain.Message) error {
	args := m.Called(ctx, msg)
	return args.Error(0)
}

func (m *mockMessageRepository) FindByID(ctx context.Context, id int64) (*domain.Message, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Message), args.Error(1)
}

func (m *mockMessageRepository) List(ctx context.Context, page, limit int) ([]domain.Message, int64, error) {
	args := m.Called(ctx, page, limit)
	return args.Get(0).([]domain.Message), args.Get(1).(int64), args.Error(2)
}

func (m *mockMessageRepository) Update(ctx context.Context, msg *domain.Message) error {
	args := m.Called(ctx, msg)
	return args.Error(0)
}

func (m *mockMessageRepository) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type mockStorageService struct {
	mock.Mock
}

func (m *mockStorageService) EnsureBucket(ctx context.Context, bucket string) error {
	return nil
}
func (m *mockStorageService) UploadFile(ctx context.Context, bucket, filename, contentType string, reader io.Reader) (*storage.UploadResult, error) {
	args := m.Called(ctx, bucket, filename, contentType, reader)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*storage.UploadResult), args.Error(1)
}
func (m *mockStorageService) DeleteFile(ctx context.Context, bucket, filename string) error {
	args := m.Called(ctx, bucket, filename)
	return args.Error(0)
}
func (m *mockStorageService) GetFileURL(bucket, filename string) string {
	return "http://gcs/" + filename
}

func TestCreateMessageUseCase(t *testing.T) {
	repo := new(mockMessageRepository)
	stor := new(mockStorageService)
	memCache := cache.NewMemoryCache()
	uc := NewCreateMessageUseCase(repo, stor, nil, memCache, "photos")

	ctx := context.Background()
	input := dto.CreateMessageInput{
		Name:    "Greeting",
		Content: "Hello World",
	}

	reader := bytes.NewReader([]byte("fake image data"))
	stor.On("UploadFile", ctx, "photos", mock.AnythingOfType("string"), "image/png", reader).
		Return(&storage.UploadResult{URL: "http://gcs/uploaded.png"}, nil)

	repo.On("Create", ctx, mock.AnythingOfType("*domain.Message")).Return(nil)

	res := uc.Execute(ctx, input, reader, "test.png", "image/png")
	require.True(t, res.Success)
	assert.Equal(t, "Greeting", res.Data.Name)
	assert.Equal(t, "Hello World", res.Data.Content)
	assert.Equal(t, "DRAFT", res.Data.Status)
	assert.Equal(t, int64(1), res.Data.Version)
	assert.NotNil(t, res.Data.Image)
	assert.Equal(t, "http://gcs/uploaded.png", *res.Data.Image)
	repo.AssertExpectations(t)
	stor.AssertExpectations(t)
}

func TestTransitionMessageUseCase(t *testing.T) {
	repo := new(mockMessageRepository)
	memCache := cache.NewMemoryCache()
	uc := NewTransitionMessageUseCase(repo, memCache)

	ctx := context.Background()
	msg := &domain.Message{
		ID:        99,
		Name:      "Note",
		Content:   "Draft note",
		Status:    domain.MessageStatusDraft,
		Version:   2,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	repo.On("FindByID", ctx, int64(99)).Return(msg, nil)
	repo.On("Update", ctx, mock.MatchedBy(func(m *domain.Message) bool {
		return m.ID == 99 && m.Status == domain.MessageStatusPublished
	})).Return(nil)

	// Valid transition: DRAFT -> PUBLISHED
	res := uc.Execute(ctx, 99, dto.TransitionMessageInput{
		Status:  "PUBLISHED",
		Version: 2,
	})
	require.True(t, res.Success)
	assert.Equal(t, "PUBLISHED", res.Data.Status)
}

func TestTransitionMessageInvalidTransition(t *testing.T) {
	repo := new(mockMessageRepository)
	memCache := cache.NewMemoryCache()
	uc := NewTransitionMessageUseCase(repo, memCache)

	ctx := context.Background()
	msg := &domain.Message{
		ID:      100,
		Status:  domain.MessageStatusArchived, // Terminal
		Version: 1,
	}

	repo.On("FindByID", ctx, int64(100)).Return(msg, nil)

	// Invalid transition: ARCHIVED -> PUBLISHED
	res := uc.Execute(ctx, 100, dto.TransitionMessageInput{
		Status:  "PUBLISHED",
		Version: 1,
	})
	assert.False(t, res.Success)
	assert.ErrorIs(t, res.Error, sharedDomain.ErrInvalidStateTransition)
}

func TestUpdateMessageOptimisticLockConflict(t *testing.T) {
	repo := new(mockMessageRepository)
	memCache := cache.NewMemoryCache()
	uc := NewUpdateMessageUseCase(repo, nil, memCache, "photos")

	ctx := context.Background()
	msg := &domain.Message{
		ID:      101,
		Version: 3, // Current DB version is 3
	}

	repo.On("FindByID", ctx, int64(101)).Return(msg, nil)

	name := "Changed Name"
	res := uc.Execute(ctx, 101, dto.UpdateMessageInput{
		Name:    &name,
		Version: 2, // Client sent stale version 2
	}, nil, "", "")

	assert.False(t, res.Success)
	assert.ErrorIs(t, res.Error, sharedDomain.ErrOptimisticLockConflict)
}
