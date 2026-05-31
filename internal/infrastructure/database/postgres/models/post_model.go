package models

import (
	"like-api/internal/domain/entities"
	"time"
)

// PostModel is the GORM representation of a post
type PostModel struct {
	ID         string    `gorm:"primaryKey;type:varchar(36)"`
	UserID     string    `gorm:"not null;index;type:varchar(36)"`
	Title      string    `gorm:"not null;type:varchar(255)"`
	Content    string    `gorm:"type:text"`
	LikesCount int       `gorm:"default:0;not null"`
	CreatedAt  time.Time `gorm:"autoCreateTime;index"` // index for cursor-based pagination
}

func (PostModel) TableName() string { return "posts" }

func (m PostModel) ToEntity() entities.Post {
	return entities.Post{
		ID:      m.ID,
		UserID:  m.UserID,
		Title:   m.Title,
		Content: m.Content,
		Likes:   m.LikesCount,
	}
}

func PostModelFromEntity(p entities.Post) PostModel {
	return PostModel{
		ID:         p.ID,
		UserID:     p.UserID,
		Title:      p.Title,
		Content:    p.Content,
		LikesCount: p.Likes,
	}
}
