package services

import (
	"database/sql"
	"errors"
	"fmt"
	"main/database"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var secret = []byte("secret_key")

// connection
func Login(email, password string) (string, error) {
	var hashedPassword string
	var userId int

	err := database.DB.QueryRow(
		"SELECT id,password FROM users WHERE email=$1",
		email,
	).Scan(&userId, &hashedPassword)

	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("invalid email")
		}
		return "", err
	}

	if err := checkPassword(hashedPassword, password); err != nil {
		return "", fmt.Errorf("incorrect password")
	}

	//err = database.DB.QueryRow("SELECT username FROM users WHERE email=$1").Scan(&username)

	token, err := createToken(userId)
	if err != nil {
		return "", err
	}

	//mt.Println(token)

	return token, err
}

// inscription
func Register(email, username, password string) error {
	hashedPassword, err := hashPassword(password)
	if err != nil {
		return fmt.Errorf("error hash password: %w", err)
	}

	//creation de l'utilisateur dans la DB
	_, err = database.DB.Exec(
		"INSERT INTO users (email, username, password) VALUES ($1,$2,$3)",
		email,
		username,
		string(hashedPassword),
	)

	//l'utilisateur existe deja ou l'email est deja utilisé
	if err != nil {
		if strings.Contains(err.Error(), "unique constraint") {
			return fmt.Errorf("email or username already used")
		}
		return err
	}

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
func createToken(userId int) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userId,
		"exp":     jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // token valide 24h
	})

	return token.SignedString(secret)
}

// regarde si le token de session est valide
func ValidateToken(tokenString string) (int, error) {

	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})

	// fmt.Println("CLAIMS:", claims)
	// fmt.Println("USER_ID RAW:", claims["user_id"])

	if err != nil {
		return 0, err
	}
	userIDFloat, ok := claims["user_id"].(float64)

	//fmt.Print(userIDFloat)

	if !ok {
		return 0, errors.New("user_id not found")
	}

	return int(userIDFloat), nil
}
