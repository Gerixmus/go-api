package database

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gerixmus/go-api/models"
	"github.com/gin-gonic/gin"
)

// getAlbums responds with the list of all albums as JSON.
// @Summary Albums
// @Description Get all the albums
// @Tags Albums
// @Accept json
// @Produce json
// @Success 200 {string} string "OK"
// @Router /albums [get]
func GetAlbums(c *gin.Context) {
	var albums []models.Album

	// Query the database
	rows, err := DB.Query("SELECT title, artist, price FROM albums")
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "error querying the database"})
		return
	}
	defer rows.Close()

	// Iterate over the rows and populate the albums slice
	for rows.Next() {
		var a models.Album
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
func PostAlbums(c *gin.Context) {
	var newAlbum models.Album

	// Bind the received JSON to newAlbum
	if err := c.BindJSON(&newAlbum); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid request payload"})
		return
	}

	// Insert the new album into the database
	result, err := DB.Exec("INSERT INTO albums (title, artist, price) VALUES (?, ?, ?)", newAlbum.Title, newAlbum.Artist, newAlbum.Price)
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
func GetAlbumByID(c *gin.Context) {
	id := c.Param("id")

	var a models.Album

	// Query the database for the album with the specified ID
	row := DB.QueryRow("SELECT title, artist, price FROM albums WHERE id = ?", id)
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
