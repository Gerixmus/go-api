package main

import (
	"log"

	"github.com/gerixmus/go-api/api"
	_ "github.com/gerixmus/go-api/docs"

	"github.com/gerixmus/go-api/database"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	// Load .env file for local development
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on environment variables")
	}
}

// @title My Gin API
// @version 1.0
// @description This is a sample Gin API with Swagger documentation.
// @host localhost:8080
// @BasePath /
func main() {
	database.Connect()
	// Initialize the Gin router
	router := gin.Default()
	api.SetupRoutes(router)
	router.Run("0.0.0.0:8080")
}
