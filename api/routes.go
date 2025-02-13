package api

import (
	"github.com/gerixmus/go-api/database"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRoutes(router *gin.Engine) {
	router.GET("/albums", database.GetAlbums)
	router.GET("/albums/:id", database.GetAlbumByID)
	router.POST("/albums", database.PostAlbums)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
