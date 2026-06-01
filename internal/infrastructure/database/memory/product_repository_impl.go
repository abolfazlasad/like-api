package memory

import (
	"errors"
	"like-api/internal/domain/entities"
	"like-api/internal/domain/repositories"
	"sync"
)

type productRepositoryImpl struct {
	products map[string]entities.Product
	mu       sync.RWMutex
}

func NewProductRepository() repositories.ProductRepository {
	return &productRepositoryImpl{
		products: make(map[string]entities.Product),
	}
}

func (r *productRepositoryImpl) Create(product entities.Product) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.products[product.ID] = product
	return nil
}

func (r *productRepositoryImpl) FindByID(id string) (entities.Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	product, exists := r.products[id]
	if !exists {
		return entities.Product{}, errors.New("product not found")
	}
	return product, nil
}

func (r *productRepositoryImpl) FindByVideoID(videoID string) (entities.Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, p := range r.products {
		if p.VideoID == videoID {
			return p, nil
		}
	}
	return entities.Product{}, errors.New("product not found")
}

func (r *productRepositoryImpl) Exists(id string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	_, exists := r.products[id]
	return exists
}
