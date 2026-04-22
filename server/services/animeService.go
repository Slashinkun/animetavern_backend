package services

import (
	"fmt"
	"main/database"
)

// regarde si l'anime est deja dans la liste d'utilisateur (pour le bouton ajouter a la liste)
func IsAnimeInUserList(userId int, animeID int) bool {
	var exists bool
	query := `
        SELECT EXISTS(
            SELECT 1
            FROM user_anime
            WHERE anime_id = $1
              AND user_id = (SELECT id FROM users WHERE username = $2)
        )
    `
	err := database.DB.QueryRow(query, animeID, userId).Scan(&exists)
	if err != nil {
		fmt.Println("Erreur SQL:", err)
		return false
	}
	return exists
}

// regarde si l'anime est deja dans les favoris de l'utilisateur (pour le bouton ajouter au favoris)
func IsAnimeUserFavorite(userID int, animeID int) bool {
	var exists bool
	query := `
        SELECT EXISTS(
            SELECT 1
            FROM user_anime
            WHERE anime_id = $1
              AND user_id = (SELECT id FROM users WHERE username = $2) AND favorite = true
        )
    `
	err := database.DB.QueryRow(query, animeID, userID).Scan(&exists)
	if err != nil {
		fmt.Println("Erreur SQL:", err)
		return false
	}
	return exists
}

// on met en cache les données de l'anime pour eviter de surcharger l'api de JIKAN (60 req par min)
func PutAnimeInCache(animeID int, title string, image_url string, episodes int) {
	var exists bool

	// vérifier si l'anime est déjà dans la DB
	queryVerif := `SELECT EXISTS(SELECT 1 FROM anime WHERE id = $1)`
	err := database.DB.QueryRow(queryVerif, animeID).Scan(&exists)
	if err != nil {
		fmt.Println("Erreur SQL:", err)
		return
	}
	if exists {
		return
	}

	// insérer l'anime
	query := `INSERT INTO anime (id, title, image, episodes) VALUES ($1, $2, $3, $4)`
	_, err = database.DB.Exec(query, animeID, title, image_url, episodes)
	if err != nil {
		fmt.Println("Erreur insertion anime:", err)
	}
}
