package productusecase

import (
	"like-api/internal/domain/entities"
	"like-api/internal/domain/repositories"
)

type GetProductInput struct {
	ProductID string
}

type GetProductOutput struct {
	Product       entities.Product
	Error         error
	ProductExists bool
}

type GetProductUseCase interface {
	Execute(input GetProductInput) GetProductOutput
}

type getProductUseCaseImpl struct {
	productRepo repositories.ProductRepository
}

func NewGetProductUseCase(productRepo repositories.ProductRepository) GetProductUseCase {
	return &getProductUseCaseImpl{productRepo: productRepo}
}

func (uc *getProductUseCaseImpl) Execute(input GetProductInput) GetProductOutput {
	product, err := uc.productRepo.FindByID(input.ProductID)
	if err != nil {
		return GetProductOutput{ProductExists: false}
	}
	return GetProductOutput{Product: product, ProductExists: true}
}
