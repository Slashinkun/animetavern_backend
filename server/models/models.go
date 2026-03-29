package models

import (
	"encoding/json"
)

type User struct {
	id       int
	username string
	password string
}

type AnimeResponse struct {
	Anime      json.RawMessage `json:"anime"`
	IsInList   bool            `json:"isInList"`
	IsFavorite bool            `json:"isFavorite"`
}

type Anime struct {
	ID       int    `json:"mal_id"`
	Title    string `json:"title"`
	Image    string `json:"image_url"`
	Episodes int    `json:"episodes"`
}

type AnimeData struct {
	MalID  int    `json:"mal_id"`
	Title  string `json:"title"`
	Images struct {
		JPG struct {
			ImageURL string `json:"image_url"`
		} `json:"jpg"`
	} `json:"images"`
	Episodes int `json:"episodes"`
}

type AnimeJikanResponse struct {
	Data AnimeData `json:"data"`
}
