package services

import (
	"fmt"
	"main/database"
	"main/models"
)

// cherche si l'utilisateur existe
func GetUserByID(userID int) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM users where id = $1)`
	err := database.DB.QueryRow(query, userID).Scan(&exists)

	if err != nil {
		fmt.Println("Erreur SQL:", err)
		return false, err
	}

	return exists, nil
}

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

//A FAIRE ???
// func GetAnimeReviews(animeID int) (models.AnimeReviews,error){

// }

func GetUserReviews(userID int) (models.UserReviews, error) {
	var data models.UserReviews
	data.Reviews = []models.AnimeUserReview{}

	query := `
		SELECT 
			r.id, r.user_id, r.anime_id,
    		r.content, r.rating, r.created_at,
    		a.title, a.image
		FROM reviews r
		JOIN anime a ON r.anime_id = a.id
		WHERE r.user_id = $1;`

	rows, err := database.DB.Query(query, userID)

	if err != nil {
		return data, err
	}
	defer rows.Close()

	for rows.Next() {
		var review models.AnimeUserReview

		err := rows.Scan(
			&review.ID,
			&review.UserID,
			&review.AnimeID,
			&review.Content,
			&review.Rating,
			&review.CreatedAt,
			&review.AnimeTitle,
			&review.AnimeImage,
		)
		if err != nil {
			return data, err
		}

		data.Reviews = append(data.Reviews, review)
	}

	if err := rows.Err(); err != nil {
		return data, err
	}

	return data, nil
}

func GetUsername(userId int) (string, error) {
	var username string

	query := `SELECT username FROM users WHERE id = $1`
	err := database.DB.QueryRow(query, userId).Scan(&username)
	if err != nil {
		// retourne l'erreur pour que le handler puisse gérer "user not found"
		return "", fmt.Errorf("user %s not found: %w", username, err)
	}
	return username, nil
}
