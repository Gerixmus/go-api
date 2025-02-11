package models

// album represents data about a record album.
// @Description Album represents data about a record album.
// @Description Note: This is a sample struct for demonstration purposes.
type Album struct {
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}
