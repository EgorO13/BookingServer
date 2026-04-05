package Service

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"task/internal/Models"
)

type mockUserRepo struct {
	mock.Mock
}

func (m *mockUserRepo) FindByID(id uuid.UUID) (*Models.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Models.User), args.Error(1)
}

func TestAuthService_GetUserByID_Success(t *testing.T) {
	mockRepo := new(mockUserRepo)
	svc := NewAuthService(mockRepo)
	userID := uuid.New()
	expectedUser := &Models.User{ID: userID, Email: "test@example.com", Role: "user"}
	mockRepo.On("FindByID", userID).Return(expectedUser, nil)
	user, err := svc.GetUserByID(userID)
	assert.NoError(t, err)
	assert.Equal(t, expectedUser, user)
	mockRepo.AssertExpectations(t)
}

func TestAuthService_GetUserByID_NotFound(t *testing.T) {
	mockRepo := new(mockUserRepo)
	svc := NewAuthService(mockRepo)
	userID := uuid.New()
	mockRepo.On("FindByID", userID).Return(nil, nil)
	user, err := svc.GetUserByID(userID)
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Contains(t, err.Error(), "user not found")
	mockRepo.AssertExpectations(t)
}

func TestAuthService_GetUserByID_RepoError(t *testing.T) {
	mockRepo := new(mockUserRepo)
	svc := NewAuthService(mockRepo)
	userID := uuid.New()
	mockRepo.On("FindByID", userID).Return(nil, assert.AnError)
	user, err := svc.GetUserByID(userID)
	assert.Error(t, err)
	assert.Nil(t, user)
	mockRepo.AssertExpectations(t)
}
