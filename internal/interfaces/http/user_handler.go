package http

import (
	"like-api/internal/application/dto/response"
	userusecase "like-api/internal/application/usecases/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	getUsersUseCase userusecase.GetUsersUseCase
}

func NewUserHandler(getUsersUseCase userusecase.GetUsersUseCase) *UserHandler {
	return &UserHandler{
		getUsersUseCase: getUsersUseCase,
	}
}

// @Summary Get all users
// @Description Get list of all predefined users
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {object} response.Response
// @Router /api/v1/users [get]
func (h *UserHandler) GetUsers(c *gin.Context) {
	// Execute use case
	input := userusecase.GetUsersInput{}
	output := h.getUsersUseCase.Execute(input)

	if output.Error != nil {
		c.JSON(
			http.StatusInternalServerError,
			response.Response{Success: false, Message: "Failed to fetch users"},
		)
		return
	}

	// Return response
	c.JSON(http.StatusOK, response.Response{Success: true, Data: output.Users})
}
