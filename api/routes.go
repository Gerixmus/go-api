package api

import (
	"github.com/gerixmus/go-api/database"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	router.GET("/albums", database.GetAlbums)
	router.GET("/albums/:id", database.GetAlbumByID)
	router.POST("/albums", database.PostAlbums)
}
