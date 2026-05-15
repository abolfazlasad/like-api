package http

import (
	"like-api/internal/application/dto/request"
	"like-api/internal/application/dto/response"
	likeusecase "like-api/internal/application/usecases/like"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LikeHandler struct {
	likePostUseCase          likeusecase.LikePostUseCase
	unlikePostUseCase        likeusecase.UnlikePostUseCase
	getUserLikedPostsUseCase likeusecase.GetUserLikedPostsUseCase
}

func NewLikeHandler(
	likePostUseCase likeusecase.LikePostUseCase,
	unlikePostUseCase likeusecase.UnlikePostUseCase,
	getUserLikedPostsUseCase likeusecase.GetUserLikedPostsUseCase,
) *LikeHandler {
	return &LikeHandler{
		likePostUseCase:          likePostUseCase,
		unlikePostUseCase:        unlikePostUseCase,
		getUserLikedPostsUseCase: getUserLikedPostsUseCase,
	}
}

// @Summary Like a post
// @Description Like a specific post as a user
// @Tags likes
// @Accept json
// @Produce json
// @Param like body request.LikePostRequest true "Like request"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 409 {object} response.Response
// @Router /likes [post]
func (h *LikeHandler) LikePost(c *gin.Context) {
	var req request.LikePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Response{Success: false, Message: "Invalid request: " + err.Error()})
		return
	}

	input := likeusecase.LikePostInput{
		UserID: req.UserID,
		PostID: req.PostID,
	}
	output := h.likePostUseCase.Execute(input)

	if !output.UserExists {
		c.JSON(http.StatusNotFound, response.Response{Success: false, Message: "User not found"})
		return
	}

	if !output.PostExists {
		c.JSON(http.StatusNotFound, response.Response{Success: false, Message: "Post not found"})
		return
	}

	if output.AlreadyLiked {
		c.JSON(http.StatusConflict, response.Response{Success: false, Message: "Post already liked by this user"})
		return
	}

	if output.Error != nil {
		c.JSON(http.StatusInternalServerError, response.Response{Success: false, Message: "Failed to like post"})
		return
	}

	c.JSON(http.StatusOK, response.Response{
		Success: true,
		Message: "Post liked successfully",
		Data:    output.Like,
	})
}

// @Summary Unlike a post
// @Description Remove like from a post
// @Tags likes
// @Accept json
// @Produce json
// @Param unlike body request.UnlikePostRequest true "Unlike request"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /likes [delete]
func (h *LikeHandler) UnlikePost(c *gin.Context) {
	var req request.UnlikePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Response{Success: false, Message: "Invalid request: " + err.Error()})
		return
	}

	input := likeusecase.UnlikePostInput{
		UserID: req.UserID,
		PostID: req.PostID,
	}
	output := h.unlikePostUseCase.Execute(input)

	if !output.LikeExists {
		c.JSON(http.StatusNotFound, response.Response{Success: false, Message: "Like not found"})
		return
	}

	if output.Error != nil {
		c.JSON(http.StatusInternalServerError, response.Response{Success: false, Message: "Failed to unlike post"})
		return
	}

	c.JSON(http.StatusOK, response.Response{
		Success: true,
		Message: "Post unliked successfully",
	})
}

// @Summary Get user's liked posts
// @Description Get all posts liked by a specific user
// @Tags likes
// @Accept json
// @Produce json
// @Param user_id path string true "User ID"
// @Success 200 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /users/{user_id}/liked-posts [get]
func (h *LikeHandler) GetUserLikedPosts(c *gin.Context) {
	userID := c.Param("user_id")

	input := likeusecase.GetUserLikedPostsInput{
		UserID: userID,
	}
	output := h.getUserLikedPostsUseCase.Execute(input)

	if !output.UserExists {
		c.JSON(http.StatusNotFound, response.Response{Success: false, Message: "User not found"})
		return
	}

	if output.Error != nil {
		c.JSON(http.StatusInternalServerError, response.Response{Success: false, Message: "Failed to fetch liked posts"})
		return
	}

	c.JSON(http.StatusOK, response.Response{Success: true, Data: output.Posts})
}
