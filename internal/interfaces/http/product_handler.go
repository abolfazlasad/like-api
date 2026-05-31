package http

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"like-api/internal/application/dto/request"
	"like-api/internal/application/dto/response"
	productusecase "like-api/internal/application/usecases/product"
)

type ProductHandler struct {
	createProductUC     productusecase.CreateProductUseCase
	getProductUC        productusecase.GetProductUseCase
	getProductByVideoUC productusecase.GetProductByVideoUseCase
}

func NewProductHandler(
	createProductUC productusecase.CreateProductUseCase,
	getProductUC productusecase.GetProductUseCase,
	getProductByVideoUC productusecase.GetProductByVideoUseCase,
) *ProductHandler {
	return &ProductHandler{
		createProductUC:     createProductUC,
		getProductUC:        getProductUC,
		getProductByVideoUC: getProductByVideoUC,
	}
}

// @Summary Create a product linked to a video
// @Tags products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body request.CreateProductRequest true "Product data"
// @Success 201 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 409 {object} response.Response
// @Router /api/v1/products [post]
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var req request.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Response{Success: false, Message: err.Error()})
		return
	}

	output := h.createProductUC.Execute(productusecase.CreateProductInput{
		VideoID:  req.VideoID,
		Name:     req.Name,
		Price:    req.Price,
		ImageURL: req.ImageURL,
	})

	if !output.VideoExists {
		c.JSON(http.StatusNotFound, response.Response{Success: false, Message: "Video not found"})
		return
	}
	if output.Error != nil {
		c.JSON(http.StatusConflict, response.Response{Success: false, Message: output.Error.Error()})
		return
	}

	c.JSON(http.StatusCreated, response.Response{Success: true, Data: output.Product})
}

// @Summary Get a product by ID
// @Tags products
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /api/v1/products/{id} [get]
func (h *ProductHandler) GetProduct(c *gin.Context) {
	productID := c.Param("id")

	output := h.getProductUC.Execute(productusecase.GetProductInput{ProductID: productID})
	if !output.ProductExists {
		c.JSON(http.StatusNotFound, response.Response{Success: false, Message: "Product not found"})
		return
	}

	c.JSON(http.StatusOK, response.Response{Success: true, Data: output.Product})
}

// @Summary Get the product attached to a video
// @Tags products
// @Produce json
// @Param id path string true "Video ID"
// @Success 200 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /api/v1/videos/{id}/product [get]
func (h *ProductHandler) GetProductByVideo(c *gin.Context) {
	videoID := c.Param("id")

	output := h.getProductByVideoUC.Execute(productusecase.GetProductByVideoInput{VideoID: videoID})

	if !output.VideoExists {
		c.JSON(http.StatusNotFound, response.Response{Success: false, Message: "Video not found"})
		return
	}
	if !output.ProductExists {
		c.JSON(http.StatusNotFound, response.Response{Success: false, Message: "No product linked to this video"})
		return
	}

	c.JSON(http.StatusOK, response.Response{Success: true, Data: output.Product})
}
