package services

import (
	"fmt"
	"main/database"
	"main/models"
)

// recuperer les infos de l'utlisateur pour sa page perso
func GetUserData(userID int) (models.UserData, error) {
	var data models.UserData
	data.Animes = []models.UserAnime{}
	data.Favorites = []models.Anime{}

	username, err := GetUsername(userID)

	if err != nil {
		fmt.Println(err)
	}
	data.Username = username

	queryList := `SELECT a.id, a.title, a.episodes, a.image, ua.viewed_episodes,ua.status
				FROM anime a
				JOIN user_anime ua ON a.id = ua.anime_id
				WHERE ua.user_id = $1;`

	rowsList, err := database.DB.Query(queryList, userID)
	if err != nil {
		fmt.Println(err)
		return data, err

	}
	defer rowsList.Close()

	for rowsList.Next() {
		var a models.UserAnime
		if err := rowsList.Scan(&a.ID, &a.Title, &a.Episodes, &a.Image, &a.ViewedEpisodes, &a.Status); err != nil {
			return data, err
		}
		data.Animes = append(data.Animes, a)
	}

	queryFavorite := `
        SELECT a.id, a.title, a.episodes, a.image
        FROM anime a
        JOIN user_anime ua ON a.id = ua.anime_id
        WHERE ua.user_id = $1 AND ua.favorite = true;`

	rowsFavorites, err := database.DB.Query(queryFavorite, userID)
	if err != nil {
		return data, err
	}
	defer rowsFavorites.Close()

	for rowsFavorites.Next() {
		var a models.Anime
		if err := rowsFavorites.Scan(&a.ID, &a.Title, &a.Episodes, &a.Image); err != nil {
			return data, err
		}
		data.Favorites = append(data.Favorites, a)
	}

	return data, nil

}

func GetUsername(userId int) (string, error) {
	var username string
	//fmt.Println(username)
	query := `SELECT username FROM users WHERE id = $1`
	err := database.DB.QueryRow(query, userId).Scan(&username)
	if err != nil {
		// retourne l'erreur pour que le handler puisse gérer "user not found"
		return "", fmt.Errorf("user %s not found: %w", username, err)
	}
	return username, nil
}
