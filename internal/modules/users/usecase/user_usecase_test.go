package usecase

import (
	"context"
	"testing"
	"time"

	rulesDomain "github.com/dutra/recodationsApp/internal/modules/rules/domain"
	"github.com/dutra/recodationsApp/internal/modules/users/domain"
	"github.com/dutra/recodationsApp/internal/modules/users/dto"
	sharedDomain "github.com/dutra/recodationsApp/internal/shared/domain"
	"github.com/dutra/recodationsApp/internal/shared/infrastructure/cache"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockUserRepository struct {
	mock.Mock
}

func (m *mockUserRepository) Create(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *mockUserRepository) FindByID(ctx context.Context, id int64) (*domain.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *mockUserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *mockUserRepository) List(ctx context.Context, page, limit int) ([]domain.User, int64, error) {
	args := m.Called(ctx, page, limit)
	return args.Get(0).([]domain.User), args.Get(1).(int64), args.Error(2)
}

func (m *mockUserRepository) Update(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *mockUserRepository) AssignRules(ctx context.Context, userID int64, ruleIDs []int64) error {
	args := m.Called(ctx, userID, ruleIDs)
	return args.Error(0)
}

func TestCreateUserUseCaseSuccess(t *testing.T) {
	repo := new(mockUserRepository)
	uc := NewCreateUserUseCase(repo)

	ctx := context.Background()
	input := dto.CreateUserInput{
		Name:     "Alice Smith",
		Email:    "alice@example.com",
		Password: "password123",
	}

	repo.On("FindByEmail", ctx, "alice@example.com").Return(nil, sharedDomain.ErrNotFound)
	repo.On("Create", ctx, mock.AnythingOfType("*domain.User")).Return(nil)

	result := uc.Execute(ctx, input)
	require.True(t, result.Success)
	assert.Equal(t, "alice@example.com", result.Data.Email)
	assert.Equal(t, "Alice Smith", result.Data.Name)
	assert.Positive(t, result.Data.ID)
	assert.Equal(t, int64(1), result.Data.Version)
	repo.AssertExpectations(t)
}

func TestCreateUserValidationFailure(t *testing.T) {
	repo := new(mockUserRepository)
	uc := NewCreateUserUseCase(repo)

	ctx := context.Background()
	input := dto.CreateUserInput{
		Name:     "",
		Email:    "invalid-email",
		Password: "123",
	}

	result := uc.Execute(ctx, input)
	assert.False(t, result.Success)
	assert.Contains(t, result.ErrorValids, "name")
	assert.Contains(t, result.ErrorValids, "email")
	assert.Contains(t, result.ErrorValids, "password")
}

func TestGetUserUseCaseWithCache(t *testing.T) {
	repo := new(mockUserRepository)
	memCache := cache.NewMemoryCache()
	uc := NewGetUserUseCase(repo, memCache)

	ctx := context.Background()
	user := &domain.User{
		ID:        42,
		Name:      "Bob",
		Email:     "bob@example.com",
		Version:   1,
		Rules:     []rulesDomain.Rule{{ID: 1, Name: "Admin", Code: "admin"}},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	repo.On("FindByID", ctx, int64(42)).Return(user, nil).Once()

	// 1st call -> Cache miss, queries repo
	result1 := uc.Execute(ctx, 42)
	require.True(t, result1.Success)
	assert.Equal(t, "Bob", result1.Data.Name)

	// 2nd call -> Cache hit, repo is NOT called again
	result2 := uc.Execute(ctx, 42)
	require.True(t, result2.Success)
	assert.Equal(t, "Bob", result2.Data.Name)

	repo.AssertNumberOfCalls(t, "FindByID", 1)
}
