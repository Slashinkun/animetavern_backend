package main

import (
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

var jikan_url = "https://api.jikan.moe/v4"
var secret = []byte("secret_key")
var connStr = "user=myuser password=mypassword dbname=animetavern_db sslmode=disable"

func getAnimePage(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 3 {
		http.Error(w, "ID manquant", http.StatusBadRequest)
		return
	}
	id := pathParts[2]
	url := fmt.Sprintf("%s/anime/%s", jikan_url, id)
	fmt.Println(url)
	resp, err := http.Get(url)
	if err != nil {
		http.Error(w, "Erreur requete api Jikan", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		http.Error(w, "Erreur lecture reponse jikan", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
	w.Write(body)

}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func checkPassword(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func createToken(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
	})

	return token.SignedString(secret)
}

// func testConnectionDB() {
// 	db, err := sql.Open("postgres", connStr)
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer db.Close()
// 	err = db.Ping()
// 	if err != nil {
// 		panic(err)
// 	}

// 	fmt.Println("Successfully connected!")
// }

func login(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	var hashed string
	db, err := sql.Open("postgres", connStr)
	user := db.QueryRow("SELECT password FROM users WHERE username=$1", username).Scan(&hashed)
	defer db.Close()
	if user != nil {
		http.Error(w, "Utilisateur introuvable", 401)
		return
	}

	if !checkPassword(hashed, password) {
		http.Error(w, "Mot de passe incorrect", 401)
		return
	}

	//creation du token d'authentification
	token, err := createToken(username)
	if err != nil {
		http.Error(w, "Erreur token", 500)
		return
	}
	expirationTime := time.Now().Add(1 * time.Hour)

	//ecriture du cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		Expires:  expirationTime,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})

}

func userPage(w http.ResponseWriter, r *http.Request) {

}

func register(w http.ResponseWriter, r *http.Request) {

	//a faire passer au proxy react ou attendre deploiment render
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Erreur formulaire", 400)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	if username == "" || password == "" {
		http.Error(w, "Username ou password vide", 400)
		return
	}

	//on hash le mot de passe
	hashed, err := hashPassword(password)
	if err != nil {
		http.Error(w, "Erreur hash", http.StatusInternalServerError)
		fmt.Println("Hash error:", err)
		return
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		http.Error(w, "Erreur connexion DB", http.StatusInternalServerError)
		fmt.Println("DB connection error:", err)
		return
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		http.Error(w, "Impossible de joindre la DB", http.StatusInternalServerError)
		fmt.Println("DB ping error:", err)
		return
	}

	//on met les infos de l'utilisateur dans la BDD
	_, err = db.Exec("INSERT INTO users(username, password) VALUES($1, $2)", username, hashed)
	if err != nil {
		http.Error(w, "Erreur insertion DB", http.StatusInternalServerError)
		fmt.Println("DB insert error:", err)
		return
	}

	//a faire rediriger l'utilisateur a la page de login
	w.Write([]byte("Utilisateur créé"))
}

func main() {

	http.HandleFunc("/anime/", getAnimePage)

	http.HandleFunc("/register", register)
	http.HandleFunc("/login", login)
	//testConnectionDB()

	fmt.Println("Serveur démarré sur le port 8080...")
	// Démarre l'écoute sur le port 8080
	http.ListenAndServe(":8080", nil)
}
