package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"like-api/internal/application/dto/request"
	"like-api/internal/application/dto/response"
	videousecase "like-api/internal/application/usecases/video"
)

type VideoHandler struct {
	createVideoUC   videousecase.CreateVideoUseCase
	getVideoUC      videousecase.GetVideoUseCase
	getFeedUC       videousecase.GetFeedUseCase
	likeVideoUC     videousecase.LikeVideoUseCase
	unlikeVideoUC   videousecase.UnlikeVideoUseCase
	trackViewUC     videousecase.TrackViewUseCase
	getVideoStatsUC videousecase.GetVideoStatsUseCase
}

func NewVideoHandler(
	createVideoUC videousecase.CreateVideoUseCase,
	getVideoUC videousecase.GetVideoUseCase,
	getFeedUC videousecase.GetFeedUseCase,
	likeVideoUC videousecase.LikeVideoUseCase,
	unlikeVideoUC videousecase.UnlikeVideoUseCase,
	trackViewUC videousecase.TrackViewUseCase,
	getVideoStatsUC videousecase.GetVideoStatsUseCase,
) *VideoHandler {
	return &VideoHandler{
		createVideoUC:   createVideoUC,
		getVideoUC:      getVideoUC,
		getFeedUC:       getFeedUC,
		likeVideoUC:     likeVideoUC,
		unlikeVideoUC:   unlikeVideoUC,
		trackViewUC:     trackViewUC,
		getVideoStatsUC: getVideoStatsUC,
	}
}

// @Summary Upload a video
// @Tags videos
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body request.CreateVideoRequest true "Video data"
// @Success 201 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /api/v1/videos [post]
func (h *VideoHandler) CreateVideo(c *gin.Context) {
	var req request.CreateVideoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Response{Success: false, Message: err.Error()})
		return
	}

	userID, _ := c.Get("userID")

	output := h.createVideoUC.Execute(videousecase.CreateVideoInput{
		UserID:      userID.(string),
		Title:       req.Title,
		Description: req.Description,
		VideoURL:    req.VideoURL,
	})

	if output.Error != nil {
		c.JSON(http.StatusInternalServerError, response.Response{Success: false, Message: "Failed to create video"})
		return
	}

	c.JSON(http.StatusCreated, response.Response{Success: true, Data: output.Video})
}

// @Summary Get a video by ID
// @Tags videos
// @Produce json
// @Param id path string true "Video ID"
// @Success 200 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /api/v1/videos/{id} [get]
func (h *VideoHandler) GetVideo(c *gin.Context) {
	videoID := c.Param("id")

	output := h.getVideoUC.Execute(videousecase.GetVideoInput{VideoID: videoID})
	if !output.VideoExists {
		c.JSON(http.StatusNotFound, response.Response{Success: false, Message: "Video not found"})
		return
	}

	c.JSON(http.StatusOK, response.Response{Success: true, Data: output.Video})
}

// @Summary Get paginated video feed
// @Tags videos
// @Produce json
// @Param cursor query string false "Pagination cursor (RFC3339 timestamp)"
// @Param limit  query int    false "Page size (default 20, max 100)"
// @Success 200 {object} response.Response
// @Router /api/v1/feed [get]
func (h *VideoHandler) GetFeed(c *gin.Context) {
	cursor := c.Query("cursor")

	limit := 20
	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	output := h.getFeedUC.Execute(videousecase.GetFeedInput{
		Cursor: cursor,
		Limit:  limit,
	})

	if output.Error != nil {
		c.JSON(http.StatusBadRequest, response.Response{Success: false, Message: output.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, response.Response{
		Success: true,
		Data: gin.H{
			"videos":      output.Videos,
			"next_cursor": output.NextCursor,
		},
	})
}

// @Summary Like a video
// @Tags videos
// @Security BearerAuth
// @Param id path string true "Video ID"
// @Success 200 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 409 {object} response.Response
// @Router /api/v1/videos/{id}/like [post]
func (h *VideoHandler) LikeVideo(c *gin.Context) {
	videoID := c.Param("id")
	userID, _ := c.Get("userID")

	output := h.likeVideoUC.Execute(videousecase.LikeVideoInput{
		UserID:  userID.(string),
		VideoID: videoID,
	})

	if !output.UserExists {
		c.JSON(http.StatusNotFound, response.Response{Success: false, Message: "User not found"})
		return
	}
	if !output.VideoExists {
		c.JSON(http.StatusNotFound, response.Response{Success: false, Message: "Video not found"})
		return
	}
	if output.AlreadyLiked {
		c.JSON(http.StatusConflict, response.Response{Success: false, Message: "Video already liked"})
		return
	}
	if output.Error != nil {
		c.JSON(http.StatusInternalServerError, response.Response{Success: false, Message: "Failed to like video"})
		return
	}

	c.JSON(http.StatusOK, response.Response{Success: true, Message: "Video liked successfully", Data: output.Like})
}

// @Summary Unlike a video
// @Tags videos
// @Security BearerAuth
// @Param id path string true "Video ID"
// @Success 200 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /api/v1/videos/{id}/unlike [post]
func (h *VideoHandler) UnlikeVideo(c *gin.Context) {
	videoID := c.Param("id")
	userID, _ := c.Get("userID")

	output := h.unlikeVideoUC.Execute(videousecase.UnlikeVideoInput{
		UserID:  userID.(string),
		VideoID: videoID,
	})

	if !output.LikeExists {
		c.JSON(http.StatusNotFound, response.Response{Success: false, Message: "Like not found"})
		return
	}
	if output.Error != nil {
		c.JSON(http.StatusInternalServerError, response.Response{Success: false, Message: "Failed to unlike video"})
		return
	}

	c.JSON(http.StatusOK, response.Response{Success: true, Message: "Video unliked successfully"})
}

// @Summary Track a video view
// @Tags videos
// @Security BearerAuth
// @Param id path string true "Video ID"
// @Success 200 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /api/v1/videos/{id}/view [post]
func (h *VideoHandler) TrackView(c *gin.Context) {
	videoID := c.Param("id")
	userID, _ := c.Get("userID")

	output := h.trackViewUC.Execute(videousecase.TrackViewInput{
		UserID:  userID.(string),
		VideoID: videoID,
	})

	if !output.VideoExists {
		c.JSON(http.StatusNotFound, response.Response{Success: false, Message: "Video not found"})
		return
	}
	if output.Error != nil {
		c.JSON(http.StatusInternalServerError, response.Response{Success: false, Message: "Failed to track view"})
		return
	}

	c.JSON(http.StatusOK, response.Response{
		Success: true,
		Data:    gin.H{"counted": output.Counted},
	})
}

// @Summary Get video stats (bonus)
// @Tags videos
// @Produce json
// @Param id path string true "Video ID"
// @Success 200 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /api/v1/videos/{id}/stats [get]
func (h *VideoHandler) GetVideoStats(c *gin.Context) {
	videoID := c.Param("id")

	output := h.getVideoStatsUC.Execute(videousecase.GetVideoStatsInput{VideoID: videoID})
	if !output.VideoExists {
		c.JSON(http.StatusNotFound, response.Response{Success: false, Message: "Video not found"})
		return
	}

	c.JSON(http.StatusOK, response.Response{Success: true, Data: output.Stats})
}
