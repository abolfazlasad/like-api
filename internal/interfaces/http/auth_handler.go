package http

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"like-api/internal/application/dto/request"
	"like-api/internal/application/dto/response"
	authusecase "like-api/internal/application/usecases/auth"
)

// AuthHandler handles authentication-related HTTP requests.
type AuthHandler struct {
	registerUC *authusecase.RegisterUseCase
	loginUC    *authusecase.LoginUseCase
}

func NewAuthHandler(register *authusecase.RegisterUseCase, login *authusecase.LoginUseCase) *AuthHandler {
	return &AuthHandler{registerUC: register, loginUC: login}
}

// Register godoc
//
//	@Summary		Register a new user
//	@Description	Creates a new user account with role "user" and returns a signed JWT.
//	@Description	The JWT must be included as a Bearer token in the Authorization header
//	@Description	for all protected endpoints.
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			body	body		request.RegisterRequest		true	"Registration data"
//	@Success		201		{object}	response.Response{data=response.AuthResponse}
//	@Failure		400		{object}	response.Response	"Validation error"
//	@Failure		409		{object}	response.Response	"Email already registered"
//	@Router			/api/v1/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req request.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Response{Success: false, Message: err.Error()})
		return
	}

	resp, err := h.registerUC.Execute(req)
	if err != nil {
		c.JSON(http.StatusConflict, response.Response{Success: false, Message: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, response.Response{Success: true, Data: resp})
}

// Login godoc
//
//	@Summary		Login
//	@Description	Authenticates a user by email and password and returns a signed JWT.
//	@Description	Use the returned token as: `Authorization: Bearer <token>`
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			body	body		request.LoginRequest		true	"Login credentials"
//	@Success		200		{object}	response.Response{data=response.AuthResponse}
//	@Failure		400		{object}	response.Response	"Validation error"
//	@Failure		401		{object}	response.Response	"Invalid credentials"
//	@Router			/api/v1/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req request.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Response{Success: false, Message: err.Error()})
		return
	}

	resp, err := h.loginUC.Execute(req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, response.Response{Success: false, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, response.Response{Success: true, Data: resp})
}
