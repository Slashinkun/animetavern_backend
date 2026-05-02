package models

import "encoding/json"

type UserData struct {
	Username  string      `json:"username"`
	Animes    []UserAnime `json:"animes"`
	Favorites []Anime     `json:"favoriteAnimes"`
}

type AnimeResponse struct {
	Anime      json.RawMessage `json:"anime"`
	IsInList   bool            `json:"isInList"`
	IsFavorite bool            `json:"isFavorite"`
}

type RequestBody struct {
	AnimeID int `json:"anime_id"`
}

type AnimeJikanResponse struct {
	Data AnimeData `json:"data"`
}

type MeResponse struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
}

type ReviewBody struct {
	AnimeID int    `json:"anime_id"`
	Content string `json:"content"`
	Rating  int    `json:"rating"`
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
