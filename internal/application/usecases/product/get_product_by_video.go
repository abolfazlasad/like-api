package productusecase

import (
	"like-api/internal/domain/entities"
	"like-api/internal/domain/repositories"
)

type GetProductByVideoInput struct {
	VideoID string
}

type GetProductByVideoOutput struct {
	Product       entities.Product
	Error         error
	VideoExists   bool
	ProductExists bool
}

type GetProductByVideoUseCase interface {
	Execute(input GetProductByVideoInput) GetProductByVideoOutput
}

type getProductByVideoUseCaseImpl struct {
	videoRepo   repositories.VideoRepository
	productRepo repositories.ProductRepository
}

func NewGetProductByVideoUseCase(
	videoRepo repositories.VideoRepository,
	productRepo repositories.ProductRepository,
) GetProductByVideoUseCase {
	return &getProductByVideoUseCaseImpl{
		videoRepo:   videoRepo,
		productRepo: productRepo,
	}
}

func (uc *getProductByVideoUseCaseImpl) Execute(input GetProductByVideoInput) GetProductByVideoOutput {
	if !uc.videoRepo.Exists(input.VideoID) {
		return GetProductByVideoOutput{VideoExists: false}
	}

	product, err := uc.productRepo.FindByVideoID(input.VideoID)
	if err != nil {
		return GetProductByVideoOutput{VideoExists: true, ProductExists: false}
	}

	return GetProductByVideoOutput{
		Product:       product,
		VideoExists:   true,
		ProductExists: true,
	}
}
