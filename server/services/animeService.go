package services

import (
	"context"
	"fmt"
	"io"
	"main/database"
	"main/models"
	"net/http"

	"golang.org/x/time/rate"
)

// regarde si l'anime est deja dans la liste d'utilisateur (pour le bouton ajouter a la liste)
func IsAnimeInUserList(userId int, animeID int) bool {
	var exists bool

	query := `
		SELECT EXISTS(
			SELECT 1
			FROM user_anime
			WHERE anime_id = $1
			  AND user_id = $2
		)
	`

	err := database.DB.QueryRow(query, animeID, userId).Scan(&exists)
	if err != nil {
		fmt.Println("SQL error:", err)
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
              AND user_id = $2 AND favorite = true
        )
    `
	err := database.DB.QueryRow(query, animeID, userID).Scan(&exists)
	if err != nil {
		fmt.Println("SQL error:", err)
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
		fmt.Println("SQL error:", err)
		return
	}
	if exists {
		return
	}

	// insérer l'anime
	query := `INSERT INTO anime (id, title, image, episodes) VALUES ($1, $2, $3, $4)`
	_, err = database.DB.Exec(query, animeID, title, image_url, episodes)
	if err != nil {
		fmt.Println("Error caching anime:", err)
	}
}

// ajoute l'anime dans la liste de l'utilisateur
func AddAnimeToUserList(userId int, animeId int) error {

	query := `
	INSERT INTO user_anime (user_id, anime_id, favorite)
	VALUES ($1, $2, false)
	ON CONFLICT (user_id, anime_id)
	DO NOTHING;
	`

	_, err := database.DB.Exec(query, userId, animeId)

	//fmt.Println(err)

	return err
}

// enleve l'anime de la liste de l'utilisateur
func RemoveAnimeFromUserList(userId int, animeId int) error {
	query := `
	DELETE from user_anime WHERE user_id = $1 and anime_id = $2
	`

	_, err := database.DB.Exec(query, userId, animeId)

	if err != nil {
		fmt.Println("Error removing anime from userlist:", err)
	}

	return err
}

// ajoute la review de l'utilisateur dans la base de données
func AddReview(userID int, body models.RequestAddReviewBody) error {

	query := `
    INSERT INTO reviews (user_id, anime_id, content, rating)
    VALUES ($1, $2, $3, $4)
    ON CONFLICT (user_id, anime_id)
    DO UPDATE SET content = $3, rating = $4;`

	_, err := database.DB.Exec(query, userID, body.AnimeID, body.Content, body.Rating)

	return err
}

func RemoveReview(reviewID int, userID int) error {

	query := `
	DELETE FROM reviews
	WHERE id = $1 AND user_id = $2
	`

	_, err := database.DB.Exec(query, reviewID, userID)

	if err != nil {
		fmt.Println("SQL error:", err)
		return err
	}

	return nil
}

func GetAnimeReviews(animeID int) ([]models.Review, error) {

	reviews := []models.Review{}

	query := `
	SELECT 
		id,
		user_id,
		anime_id,
		content,
		rating,
		created_at
	FROM reviews
	WHERE anime_id = $1
	ORDER BY created_at DESC
	`

	rows, err := database.DB.Query(query, animeID)

	if err != nil {
		return reviews, err
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
		)

		if err != nil {
			return reviews, err
		}

		reviews = append(reviews, review)
	}

	if err := rows.Err(); err != nil {
		return reviews, err
	}

	return reviews, nil
}

// verifie si l'utilisateur a deja ecrit une review
func AlreadyWroteReview(userID int, animeID int) bool {
	var exists bool
	query := "SELECT EXISTS (SELECT 1 from reviews where anime_id=$1 and user_id=$2)"
	err := database.DB.QueryRow(query, animeID, userID).Scan(&exists)

	if err != nil {
		fmt.Println("SQL error:", err)
		return false
	}

	return exists

}

func IsReviewFromUser(userID int, reviewID int) (bool, error) {

	var exists bool

	query := `
	SELECT EXISTS(
		SELECT 1
		FROM reviews
		WHERE id = $1 AND user_id = $2
	)
	`

	err := database.DB.QueryRow(
		query,
		reviewID,
		userID,
	).Scan(&exists)

	if err != nil {
		return false, err
	}

	return exists, nil

}

var jikan_url = "https://api.jikan.moe/v4"
var jikanLimiter = rate.NewLimiter(3, 1)

func GetAnimeFromJikan(animeId int) ([]byte, error) {

	err := jikanLimiter.Wait(context.Background())
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/anime/%d", jikan_url, animeId)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error request api Jikan: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("jikan error: status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading reponse Jikan: %w", err)
	}

	return body, nil
}
