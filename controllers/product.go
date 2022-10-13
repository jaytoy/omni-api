package controllers

import (
	"net/http"
	"time"

	"blitzomni.com/m/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CreateProductInput struct {
	Name     string  `json:"product_name" binding:"required"`
	Category string  `json:"category"`
	ImageURL string  `json:"image"`
	Price    float64 `json:"price"`
}

type UpdateProductInput struct {
	Name     string  `json:"product_name"`
	Category string  `json:"category"`
	ImageURL string  `json:"image"`
	Price    float64 `json:"price"`
}

type ProductController struct {
	DB *gorm.DB
}

func NewProductController(DB *gorm.DB) ProductController {
	return ProductController{DB}
}

// Create a new product
func (pc *ProductController) Create(c *gin.Context) {
	var payload CreateProductInput
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
	}

	newProduct := models.Product{
		Name:      payload.Name,
		Category:  payload.Category,
		ImageURL:  payload.ImageURL,
		Price:     payload.Price,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	result := pc.DB.Create(&newProduct)

	if result.Error != nil {
		c.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": "Something bad happened"})
		return
	}
}

// Find all products
func ViewAll() {

}

// Find a product
func ViewById() {

}

func Edit() {

}

func Delete() {

}
