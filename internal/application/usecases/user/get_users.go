package userusecase

import (
	"like-api/internal/domain/entities"
	"like-api/internal/domain/repositories"
)

// GetUsersInput represents input for get users use case
type GetUsersInput struct{}

// GetUsersOutput represents output for get users use case
type GetUsersOutput struct {
	Users []entities.User
	Error error
}

// GetUsersUseCase defines the interface for getting all users
type GetUsersUseCase interface {
	Execute(input GetUsersInput) GetUsersOutput
}

// getUsersUseCaseImpl implements GetUsersUseCase
type getUsersUseCaseImpl struct {
	userRepo repositories.UserRepository
}

// NewGetUsersUseCase creates a new get users use case
func NewGetUsersUseCase(userRepo repositories.UserRepository) GetUsersUseCase {
	return &getUsersUseCaseImpl{
		userRepo: userRepo,
	}
}

// Execute gets all users
func (uc *getUsersUseCaseImpl) Execute(input GetUsersInput) GetUsersOutput {
	users, err := uc.userRepo.FindAll()
	if err != nil {
		return GetUsersOutput{
			Users: nil,
			Error: err,
		}
	}

	return GetUsersOutput{
		Users: users,
		Error: nil,
	}
}
