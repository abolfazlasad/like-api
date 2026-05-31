package models

import (
	"like-api/internal/domain/entities"
	"time"
)

// LikeModel is the GORM representation of a "like" record.
// A unique index on (user_id, post_id) prevents duplicate likes
// at the database level as well.
type LikeModel struct {
	ID        string    `gorm:"primaryKey;type:varchar(36)"`
	UserID    string    `gorm:"not null;type:varchar(36);uniqueIndex:idx_user_post"`
	PostID    string    `gorm:"not null;type:varchar(36);uniqueIndex:idx_user_post;index"` // Separate index for FindByPostID queries.
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

func (LikeModel) TableName() string { return "likes" }

func (m LikeModel) ToEntity() entities.Like {
	return entities.Like{
		ID:        m.ID,
		UserID:    m.UserID,
		PostID:    m.PostID,
		CreatedAt: m.CreatedAt.Format(time.RFC3339),
	}
}

func LikeModelFromEntity(l entities.Like) LikeModel {
	var createdAt time.Time
	if l.CreatedAt != "" {
		createdAt, _ = time.Parse(time.RFC3339, l.CreatedAt)
	}
	return LikeModel{
		ID:        l.ID,
		UserID:    l.UserID,
		PostID:    l.PostID,
		CreatedAt: createdAt,
	}
}
