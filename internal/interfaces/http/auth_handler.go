package http

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"like-api/internal/application/dto/request"
	"like-api/internal/application/dto/response"
	authusecase "like-api/internal/application/usecases/auth"
)

type AuthHandler struct {
	registerUC *authusecase.RegisterUseCase
	loginUC    *authusecase.LoginUseCase
}

func NewAuthHandler(register *authusecase.RegisterUseCase, login *authusecase.LoginUseCase) *AuthHandler {
	return &AuthHandler{registerUC: register, loginUC: login}
}

// @Summary Register a new user
// @Tags auth
// @Accept json
// @Produce json
// @Param body body request.RegisterRequest true "Register request"
// @Success 201 {object} response.AuthResponse
// @Failure 400 {object} response.Response
// @Failure 409 {object} response.Response
// @Router /api/v1/auth/register [post]
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

// @Summary Login
// @Tags auth
// @Accept json
// @Produce json
// @Param body body request.LoginRequest true "Login request"
// @Success 200 {object} response.AuthResponse
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /api/v1/auth/login [post]
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
