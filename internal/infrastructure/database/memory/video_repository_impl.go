package memory

import (
	"errors"
	"like-api/internal/domain/entities"
	"like-api/internal/domain/repositories"
	"sort"
	"sync"
	"time"
)

type videoRepositoryImpl struct {
	videos map[string]entities.Video
	mu     sync.RWMutex
}

func NewVideoRepository() repositories.VideoRepository {
	return &videoRepositoryImpl{
		videos: make(map[string]entities.Video),
	}
}

func (r *videoRepositoryImpl) Create(video entities.Video) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.videos[video.ID] = video
	return nil
}

func (r *videoRepositoryImpl) FindByID(id string) (entities.Video, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	video, exists := r.videos[id]
	if !exists {
		return entities.Video{}, errors.New("video not found")
	}
	return video, nil
}

// FindAll implements cursor-based pagination over an in-memory map.
// Videos are sorted by created_at DESC; the cursor is an RFC3339Nano timestamp
// marking the exclusive upper bound (same contract as the postgres impl).
func (r *videoRepositoryImpl) FindAll(cursor string, limit int) ([]entities.Video, string, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}

	r.mu.RLock()
	all := make([]entities.Video, 0, len(r.videos))
	for _, v := range r.videos {
		all = append(all, v)
	}
	r.mu.RUnlock()

	sort.Slice(all, func(i, j int) bool {
		return all[i].CreatedAt > all[j].CreatedAt
	})

	if cursor != "" {
		cursorTime, err := time.Parse(time.RFC3339Nano, cursor)
		if err != nil {
			return nil, "", errors.New("invalid cursor format")
		}
		filtered := all[:0]
		for _, v := range all {
			t, _ := time.Parse(time.RFC3339Nano, v.CreatedAt)
			if t.Before(cursorTime) {
				filtered = append(filtered, v)
			}
		}
		all = filtered
	}

	hasMore := len(all) > limit
	if hasMore {
		all = all[:limit]
	}

	nextCursor := ""
	if hasMore && len(all) > 0 {
		nextCursor = all[len(all)-1].CreatedAt
	}

	return all, nextCursor, nil
}

func (r *videoRepositoryImpl) FindByUserID(userID string) ([]entities.Video, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	videos := make([]entities.Video, 0)
	for _, v := range r.videos {
		if v.UserID == userID {
			videos = append(videos, v)
		}
	}
	return videos, nil
}

func (r *videoRepositoryImpl) Exists(id string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	_, exists := r.videos[id]
	return exists
}

func (r *videoRepositoryImpl) IncrementLikes(videoID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	video, exists := r.videos[videoID]
	if !exists {
		return errors.New("video not found")
	}
	video.LikesCount++
	r.videos[videoID] = video
	return nil
}

func (r *videoRepositoryImpl) DecrementLikes(videoID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	video, exists := r.videos[videoID]
	if !exists {
		return errors.New("video not found")
	}
	if video.LikesCount > 0 {
		video.LikesCount--
	}
	r.videos[videoID] = video
	return nil
}

func (r *videoRepositoryImpl) IncrementViews(videoID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	video, exists := r.videos[videoID]
	if !exists {
		return errors.New("video not found")
	}
	video.ViewsCount++
	r.videos[videoID] = video
	return nil
}
