package authusecase

import (
	"errors"

	"golang.org/x/crypto/bcrypt"

	"like-api/internal/application/dto/request"
	"like-api/internal/application/dto/response"
	"like-api/internal/domain/repositories"
)

type LoginUseCase struct {
	userRepo  repositories.UserRepository
	jwtSecret string
}

func NewLoginUseCase(userRepo repositories.UserRepository, jwtSecret string) *LoginUseCase {
	return &LoginUseCase{userRepo: userRepo, jwtSecret: jwtSecret}
}

func (uc *LoginUseCase) Execute(req request.LoginRequest) (*response.AuthResponse, error) {
	user, err := uc.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	token, err := generateToken(user, uc.jwtSecret)
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
		},
	}, nil
}
