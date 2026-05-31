package models

import "like-api/internal/domain/entities"

// ProductModel is the GORM representation of a product.
// A video can have at most one associated product (one-to-one via unique index on video_id).
type ProductModel struct {
	ID       string  `gorm:"primaryKey;type:varchar(36)"`
	VideoID  string  `gorm:"not null;uniqueIndex;type:varchar(36)"`
	Name     string  `gorm:"not null;type:varchar(255)"`
	Price    float64 `gorm:"not null;type:numeric(12,2)"`
	ImageURL string  `gorm:"not null;type:varchar(2048)"`
}

func (ProductModel) TableName() string { return "products" }

func (m ProductModel) ToEntity() entities.Product {
	return entities.Product{
		ID:       m.ID,
		VideoID:  m.VideoID,
		Name:     m.Name,
		Price:    m.Price,
		ImageURL: m.ImageURL,
	}
}

func ProductModelFromEntity(p entities.Product) ProductModel {
	return ProductModel{
		ID:       p.ID,
		VideoID:  p.VideoID,
		Name:     p.Name,
		Price:    p.Price,
		ImageURL: p.ImageURL,
	}
}
