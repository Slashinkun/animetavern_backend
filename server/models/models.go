package models

import (
	"encoding/json"
)

type User struct {
	ID       int
	Username string
	Password string
}

type UserData struct {
	Username  string  `json:"username"`
	Animes    []Anime `json:"animes"`
	Favorites []Anime `json:"favoriteAnimes"`
}

type AnimeResponse struct {
	Anime      json.RawMessage `json:"anime"`
	IsInList   bool            `json:"isInList"`
	IsFavorite bool            `json:"isFavorite"`
}

// structure pour le cache
type Anime struct {
	ID       int    `json:"mal_id"`
	Title    string `json:"title"`
	Image    string `json:"image_url"`
	Episodes int    `json:"episodes"`
}

type AnimeData struct {
	MalID  int    `json:"mal_id"`
	Title  string `json:"title_english"`
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

type MeResponse struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
}
