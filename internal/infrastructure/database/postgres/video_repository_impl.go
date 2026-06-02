package postgres

import (
	"errors"
	"time"

	"like-api/internal/domain/entities"
	"like-api/internal/domain/repositories"
	"like-api/internal/infrastructure/database/postgres/models"

	"gorm.io/gorm"
)

type videoRepositoryImpl struct {
	writeDB *gorm.DB
	readDB  *gorm.DB
}

func newVideoRepository(db *DB) repositories.VideoRepository {
	return &videoRepositoryImpl{
		writeDB: db.WriteDB,
		readDB:  db.ReadDB,
	}
}

func (r *videoRepositoryImpl) Create(video entities.Video) error {
	model := models.VideoModelFromEntity(video)
	return r.writeDB.Create(&model).Error
}

func (r *videoRepositoryImpl) FindByID(id string) (entities.Video, error) {
	var model models.VideoModel
	err := r.readDB.Where("id = ?", id).First(&model).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return entities.Video{}, errors.New("video not found")
	}
	if err != nil {
		return entities.Video{}, err
	}
	return model.ToEntity(), nil
}

// FindAll implements cursor-based pagination.
//
// The cursor encodes the created_at timestamp of the last seen video.
// On each page we fetch rows where created_at < cursor (exclusive), ordered
// by created_at DESC so the newest items appear first.  This is more
// efficient than OFFSET-based pagination because:
//   - It never re-scans skipped rows.
//   - The index on created_at makes the WHERE clause a range scan.
//   - It stays stable even when new videos are inserted between pages.
//
// The next cursor returned to the caller is the created_at value of the
// last item in the current page, encoded as RFC3339Nano.  An empty string
// means there are no more pages.
func (r *videoRepositoryImpl) FindAll(cursor string, limit int) ([]entities.Video, string, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}

	query := r.readDB.Order("created_at DESC").Limit(limit + 1)

	if cursor != "" {
		cursorTime, err := time.Parse(time.RFC3339Nano, cursor)
		if err != nil {
			return nil, "", errors.New("invalid cursor format")
		}
		query = query.Where("created_at < ?", cursorTime)
	}

	var videoModels []models.VideoModel
	if err := query.Find(&videoModels).Error; err != nil {
		return nil, "", err
	}

	hasMore := len(videoModels) > limit
	if hasMore {
		videoModels = videoModels[:limit]
	}

	videos := make([]entities.Video, len(videoModels))
	for i, m := range videoModels {
		videos[i] = m.ToEntity()
	}

	nextCursor := ""
	if hasMore && len(videoModels) > 0 {
		lastCreatedAt, _ := time.Parse(time.RFC3339Nano, videos[len(videos)-1].CreatedAt)
		nextCursor = lastCreatedAt.Format(time.RFC3339Nano)
	}

	return videos, nextCursor, nil
}

func (r *videoRepositoryImpl) FindByUserID(userID string) ([]entities.Video, error) {
	var videoModels []models.VideoModel
	if err := r.readDB.Where("user_id = ?", userID).Order("created_at DESC").Find(&videoModels).Error; err != nil {
		return nil, err
	}
	videos := make([]entities.Video, len(videoModels))
	for i, m := range videoModels {
		videos[i] = m.ToEntity()
	}
	return videos, nil
}

func (r *videoRepositoryImpl) Exists(id string) bool {
	var count int64
	r.readDB.Model(&models.VideoModel{}).Where("id = ?", id).Count(&count)
	return count > 0
}

func (r *videoRepositoryImpl) IncrementLikes(videoID string) error {
	result := r.writeDB.Model(&models.VideoModel{}).
		Where("id = ?", videoID).
		Update("likes_count", gorm.Expr("likes_count + 1"))
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("video not found")
	}
	return nil
}

func (r *videoRepositoryImpl) DecrementLikes(videoID string) error {
	result := r.writeDB.Model(&models.VideoModel{}).
		Where("id = ?", videoID).
		Update("likes_count", gorm.Expr("GREATEST(likes_count - 1, 0)"))
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("video not found")
	}
	return nil
}

// IncrementViews atomically increments the view counter.
//
// Trade-off note: writing a view synchronously on every request is simple
// but becomes a bottleneck under high traffic.  A production system would
// buffer view events in Redis and flush them via a background worker
// (see the Bonus section in the project requirements).  The atomic SQL
// update here at least prevents lost-update races at the DB level.
func (r *videoRepositoryImpl) IncrementViews(videoID string) error {
	result := r.writeDB.Model(&models.VideoModel{}).
		Where("id = ?", videoID).
		Update("views_count", gorm.Expr("views_count + 1"))
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("video not found")
	}
	return nil
}
