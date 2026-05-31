package productusecase

import (
	"errors"
	"like-api/internal/domain/entities"
	"like-api/internal/domain/repositories"

	"github.com/google/uuid"
)

type CreateProductInput struct {
	VideoID  string
	Name     string
	Price    float64
	ImageURL string
}

type CreateProductOutput struct {
	Product     entities.Product
	Error       error
	VideoExists bool
}

type CreateProductUseCase interface {
	Execute(input CreateProductInput) CreateProductOutput
}

type createProductUseCaseImpl struct {
	videoRepo   repositories.VideoRepository
	productRepo repositories.ProductRepository
}

func NewCreateProductUseCase(
	videoRepo repositories.VideoRepository,
	productRepo repositories.ProductRepository,
) CreateProductUseCase {
	return &createProductUseCaseImpl{
		videoRepo:   videoRepo,
		productRepo: productRepo,
	}
}

func (uc *createProductUseCaseImpl) Execute(input CreateProductInput) CreateProductOutput {
	if !uc.videoRepo.Exists(input.VideoID) {
		return CreateProductOutput{VideoExists: false}
	}

	// Prevent attaching a second product to the same video.
	if _, err := uc.productRepo.FindByVideoID(input.VideoID); err == nil {
		return CreateProductOutput{
			VideoExists: true,
			Error:       errors.New("video already has an associated product"),
		}
	}

	product := entities.Product{
		ID:       uuid.New().String(),
		VideoID:  input.VideoID,
		Name:     input.Name,
		Price:    input.Price,
		ImageURL: input.ImageURL,
	}

	if err := uc.productRepo.Create(product); err != nil {
		return CreateProductOutput{VideoExists: true, Error: err}
	}

	return CreateProductOutput{Product: product, VideoExists: true}
}
