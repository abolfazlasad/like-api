package postgres

import (
	"errors"

	"like-api/internal/domain/entities"
	"like-api/internal/domain/repositories"
	"like-api/internal/infrastructure/database/postgres/models"

	"gorm.io/gorm"
)

// postRepositoryImpl implements the PostRepository interface using PostgreSQL and GORM.
type postRepositoryImpl struct {
	writeDB *gorm.DB
	readDB  *gorm.DB
}

// newPostRepository creates a new post repository instance.
func newPostRepository(db *DB) repositories.PostRepository {
	return &postRepositoryImpl{
		writeDB: db.WriteDB,
		readDB:  db.ReadDB,
	}
}

// Create persists a new post.
func (r *postRepositoryImpl) Create(post entities.Post) error {
	model := models.PostModelFromEntity(post)
	return r.writeDB.Create(&model).Error
}

// FindAll retrieves all posts.
func (r *postRepositoryImpl) FindAll() ([]entities.Post, error) {
	var postModels []models.PostModel

	if err := r.readDB.Find(&postModels).Error; err != nil {
		return nil, err
	}

	posts := make([]entities.Post, len(postModels))
	for i, m := range postModels {
		posts[i] = m.ToEntity()
	}

	return posts, nil
}

// FindByID retrieves a post by its ID.
func (r *postRepositoryImpl) FindByID(id string) (entities.Post, error) {
	var model models.PostModel

	err := r.readDB.Where("id = ?", id).First(&model).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return entities.Post{}, errors.New("post not found")
	}
	if err != nil {
		return entities.Post{}, err
	}

	return model.ToEntity(), nil
}

// FindByUserID retrieves all posts created by a specific user.
func (r *postRepositoryImpl) FindByUserID(userID string) ([]entities.Post, error) {
	var postModels []models.PostModel

	if err := r.readDB.Where("user_id = ?", userID).Find(&postModels).Error; err != nil {
		return nil, err
	}

	posts := make([]entities.Post, len(postModels))
	for i, m := range postModels {
		posts[i] = m.ToEntity()
	}

	return posts, nil
}

// Update modifies non-atomic fields of a post.
// For likes_count updates, use IncrementLikes or DecrementLikes instead.
func (r *postRepositoryImpl) Update(post entities.Post) error {
	return r.writeDB.Model(&models.PostModel{}).
		Where("id = ?", post.ID).
		Updates(map[string]interface{}{
			"title":   post.Title,
			"content": post.Content,
		}).Error
}

// IncrementLikes atomically increments the likes_count field.
// This operation is safe from race conditions.
func (r *postRepositoryImpl) IncrementLikes(postID string) error {
	result := r.writeDB.Model(&models.PostModel{}).
		Where("id = ?", postID).
		Update("likes_count", gorm.Expr("likes_count + 1"))

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("post not found")
	}

	return nil
}

// DecrementLikes atomically decrements likes_count while preventing negative values.
func (r *postRepositoryImpl) DecrementLikes(postID string) error {
	result := r.writeDB.Model(&models.PostModel{}).
		Where("id = ?", postID).
		Update("likes_count", gorm.Expr("GREATEST(likes_count - 1, 0)"))

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("post not found")
	}

	return nil
}

// Exists checks whether a post with the given ID exists.
func (r *postRepositoryImpl) Exists(id string) bool {
	var count int64

	r.readDB.Model(&models.PostModel{}).
		Where("id = ?", id).
		Count(&count)

	return count > 0
}
