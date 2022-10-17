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
func (pc *ProductController) ViewAll(c *gin.Context) {
	var products []models.Product
	result := pc.DB.Find(&products)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "No products!"})
		return
	}
	c.JSON(http.StatusOK, products)
}

// Find a product
func (pc *ProductController) ViewById(c *gin.Context) {
	var product models.Product
	result := pc.DB.First(&product, "id = ?", c.Param("id"))

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Product not found!"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": product})
}

// Update a product
func (pc *ProductController) EditById(c *gin.Context) {
	var payload *UpdateProductInput
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	var updatedProduct models.Product
	result := pc.DB.First(&updatedProduct, "id = ?", c.Param("id"))
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Product not found!"})
		return
	}

	productToUpdate := models.Product{
		Name:      payload.Name,
		Category:  payload.Category,
		ImageURL:  payload.ImageURL,
		Price:     payload.Price,
		UpdatedAt: time.Now(),
	}

	pc.DB.Model(&updatedProduct).Updates(productToUpdate)

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": updatedProduct})

}

func (pc *ProductController) Delete(c *gin.Context) {
	result := pc.DB.Delete(&models.Product{}, "id = ?", c.Param("id"))

	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "No product with the given id exists"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
