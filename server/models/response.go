package models

import "encoding/json"

//Structures pour la réponse du serveur (et la réponse de l'API externe)

type UserData struct {
	Username  string      `json:"username"`
	Animes    []UserAnime `json:"animes"`
	Favorites []Anime     `json:"favoriteAnimes"`
}

type UserReviews struct {
	Reviews []AnimeUserReview `json:"reviews"`
}

type AnimeResponse struct {
	Anime      json.RawMessage `json:"anime"`
	IsInList   bool            `json:"isInList"`
	IsFavorite bool            `json:"isFavorite"`
	Reviews    []Reviews       `json:"reviews"`
}

type AnimeUserReview struct {
	ID         int    `json:"id"`
	UserID     int    `json:"user_id"`
	AnimeID    int    `json:"anime_id"`
	AnimeTitle string `json:"anime_title"`
	AnimeImage string `json:"anime_image"`
	Content    string `json:"content"`
	Rating     int    `json:"rating"`
	CreatedAt  string `json:"created_at"`
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
