package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "go-api/docs"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
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
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPassword, dbHost, dbPort, dbName)
	// Connect to the MySQL database
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	// Ensure the connection is available
	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}

	// Initialize the Gin router
	router := gin.Default()
	router.GET("/albums", getAlbums)
	router.GET("/albums/:id", getAlbumByID)
	router.POST("/albums", postAlbums)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.Run("localhost:8080")
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

// getAlbumByID locates the album whose ID value matches the id
// parameter sent by the client, then returns that album as a response.
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
