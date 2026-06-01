package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"like-api/internal/application/dto/request"
	"like-api/internal/application/dto/response"
	videousecase "like-api/internal/application/usecases/video"
)

// VideoHandler handles all video-related HTTP requests.
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

// CreateVideo godoc
//
//	@Summary		Upload a video
//	@Description	Creates a new video/reel owned by the authenticated user.
//	@Tags			Videos
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			body	body		request.CreateVideoRequest	true	"Video metadata"
//	@Success		201		{object}	response.Response{data=entities.Video}
//	@Failure		400		{object}	response.Response	"Validation error"
//	@Failure		401		{object}	response.Response	"Missing or invalid JWT"
//	@Failure		500		{object}	response.Response	"Internal server error"
//	@Router			/api/v1/videos [post]
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

// GetVideo godoc
//
//	@Summary		Get a video by ID
//	@Description	Returns a single video. Accessible by anyone. If a valid JWT is provided the request is attributed to that user.
//	@Tags			Videos
//	@Produce		json
//	@Param			Authorization	header		string	false	"Bearer {token} — optional, enables per-user features"
//	@Param			id				path		string	true	"Video ID"
//	@Success		200				{object}	response.Response{data=entities.Video}
//	@Failure		404				{object}	response.Response	"Video not found"
//	@Router			/api/v1/videos/{id} [get]
func (h *VideoHandler) GetVideo(c *gin.Context) {
	videoID := c.Param("id")

	output := h.getVideoUC.Execute(videousecase.GetVideoInput{VideoID: videoID})
	if !output.VideoExists {
		c.JSON(http.StatusNotFound, response.Response{Success: false, Message: "Video not found"})
		return
	}

	c.JSON(http.StatusOK, response.Response{Success: true, Data: output.Video})
}

// GetFeed godoc
//
//	@Summary		Get paginated video feed
//	@Description	Returns a cursor-paginated list of videos ordered by creation time (newest first).
//	@Description	Authentication is optional — when a valid JWT is supplied the caller is identified
//	@Description	for future personalisation features.
//	@Description
//	@Description	**Cursor pagination:** pass the `next_cursor` value from the previous response as
//	@Description	the `cursor` query parameter to fetch the next page. An empty `next_cursor` in the
//	@Description	response means there are no more pages.
//	@Tags			Feed
//	@Produce		json
//	@Param			Authorization	header		string	false	"Bearer {token} — optional"
//	@Param			cursor			query		string	false	"RFC3339 timestamp cursor from previous page"
//	@Param			limit			query		int		false	"Page size (1–100, default 20)"
//	@Success		200				{object}	response.Response{data=response.FeedResponseData}
//	@Failure		400				{object}	response.Response	"Invalid cursor format"
//	@Router			/api/v1/feed [get]
func (h *VideoHandler) GetFeed(c *gin.Context) {
	cursor := c.Query("cursor")

	limit := 20
	if l := c.Query("limit"); l != "" {
		parsed, err := strconv.Atoi(l)
		if err != nil {
			c.JSON(http.StatusBadRequest, response.Response{
				Success: false,
				Message: "limit must be a number",
			})
			return
		}

		if parsed <= 0 {
			c.JSON(http.StatusBadRequest, response.Response{
				Success: false,
				Message: "limit must be greater than 0",
			})
			return
		}

		if parsed > 100 {
			c.JSON(http.StatusBadRequest, response.Response{
				Success: false,
				Message: "limit must be <= 100",
			})
			return
		}

		limit = parsed
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
		Data: response.FeedResponseData{
			Videos:     output.Videos,
			NextCursor: output.NextCursor,
		},
	})
}

// LikeVideo godoc
//
//	@Summary		Like a video
//	@Description	Adds a like from the authenticated user to the specified video.
//	@Description	Returns 409 if the video is already liked by this user.
//	@Tags			Interactions
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		string	true	"Video ID"
//	@Success		200	{object}	response.Response{data=entities.Like}
//	@Failure		401	{object}	response.Response	"Missing or invalid JWT"
//	@Failure		404	{object}	response.Response	"Video not found"
//	@Failure		409	{object}	response.Response	"Video already liked by this user"
//	@Failure		500	{object}	response.Response	"Internal server error"
//	@Router			/api/v1/videos/{id}/like [post]
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

// UnlikeVideo godoc
//
//	@Summary		Unlike a video
//	@Description	Removes the authenticated user's like from the specified video.
//	@Description	Returns 404 if the user has not liked this video.
//	@Tags			Interactions
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		string	true	"Video ID"
//	@Success		200	{object}	response.Response	"Video unliked successfully"
//	@Failure		401	{object}	response.Response	"Missing or invalid JWT"
//	@Failure		404	{object}	response.Response	"Like not found"
//	@Failure		500	{object}	response.Response	"Internal server error"
//	@Router			/api/v1/videos/{id}/unlike [post]
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

// TrackView godoc
//
//	@Summary		Track a video view
//	@Description	Records a view for the specified video. Authentication is optional.
//	@Description
//	@Description	**Deduplication window: 1 hour.**
//	@Description	- Authenticated: deduplicated per `userID + videoID`.
//	@Description	- Anonymous: deduplicated per `clientIP + videoID`.
//	@Description
//	@Description	The response field `counted` indicates whether this call actually
//	@Description	incremented the view counter (false = duplicate within the window).
//	@Tags			Interactions
//	@Produce		json
//	@Param			Authorization	header		string	false	"Bearer {token} — optional"
//	@Param			id				path		string	true	"Video ID"
//	@Success		200				{object}	response.Response{data=response.TrackViewResponseData}
//	@Failure		404				{object}	response.Response	"Video not found"
//	@Failure		500				{object}	response.Response	"Internal server error"
//	@Router			/api/v1/videos/{id}/view [post]
func (h *VideoHandler) TrackView(c *gin.Context) {
	videoID := c.Param("id")

	// Use userID when authenticated, fall back to client IP for anonymous viewers.
	callerID, _ := c.Get("userID")
	viewerKey, _ := callerID.(string)
	if viewerKey == "" {
		viewerKey = "anon:" + c.ClientIP()
	}

	output := h.trackViewUC.Execute(videousecase.TrackViewInput{
		UserID:  viewerKey,
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
		Data:    response.TrackViewResponseData{Counted: output.Counted},
	})
}

// GetVideoStats godoc
//
//	@Summary		Get video analytics
//	@Description	Returns view count, like count, and engagement rate for the specified video.
//	@Description	Engagement rate = (likes / views) × 100. Returns 0 when views = 0.
//	@Tags			Analytics
//	@Produce		json
//	@Param			Authorization	header		string	false	"Bearer {token} — optional"
//	@Param			id				path		string	true	"Video ID"
//	@Success		200				{object}	response.Response{data=videousecase.VideoStats}
//	@Failure		404				{object}	response.Response	"Video not found"
//	@Router			/api/v1/videos/{id}/stats [get]
func (h *VideoHandler) GetVideoStats(c *gin.Context) {
	videoID := c.Param("id")

	output := h.getVideoStatsUC.Execute(videousecase.GetVideoStatsInput{VideoID: videoID})
	if !output.VideoExists {
		c.JSON(http.StatusNotFound, response.Response{Success: false, Message: "Video not found"})
		return
	}

	c.JSON(http.StatusOK, response.Response{Success: true, Data: output.Stats})
}
