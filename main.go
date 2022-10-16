package main

import (
	"log"

	"blitzomni.com/m/controllers"
	"blitzomni.com/m/database"
	"blitzomni.com/m/routes"
	"blitzomni.com/m/utils"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var (
	server            *gin.Engine
	AuthController    controllers.AuthController
	AuthRoute         routes.AuthRoute
	ProductController controllers.ProductController
	ProductRoute      routes.ProductRoute
)

func init() {
	config, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatal("Could not load environment variables", err)
	}

	database.ConnectDB(&config)

	AuthController = controllers.NewAuthController(database.DB)
	AuthRoute = routes.NewAuthRoute(AuthController)
	ProductController = controllers.NewProductController(database.DB)
	ProductRoute = routes.NewProductRoute(ProductController)

	server = gin.Default()
}

func main() {
	config, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatal("? Could not load environment variables", err)
	}

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:8080", config.ClientOrigin}
	corsConfig.AllowCredentials = true

	server.Use(cors.New(corsConfig))

	router := server.Group("/api")

	AuthRoute.AuthRoute(router)
	ProductRoute.ProductRoute(router)
	log.Fatal(server.Run(":" + config.ServerPort))

}
