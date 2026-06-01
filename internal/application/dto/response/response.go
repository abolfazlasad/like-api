package response

import "like-api/internal/domain/entities"

// Response represents a standard API response
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

type FeedResponseData struct {
	Videos     []entities.Video
	NextCursor string
}

type TrackViewResponseData struct {
	Counted bool
}
