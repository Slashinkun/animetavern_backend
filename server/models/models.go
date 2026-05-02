package models

type User struct {
	ID       int
	Username string
	Password string
}

// structure pour le cache
type Anime struct {
	ID       int    `json:"mal_id"`
	Title    string `json:"title"`
	Image    string `json:"image_url"`
	Episodes int    `json:"episodes"`
}

// structure pour le profil de l'utilisateur
type UserAnime struct {
	ID             int    `json:"mal_id"`
	Title          string `json:"title"`
	Image          string `json:"image_url"`
	Episodes       int    `json:"episodes"`
	ViewedEpisodes int    `json:"viewed_episodes"`
	Status         string `json:"status"`
}
