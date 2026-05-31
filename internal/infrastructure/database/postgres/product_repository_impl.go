package postgres

import (
	"errors"

	"like-api/internal/domain/entities"
	"like-api/internal/domain/repositories"
	"like-api/internal/infrastructure/database/postgres/models"

	"gorm.io/gorm"
)

type productRepositoryImpl struct {
	writeDB *gorm.DB
	readDB  *gorm.DB
}

func newProductRepository(db *DB) repositories.ProductRepository {
	return &productRepositoryImpl{
		writeDB: db.WriteDB,
		readDB:  db.ReadDB,
	}
}

func (r *productRepositoryImpl) Create(product entities.Product) error {
	model := models.ProductModelFromEntity(product)
	return r.writeDB.Create(&model).Error
}

func (r *productRepositoryImpl) FindByID(id string) (entities.Product, error) {
	var model models.ProductModel
	err := r.readDB.Where("id = ?", id).First(&model).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return entities.Product{}, errors.New("product not found")
	}
	if err != nil {
		return entities.Product{}, err
	}
	return model.ToEntity(), nil
}

// FindByVideoID retrieves the single product associated with a video.
// The unique index on video_id guarantees at most one result.
func (r *productRepositoryImpl) FindByVideoID(videoID string) (entities.Product, error) {
	var model models.ProductModel
	err := r.readDB.Where("video_id = ?", videoID).First(&model).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return entities.Product{}, errors.New("product not found")
	}
	if err != nil {
		return entities.Product{}, err
	}
	return model.ToEntity(), nil
}

func (r *productRepositoryImpl) Exists(id string) bool {
	var count int64
	r.readDB.Model(&models.ProductModel{}).Where("id = ?", id).Count(&count)
	return count > 0
}
