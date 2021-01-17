package common

// Item is part of Spotify's JSON response and contains
// the Spotify ID, name of the artist, and the popularity
type Item struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Popularity int    `json:"popularity"`
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
