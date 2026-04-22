package handlers

import (
	"encoding/json"
	"fmt"
	"main/models"
	"main/services"
	"net/http"
	"time"
)

// POST /register
func RegisterHandler(w http.ResponseWriter, r *http.Request) {

	//preflight
	if r.Method == http.MethodOptions {
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

	email := r.FormValue("email")
	username := r.FormValue("username")
	password := r.FormValue("password")

	if username == "" || password == "" || email == "" {
		http.Error(w, "Username, password ou email vide", 400)
		return
	}

	err := services.Register(email, username, password)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), 500)
		return
	}

	w.WriteHeader(http.StatusCreated)

}

// POST /login
func LoginHandler(w http.ResponseWriter, r *http.Request) {

	//recuperer la valeur des champs du formulaire
	email := r.FormValue("email")
	password := r.FormValue("password")

	//connexion
	token, err := services.Login(email, password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	expirationTime := time.Now().Add(24 * time.Hour)

	//ecriture du cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    token,
		Expires:  expirationTime,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})

}

// deconnection
func LogoutHandler(w http.ResponseWriter, r *http.Request) {

	//on change le temps d'expiration du cookie pour qu'il disparaisse immediatement
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
	})
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Déconnecté"))
}

// pour l'authentification
func MeHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "method not allowed",
		})
		return
	}

	cookie, err := r.Cookie("session_token")

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "not logged in",
		})
		return
	}

	//fmt.Println("Cookies:", r.Cookies())

	userId, err := services.ValidateToken(cookie.Value)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "invalid token",
		})
		return
	}

	username, err := services.GetUsername(userId)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "user not found",
		})
		return
	}

	json.NewEncoder(w).Encode(models.MeResponse{
		ID:       userId,
		Username: username,
	})
}
