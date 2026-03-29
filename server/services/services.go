package services

import (
	"errors"
	"fmt"
	"main/database"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var connStr = "user=myuser password=mypassword dbname=animetavern_db sslmode=disable"

var secret = []byte("secret_key")

// connection
func Login(username, password string) (string, error) {
	var hashedPassword string
	database.DB.QueryRow(
		"SELECT password FROM users WHERE username=$1",
		username,
	).Scan(&hashedPassword)

	err := checkPassword(hashedPassword, password)
	token, err := createToken(username)

	return token, err
}

// inscription
func Register(username, password string) error {
	hashedPassword, _ := hashPassword(password)
	_, err := database.DB.Exec(
		"INSERT INTO users (username, password) VALUES ($1, $2)",
		username,
		string(hashedPassword),
	)

	return err
}

// on hash le password pour eviter de stocker en dur le mdp
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func checkPassword(hash, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err
}

// cree le token de session
func createToken(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
	})

	return token.SignedString(secret)
}

// regarde si le token de session est valide
func ValidateToken(tokenString string) (string, error) {

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})

	if err != nil || !token.Valid {
		return "", errors.New("token invalide")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("claims invalides")
	}

	username, ok := claims["username"].(string)
	if !ok {
		return "", errors.New("username introuvable")
	}

	return username, nil
}

func IsAnimeInUserList(username string, animeID int) bool {
	var exists bool
	query := `
        SELECT EXISTS(
            SELECT 1
            FROM user_anime
            WHERE anime_id = $1
              AND user_id = (SELECT id FROM users WHERE username = $2)
        )
    `
	err := database.DB.QueryRow(query, animeID, username).Scan(&exists)
	if err != nil {
		fmt.Println("Erreur SQL:", err)
		return false
	}
	return exists
}

func IsAnimeUserFavorite(username string, animeID int) bool {
	var exists bool
	query := `
        SELECT EXISTS(
            SELECT 1
            FROM user_anime
            WHERE anime_id = $1
              AND user_id = (SELECT id FROM users WHERE username = $2) AND favorite = true
        )
    `
	err := database.DB.QueryRow(query, animeID, username).Scan(&exists)
	if err != nil {
		fmt.Println("Erreur SQL:", err)
		return false
	}
	return exists
}

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
