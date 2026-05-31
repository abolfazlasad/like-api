package videousecase

import (
	"like-api/internal/domain/repositories"
	"sync"
	"time"
)

// viewRecord tracks the last view time for a (userID, videoID) pair.
type viewRecord struct {
	lastSeen time.Time
}

// viewTracker is an in-process deduplication window.
//
// Trade-off: this is intentionally simple — a production system would use a
// Redis SET with a TTL (e.g. SETEX "view:{userID}:{videoID}" 3600 1) so the
// window survives restarts and works across multiple instances.  The in-memory
// map is fine for a single-instance deployment or a take-home challenge.
type viewTracker struct {
	mu      sync.Mutex
	records map[string]viewRecord
	window  time.Duration
}

var globalViewTracker = &viewTracker{
	records: make(map[string]viewRecord),
	window:  time.Hour,
}

func (t *viewTracker) shouldCount(userID, videoID string) bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	key := userID + ":" + videoID
	rec, ok := t.records[key]
	now := time.Now()

	if ok && now.Sub(rec.lastSeen) < t.window {
		return false
	}

	t.records[key] = viewRecord{lastSeen: now}
	return true
}

// ---

type TrackViewInput struct {
	UserID  string
	VideoID string
}

type TrackViewOutput struct {
	Counted     bool
	VideoExists bool
	Error       error
}

type TrackViewUseCase interface {
	Execute(input TrackViewInput) TrackViewOutput
}

type trackViewUseCaseImpl struct {
	videoRepo   repositories.VideoRepository
	viewTracker *viewTracker
}

func NewTrackViewUseCase(videoRepo repositories.VideoRepository) TrackViewUseCase {
	return &trackViewUseCaseImpl{
		videoRepo:   videoRepo,
		viewTracker: globalViewTracker,
	}
}

func (uc *trackViewUseCaseImpl) Execute(input TrackViewInput) TrackViewOutput {
	if !uc.videoRepo.Exists(input.VideoID) {
		return TrackViewOutput{VideoExists: false}
	}

	if !uc.viewTracker.shouldCount(input.UserID, input.VideoID) {
		return TrackViewOutput{VideoExists: true, Counted: false}
	}

	if err := uc.videoRepo.IncrementViews(input.VideoID); err != nil {
		return TrackViewOutput{VideoExists: true, Error: err}
	}

	return TrackViewOutput{VideoExists: true, Counted: true}
}
