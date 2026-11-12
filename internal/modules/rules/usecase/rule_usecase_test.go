package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/dutra/recodationsApp/internal/modules/rules/domain"
	"github.com/dutra/recodationsApp/internal/modules/rules/dto"
	sharedDomain "github.com/dutra/recodationsApp/internal/shared/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockRuleRepository struct {
	mock.Mock
}

func (m *mockRuleRepository) Create(ctx context.Context, rule *domain.Rule) error {
	args := m.Called(ctx, rule)
	return args.Error(0)
}

func (m *mockRuleRepository) FindByID(ctx context.Context, id int64) (*domain.Rule, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Rule), args.Error(1)
}

func (m *mockRuleRepository) FindByCode(ctx context.Context, code string) (*domain.Rule, error) {
	args := m.Called(ctx, code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Rule), args.Error(1)
}

func (m *mockRuleRepository) List(ctx context.Context) ([]domain.Rule, error) {
	args := m.Called(ctx)
	return args.Get(0).([]domain.Rule), args.Error(1)
}

func (m *mockRuleRepository) FindRulesByUserID(ctx context.Context, userID int64) ([]domain.Rule, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]domain.Rule), args.Error(1)
}

func (m *mockRuleRepository) AssignRuleToUser(ctx context.Context, userID int64, ruleID int64) error {
	args := m.Called(ctx, userID, ruleID)
	return args.Error(0)
}

func (m *mockRuleRepository) RemoveRuleFromUser(ctx context.Context, userID int64, ruleID int64) error {
	args := m.Called(ctx, userID, ruleID)
	return args.Error(0)
}

func TestCreateRuleUseCaseSuccess(t *testing.T) {
	repo := new(mockRuleRepository)
	uc := NewCreateRuleUseCase(repo)

	ctx := context.Background()
	input := dto.CreateRuleInput{
		Name:        "Test Permission",
		Code:        "test:perm",
		Description: "A unit test rule",
	}

	repo.On("FindByCode", ctx, "test:perm").Return(nil, sharedDomain.ErrNotFound)
	repo.On("Create", ctx, mock.AnythingOfType("*domain.Rule")).Return(nil)

	result := uc.Execute(ctx, input)
	require.True(t, result.Success)
	assert.Equal(t, "test:perm", result.Data.Code)
	assert.Equal(t, "Test Permission", result.Data.Name)
	assert.Positive(t, result.Data.ID)
	repo.AssertExpectations(t)
}

func TestCreateRuleUseCaseValidationFail(t *testing.T) {
	repo := new(mockRuleRepository)
	uc := NewCreateRuleUseCase(repo)

	ctx := context.Background()
	input := dto.CreateRuleInput{
		Name: "",
		Code: "",
	}

	result := uc.Execute(ctx, input)
	assert.False(t, result.Success)
	assert.NotEmpty(t, result.ErrorValids)
	assert.Contains(t, result.ErrorValids, "name")
	assert.Contains(t, result.ErrorValids, "code")
}

func TestListRulesUseCaseSuccess(t *testing.T) {
	repo := new(mockRuleRepository)
	uc := NewListRulesUseCase(repo)

	ctx := context.Background()
	rules := []domain.Rule{
		{ID: 1, Name: "Read", Code: "read", CreatedAt: time.Now()},
		{ID: 2, Name: "Write", Code: "write", CreatedAt: time.Now()},
	}

	repo.On("List", ctx).Return(rules, nil)

	result := uc.Execute(ctx)
	require.True(t, result.Success)
	assert.Len(t, *result.Data, 2)
	repo.AssertExpectations(t)
}
