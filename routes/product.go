package routes

import (
	"blitzomni.com/m/controllers"
	"github.com/gin-gonic/gin"
)

type ProductRoute struct {
	productController controllers.ProductController
}

func NewProductRoute(productController controllers.ProductController) ProductRoute {
	return ProductRoute{productController}
}

func (pr *ProductRoute) ProductRoute(rg *gin.RouterGroup) {
	router := rg.Group("/products")

	router.POST("/", pr.productController.Create)
	router.GET("/", pr.productController.ViewAll)
	router.GET("/:id", pr.productController.ViewById)
	router.PUT("/:id", pr.productController.EditById)
	router.DELETE("/:id", pr.productController.Delete)
}
