package models

import "encoding/json"

//Structures pour la réponse du serveur (et la réponse de l'API externe)

type UserData struct {
	Username  string      `json:"username"`
	Animes    []UserAnime `json:"animes"`
	Favorites []Anime     `json:"favoriteAnimes"`
}

type UserReviews struct {
	Reviews []Review `json:"reviews"`
}

type AnimeResponse struct {
	Anime      json.RawMessage `json:"anime"`
	IsInList   bool            `json:"isInList"`
	IsFavorite bool            `json:"isFavorite"`
	Reviews    []Review        `json:"reviews"`
}

type RequestAddAnimeBody struct {
	AnimeID int `json:"anime_id"`
}

type RequestUpdateEpisode struct {
	Delta int `json:"delta"`
}

type RequestUpdateAnimeStatus struct {
	Status string `json:"status"`
}

type AnimeJikanResponse struct {
	Data AnimeData `json:"data"`
}

type MeResponse struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
}

type RequestAddReviewBody struct {
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
