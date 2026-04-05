package Service

import (
	"fmt"

	"task/internal/Models"
	"task/internal/Repository"

	"github.com/google/uuid"
)

type IAuthService interface {
	GetUserByID(id uuid.UUID) (*Models.User, error)
}

type AuthService struct {
	userRepo Repository.IUserRepository
}

func NewAuthService(userRepo Repository.IUserRepository) IAuthService {
	return &AuthService{userRepo: userRepo}
}

func (s *AuthService) GetUserByID(id uuid.UUID) (*Models.User, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}
	return user, nil
}
