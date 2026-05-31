package postgres

import (
	"errors"

	"like-api/internal/domain/entities"
	"like-api/internal/domain/repositories"
	"like-api/internal/infrastructure/database/postgres/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// likeRepositoryImpl implements the LikeRepository interface using PostgreSQL and GORM.
type likeRepositoryImpl struct {
	writeDB *gorm.DB
	readDB  *gorm.DB
}

// newLikeRepository creates a new like repository instance.
func newLikeRepository(db *DB) repositories.LikeRepository {
	return &likeRepositoryImpl{
		writeDB: db.WriteDB,
		readDB:  db.ReadDB,
	}
}

// Create persists a new like record.
//
// The operation uses ON CONFLICT DO NOTHING to safely handle duplicate likes.
// A unique index on (user_id, post_id) enforces uniqueness at the database level.
func (r *likeRepositoryImpl) Create(like entities.Like) error {
	model := models.LikeModelFromEntity(like)

	return r.writeDB.
		Clauses(clause.OnConflict{DoNothing: true}).
		Create(&model).
		Error
}

// Delete removes a user's like from a specific post.
func (r *likeRepositoryImpl) Delete(userID, postID string) error {
	result := r.writeDB.
		Where("user_id = ? AND post_id = ?", userID, postID).
		Delete(&models.LikeModel{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("like not found")
	}

	return nil
}

// FindByPostID retrieves all likes associated with a post.
func (r *likeRepositoryImpl) FindByPostID(postID string) ([]entities.Like, error) {
	var likeModels []models.LikeModel

	if err := r.readDB.
		Where("post_id = ?", postID).
		Find(&likeModels).
		Error; err != nil {
		return nil, err
	}

	likes := make([]entities.Like, len(likeModels))
	for i, m := range likeModels {
		likes[i] = m.ToEntity()
	}

	return likes, nil
}

// FindByUserID retrieves all likes created by a specific user.
func (r *likeRepositoryImpl) FindByUserID(userID string) ([]entities.Like, error) {
	var likeModels []models.LikeModel

	if err := r.readDB.
		Where("user_id = ?", userID).
		Find(&likeModels).
		Error; err != nil {
		return nil, err
	}

	likes := make([]entities.Like, len(likeModels))
	for i, m := range likeModels {
		likes[i] = m.ToEntity()
	}

	return likes, nil
}

// Exists checks whether a like already exists for the given user and post.
func (r *likeRepositoryImpl) Exists(userID, postID string) bool {
	var count int64

	r.readDB.Model(&models.LikeModel{}).
		Where("user_id = ? AND post_id = ?", userID, postID).
		Count(&count)

	return count > 0
}
