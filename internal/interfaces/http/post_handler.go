package http

import (
	"like-api/internal/application/dto/response"
	postusecase "like-api/internal/application/usecases/post"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PostHandler struct {
	getPostsUseCase     postusecase.GetPostsUseCase
	getUserPostsUseCase postusecase.GetUserPostsUseCase
	getPostLikesUseCase postusecase.GetPostLikesUseCase
}

func NewPostHandler(
	getPostsUseCase postusecase.GetPostsUseCase,
	getUserPostsUseCase postusecase.GetUserPostsUseCase,
	getPostLikesUseCase postusecase.GetPostLikesUseCase,
) *PostHandler {
	return &PostHandler{
		getPostsUseCase:     getPostsUseCase,
		getUserPostsUseCase: getUserPostsUseCase,
		getPostLikesUseCase: getPostLikesUseCase,
	}
}

// @Summary Get all posts
// @Description Get list of all posts
// @Tags posts
// @Accept json
// @Produce json
// @Success 200 {object} response.Response
// @Router /posts [get]
func (h *PostHandler) GetPosts(c *gin.Context) {
	// Execute use case
	input := postusecase.GetPostsInput{}
	output := h.getPostsUseCase.Execute(input)

	// Handle error
	if output.Error != nil {
		c.JSON(http.StatusInternalServerError, response.Response{Success: false, Message: "Failed to fetch posts"})
		return
	}

	// Return response
	c.JSON(http.StatusOK, response.Response{Success: true, Data: output.Posts})
}

// @Summary Get posts by user
// @Description Get all posts from a specific user
// @Tags posts
// @Accept json
// @Produce json
// @Param user_id path string true "User ID"
// @Success 200 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /users/{user_id}/posts [get]
func (h *PostHandler) GetUserPosts(c *gin.Context) {
	userID := c.Param("user_id")

	input := postusecase.GetUserPostsInput{
		UserID: userID,
	}
	output := h.getUserPostsUseCase.Execute(input)

	if !output.UserExists {
		c.JSON(http.StatusNotFound, response.Response{Success: false, Message: "User not found"})
		return
	}

	if output.Error != nil {
		c.JSON(http.StatusInternalServerError, response.Response{Success: false, Message: "Failed to fetch posts"})
		return
	}

	c.JSON(http.StatusOK, response.Response{Success: true, Data: output.Posts})
}

// @Summary Get likes for a post
// @Description Get all likes for a specific post
// @Tags likes
// @Accept json
// @Produce json
// @Param post_id path string true "Post ID"
// @Success 200 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /posts/{post_id}/likes [get]
func (h *PostHandler) GetPostLikes(c *gin.Context) {
	postID := c.Param("post_id")

	input := postusecase.GetPostLikesInput{
		PostID: postID,
	}
	output := h.getPostLikesUseCase.Execute(input)

	if !output.PostExists {
		c.JSON(http.StatusNotFound, response.Response{Success: false, Message: "Post not found"})
		return
	}

	if output.Error != nil {
		c.JSON(http.StatusInternalServerError, response.Response{Success: false, Message: "Failed to fetch likes"})
		return
	}

	c.JSON(http.StatusOK, response.Response{Success: true, Data: output.Likes})
}
