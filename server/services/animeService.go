package services

import (
	"fmt"
	"main/database"
	"main/models"
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
		fmt.Println("Erreur SQL:", err)
		return false
	}

	//fmt.Printf("inList : %t\n", exists)

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
		fmt.Println("Erreur suppression anime:", err)
	}

	return err
}

// ajoute la review de l'utilisateur dans la base de données
func AddReview(userID int, body models.ReviewBody) error {

	query := `
    INSERT INTO reviews (user_id, anime_id, content, rating)
    VALUES ($1, $2, $3, $4)
    ON CONFLICT (user_id, anime_id)
    DO UPDATE SET content = $3, rating = $4;`

	_, err := database.DB.Exec(query, userID, body.AnimeID, body.Content, body.Rating)

	//fmt.Println(err)

	return err
}

func RemoveReview(reviewID int, userID int) error {

	query := `
	DELETE FROM reviews
	WHERE id = $1 AND user_id = $2
	`

	_, err := database.DB.Exec(query, reviewID, userID)

	if err != nil {
		fmt.Println("Erreur SQL:", err)
		return err
	}

	return nil
}

// verifie si l'utilisateur a deja ecrit une review
func AlreadyWroteReview(userID int, animeID int) bool {
	var exists bool
	query := "SELECT EXISTS (SELECT 1 from reviews where anime_id=$1 and user_id=$2)"
	err := database.DB.QueryRow(query, animeID, userID).Scan(&exists)

	if err != nil {
		fmt.Println("Erreur SQL:", err)
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
