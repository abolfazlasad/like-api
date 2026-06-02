package http

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"like-api/internal/application/dto/request"
	"like-api/internal/application/dto/response"
	productusecase "like-api/internal/application/usecases/product"
)

// ProductHandler handles product-related HTTP requests.
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

// CreateProduct godoc
//
//	@Summary		Create a product linked to a video
//	@Description	Attaches a shoppable product to a video. Each video can have at most one product.
//	@Description	Returns 409 if a product is already linked to the given video.
//	@Tags			Products
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			body	body		request.CreateProductRequest	true	"Product data"
//	@Success		201		{object}	response.Response{data=entities.Product}
//	@Failure		400		{object}	response.Response	"Validation error"
//	@Failure		401		{object}	response.Response	"Missing or invalid JWT"
//	@Failure		404		{object}	response.Response	"Video not found"
//	@Failure		409		{object}	response.Response	"Product already linked to this video"
//	@Failure		500		{object}	response.Response	"Internal server error"
//	@Router			/api/v1/products [post]
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

// GetProduct godoc
//
//	@Summary		Get a product by ID
//	@Description	Returns a single product by its UUID.
//	@Tags			Products
//	@Produce		json
//	@Param			id	path		string	true	"Product ID"
//	@Success		200	{object}	response.Response{data=entities.Product}
//	@Failure		404	{object}	response.Response	"Product not found"
//	@Router			/api/v1/products/{id} [get]
func (h *ProductHandler) GetProduct(c *gin.Context) {
	productID := c.Param("id")

	output := h.getProductUC.Execute(productusecase.GetProductInput{ProductID: productID})
	if !output.ProductExists {
		c.JSON(http.StatusNotFound, response.Response{Success: false, Message: "Product not found"})
		return
	}

	c.JSON(http.StatusOK, response.Response{Success: true, Data: output.Product})
}

// GetProductByVideo godoc
//
//	@Summary		Get the product attached to a video
//	@Description	Returns the single product linked to the given video, if any.
//	@Tags			Products
//	@Produce		json
//	@Param			id				path		string	true	"Video ID"
//	@Success		200				{object}	response.Response{data=entities.Product}
//	@Failure		404				{object}	response.Response	"Video or product not found"
//	@Router			/api/v1/videos/{id}/product [get]
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
