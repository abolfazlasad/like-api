package authusecase

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"like-api/internal/application/dto/request"
	"like-api/internal/application/dto/response"
	"like-api/internal/domain/entities"
	"like-api/internal/domain/repositories"
)

type RegisterUseCase struct {
	userRepo  repositories.UserRepository
	jwtSecret string
}

func NewRegisterUseCase(userRepo repositories.UserRepository, jwtSecret string) *RegisterUseCase {
	return &RegisterUseCase{userRepo: userRepo, jwtSecret: jwtSecret}
}

func (uc *RegisterUseCase) Execute(req request.RegisterRequest) (*response.AuthResponse, error) {
	if uc.userRepo.ExistsByEmail(req.Email) {
		return nil, errors.New("email already registered")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := entities.User{
		ID:        uuid.New().String(),
		Username:  req.Username,
		Name:      req.Name,
		Email:     req.Email,
		Password:  string(hash),
		Role:      entities.RoleUser,
		CreatedAt: time.Now().Format(time.RFC3339Nano),
	}

	if err := uc.userRepo.Create(user); err != nil {
		return nil, err
	}

	token, err := GenerateToken(user, uc.jwtSecret)
	if err != nil {
		return nil, err
	}

	return &response.AuthResponse{
		Token: token,
		User: response.UserDTO{
			ID:       user.ID,
			Username: user.Username,
			Name:     user.Name,
			Email:    user.Email,
			Role:     string(user.Role),
		},
	}, nil
}
