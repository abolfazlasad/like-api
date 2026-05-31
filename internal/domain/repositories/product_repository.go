package repositories

import "like-api/internal/domain/entities"

// ProductRepository defines the data-access contract for products.
type ProductRepository interface {
	Create(product entities.Product) error
	FindByID(id string) (entities.Product, error)
	FindByVideoID(videoID string) (entities.Product, error)
	Exists(id string) bool
}
