package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/dutra/recodationsApp/internal/modules/auth/dto"
	rulesDomain "github.com/dutra/recodationsApp/internal/modules/rules/domain"
	usersDomain "github.com/dutra/recodationsApp/internal/modules/users/domain"
	sharedDomain "github.com/dutra/recodationsApp/internal/shared/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

type mockUserRepoForAuth struct {
	mock.Mock
}

func (m *mockUserRepoForAuth) Create(ctx context.Context, user *usersDomain.User) error {
	return nil
}
func (m *mockUserRepoForAuth) FindByID(ctx context.Context, id int64) (*usersDomain.User, error) {
	return nil, nil
}
func (m *mockUserRepoForAuth) FindByEmail(ctx context.Context, email string) (*usersDomain.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*usersDomain.User), args.Error(1)
}
func (m *mockUserRepoForAuth) List(ctx context.Context, page, limit int) ([]usersDomain.User, int64, error) {
	return nil, 0, nil
}
func (m *mockUserRepoForAuth) Update(ctx context.Context, user *usersDomain.User) error {
	return nil
}
func (m *mockUserRepoForAuth) AssignRules(ctx context.Context, userID int64, ruleIDs []int64) error {
	return nil
}

func TestLoginUseCaseSuccess(t *testing.T) {
	repo := new(mockUserRepoForAuth)
	uc := NewLoginUseCase(repo, "secret-key", 24*time.Hour)

	ctx := context.Background()
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	require.NoError(t, err)

	user := &usersDomain.User{
		ID:           101,
		Name:         "John Doe",
		Email:        "john@example.com",
		PasswordHash: string(hashedPassword),
		Version:      1,
		Rules: []rulesDomain.Rule{
			{ID: 1, Name: "Read", Code: "messages:read"},
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	repo.On("FindByEmail", ctx, "john@example.com").Return(user, nil)

	result := uc.Execute(ctx, dto.LoginInput{
		Email:    "john@example.com",
		Password: "password123",
	})

	require.True(t, result.Success)
	assert.NotEmpty(t, result.Data.AccessToken)
	assert.Equal(t, "Bearer", result.Data.TokenType)
	assert.Equal(t, "john@example.com", result.Data.User.Email)
	assert.Len(t, result.Data.User.Rules, 1)
	repo.AssertExpectations(t)
}

func TestLoginUseCaseInvalidPassword(t *testing.T) {
	repo := new(mockUserRepoForAuth)
	uc := NewLoginUseCase(repo, "secret-key", 24*time.Hour)

	ctx := context.Background()
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("correct-password"), bcrypt.DefaultCost)

	user := &usersDomain.User{
		ID:           101,
		Email:        "john@example.com",
		PasswordHash: string(hashedPassword),
	}

	repo.On("FindByEmail", ctx, "john@example.com").Return(user, nil)

	result := uc.Execute(ctx, dto.LoginInput{
		Email:    "john@example.com",
		Password: "wrong-password",
	})

	assert.False(t, result.Success)
	assert.ErrorIs(t, result.Error, sharedDomain.ErrInvalidCredentials)
}
