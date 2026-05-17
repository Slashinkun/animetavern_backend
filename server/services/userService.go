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
		fmt.Println("SQL error:", err)
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
		//fmt.Println(err)
		return data, err
	}
	data.Username = username

	queryList := `SELECT a.id, a.title, a.episodes, a.image, ua.viewed_episodes,ua.status
				FROM anime a
				JOIN user_anime ua ON a.id = ua.anime_id
				WHERE ua.user_id = $1;`

	rowsList, err := database.DB.Query(queryList, userID)
	if err != nil {
		//fmt.Println(err)
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

func GetUserReviews(userID int) (models.UserReviews, error) {

	var data models.UserReviews

	data.Reviews = []models.Review{}

	query := `
	SELECT 
		r.id,
		r.user_id,
		r.anime_id,
		r.content,
		r.rating,
		r.created_at,
		a.title,
		a.image
	FROM reviews r
	JOIN anime a ON r.anime_id = a.id
	WHERE r.user_id = $1
	ORDER BY r.created_at DESC
	`

	rows, err := database.DB.Query(query, userID)

	if err != nil {
		return data, err
	}

	defer rows.Close()

	for rows.Next() {

		var review models.Review

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

func UpdateEpisodes(userId int, animeId int, delta int) (models.UserAnime, error) {

	query := `
	UPDATE user_anime
	SET viewed_episodes = GREATEST(viewed_episodes + $1, 0)
	WHERE user_id = $2 AND anime_id = $3
	RETURNING anime_id, viewed_episodes, status
	`

	var result models.UserAnime

	err := database.DB.QueryRow(query, delta, userId, animeId).Scan(
		&result.ID,
		&result.ViewedEpisodes,
		&result.Status,
	)

	if err != nil {
		return result, err
	}

	return result, nil
}

func UpdateAnimeStatus(userId int, animeId int, status string) (string, error) {

	query := `
	UPDATE user_anime
	SET status = $1
	WHERE user_id = $2 AND anime_id = $3
	RETURNING status
	`

	var newStatus string

	err := database.DB.QueryRow(query, status, userId, animeId).Scan(&newStatus)

	return newStatus, err
}

func UpdateAnimeFavorite(userId int, animeId int, favorite bool) error {
	query := `
		UPDATE user_anime
		SET favorite = $1
		WHERE user_id = $2 AND anime_id = $3
	`

	_, err := database.DB.Exec(query, favorite, userId, animeId)

	return err
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
