package common

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// Output holds together all of the data
type Output struct {
	AllArtists          []Artist  `json:"all_artists"`
	FloydWarshallMatrix [][]int   `json:"floyd_warshall_matrix"`
	ShortestPathsMatrix [][]int   `json:"shortest_paths_matrix"`
	IndexIDMap          []IndexID `json:"index_id_map"`
}

// IndexID holds the index-id mapping for the
// floyd warshall matrix
type IndexID struct {
	Index int    `json:"index"`
	ID    string `json:"id"`
}

// Artist is a struct representing what is stored in the database
type Artist struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Popularity  int      `json:"popularity"`
	Recommended []string `json:"recommended"`
	ImageData   Image    `json:"image_data"`
}

// Array is a type alias to allow for a receiver function
type Array []string

// Item is part of Spotify's JSON response and contains
// the Spotify ID, name of the artist, and the popularity
type Item struct {
	ID         string  `json:"id"`
	Name       string  `json:"name"`
	Popularity int     `json:"popularity"`
	Images     []Image `json:"images"`
}

// Image holds image data for the artist
type Image struct {
	Height int    `json:"height"`
	URL    string `json:"url"`
	Width  int    `json:"width"`
}

// RecommendedResponse is the response returned from the
// recommended endpoint
type RecommendedResponse struct {
	Artists []Item `json:"artists"`
}

// Artists is a part of Spotify's JSON response and contains
// the list of artists, the next page, current offset, and previous page
type Artists struct {
	Items    []Item `json:"items"`
	Next     string `json:"next"`
	Offset   int    `json:"offset"`
	Previous string `json:"previous"`
}

// Response is part of Spotify's JSON response and contains
// the data that we need
type Response struct {
	Artists Artists `json:"artists"`
}

// NotFoundError is the error returned by Spotify
type NotFoundError struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

// NotFound is the error returned by Spotify
type NotFound struct {
	Error NotFoundError `json:"error"`
}

// Value helps with insertion of json data
func (imageData Image) Value() (driver.Value, error) {
	return json.Marshal(imageData)
}

// Scan helps with the retrieval of json data
func (imageData *Image) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &imageData)
}

// Contains returns true if the value is in the array
func (arr Array) Contains(value string) bool {
	for i := 0; i < len(arr); i++ {
		if arr[i] == value {
			return true
		}
	}
	return false
}
