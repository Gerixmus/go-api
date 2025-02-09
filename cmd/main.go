package main

import (
	"database/sql"
	"log"
	"net/http"

	_ "go-api/docs"
	"go-api/internal/database"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func init() {
	// Load .env file for local development
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on environment variables")
	}
}

// album represents data about a record album.
// @Description Album represents data about a record album.
// @Description Note: This is a sample struct for demonstration purposes.
type album struct {
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

var db *sql.DB

// @title My Gin API
// @version 1.0
// @description This is a sample Gin API with Swagger documentation.
// @host localhost:8080
// @BasePath /
func main() {
	var err error

	db, err = database.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Initialize the Gin router
	router := gin.Default()
	router.GET("/albums", getAlbums)
	router.GET("/albums/:id", getAlbumByID)
	router.POST("/albums", postAlbums)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.Run("0.0.0.0:8080")
}

// getAlbums responds with the list of all albums as JSON.
// @Summary Albums
// @Description Get all the albums
// @Tags Albums
// @Accept json
// @Produce json
// @Success 200 {string} string "OK"
// @Router /albums [get]
func getAlbums(c *gin.Context) {
	var albums []album

	// Query the database
	rows, err := db.Query("SELECT title, artist, price FROM albums")
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "error querying the database"})
		return
	}
	defer rows.Close()

	// Iterate over the rows and populate the albums slice
	for rows.Next() {
		var a album
		if err := rows.Scan(&a.Title, &a.Artist, &a.Price); err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "error scanning rows"})
			return
		}
		albums = append(albums, a)
	}

	c.IndentedJSON(http.StatusOK, albums)
}

// postAlbums adds an album from JSON received in the request body.
// @Summary Add a new album
// @Description Add a new album to the list
// @Tags Albums
// @Accept json
// @Produce json
// @Param album body album true "Album to add"
// @Success 201 {object} album "Album added"
// @Failure 400 {string} string "Bad Request"
// @Router /albums [post]
func postAlbums(c *gin.Context) {
	var newAlbum album

	// Bind the received JSON to newAlbum
	if err := c.BindJSON(&newAlbum); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid request payload"})
		return
	}

	// Insert the new album into the database
	result, err := db.Exec("INSERT INTO albums (title, artist, price) VALUES (?, ?, ?)", newAlbum.Title, newAlbum.Artist, newAlbum.Price)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "error inserting into database"})
		return
	}

	// Get the ID of the inserted album
	id, err := result.LastInsertId()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "error getting last insert ID"})
		return
	}

	log.Println(id)
	c.IndentedJSON(http.StatusCreated, newAlbum)
}

// getAlbumByID locates and returns an album as a response.
// @Summary Get an album by ID
// @Description Get an album by its ID
// @Tags Albums
// @Accept json
// @Produce json
// @Param id path string true "Album ID"
// @Success 200 {object} album "OK"
// @Failure 404 {string} string "Not Found"
// @Router /albums/{id} [get]
func getAlbumByID(c *gin.Context) {
	id := c.Param("id")

	var a album

	// Query the database for the album with the specified ID
	row := db.QueryRow("SELECT title, artist, price FROM albums WHERE id = ?", id)
	if err := row.Scan(&a.Title, &a.Artist, &a.Price); err != nil {
		if err == sql.ErrNoRows {
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
			return
		}
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "error querying the database"})
		return
	}

	c.IndentedJSON(http.StatusOK, a)
}
