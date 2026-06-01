package http

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"like-api/internal/application/dto/response"
	userusecase "like-api/internal/application/usecases/user"
)

// UserHandler handles user-related HTTP requests.
type UserHandler struct {
	getUsersUseCase userusecase.GetUsersUseCase
}

func NewUserHandler(getUsersUseCase userusecase.GetUsersUseCase) *UserHandler {
	return &UserHandler{
		getUsersUseCase: getUsersUseCase,
	}
}

// GetUsers godoc
//
//	@Summary		List all users
//	@Description	Returns a list of all registered users. Requires admin role.
//	@Tags			Admin
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	response.Response{data=[]response.UserDTO}
//	@Failure		401	{object}	response.Response	"Missing or invalid JWT"
//	@Failure		403	{object}	response.Response	"Admin access required"
//	@Failure		500	{object}	response.Response	"Internal server error"
//	@Router			/api/v1/admin/users [get]
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

	users := make([]response.UserDTO, len(output.Users))
	for i, u := range output.Users {
		users[i] = response.UserDTO{
			ID:       u.ID,
			Username: u.Username,
			Name:     u.Name,
			Email:    u.Email,
			Role:     string(u.Role),
		}
	}

	c.JSON(http.StatusOK, response.Response{Success: true, Data: users})
}
