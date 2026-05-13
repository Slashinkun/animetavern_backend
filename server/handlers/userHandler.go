package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"main/contextkeys"
	"main/models"
	"main/services"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// page utilisateur
func UserPageHandler(w http.ResponseWriter, r *http.Request) {

	//preflight
	if r.Method == http.MethodOptions {
		//w.WriteHeader(http.StatusOK)
		return
	}

	//on recuperer l'id dans l'url
	vars := mux.Vars(r)

	// fmt.Println("vars:", vars)
	// fmt.Println("id string:", vars["id"])

	idStr, ok := vars["id"]
	if !ok || idStr == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}

	userId, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	//gestion cas utilisateur n'existe pas
	userExists, err := services.GetUserByID(userId)

	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	if !userExists {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	//on regarde si l'utilisateur qui demande la page est le proprietaire ou pas
	var currentUserID int
	currentUserID = -1
	if idCtx := r.Context().Value(contextkeys.UserID); idCtx != nil {
		currentUserID = idCtx.(int)
	}
	isUser := (currentUserID == userId)

	fmt.Printf("page id : %d  currentid : %d \n", userId, currentUserID)

	//récupération des données et construction de la reponse
	userData, err := services.GetUserData(userId)
	if err != nil {
		http.Error(w, "error fetching data", http.StatusInternalServerError)
		return
	}

	response := struct {
		models.UserData
		IsUser bool `json:"is_user"`
	}{
		UserData: userData,
		IsUser:   isUser,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

}

// page des reviews ecrit par l'utilisateur
func UserReviewsHandler(w http.ResponseWriter, r *http.Request) {

	//preflight
	if r.Method == http.MethodOptions {
		//w.WriteHeader(http.StatusOK)
		return
	}

	//on recuperer l'id dans l'url
	vars := mux.Vars(r)

	idStr, ok := vars["id"]
	if !ok || idStr == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}

	userId, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	userReviews, err := services.GetUserReviews(userId)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		http.Error(w, "error fetching data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(userReviews)

}

func AddAnimeToUserList(w http.ResponseWriter, r *http.Request) {
	var body models.RequestBody

	//fmt.Println("REQUEST COOKIES:", r.Cookies())

	//on recupere l'id depuis les cookies
	userId, ok := r.Context().Value(contextkeys.UserID).(int)

	//fmt.Println(userId)

	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}

	//verifie si l'utilisateur a deja l'anime dans sa liste
	exists := services.IsAnimeInUserList(userId, body.AnimeID)
	if exists {
		http.Error(w, "Anime already in the list", http.StatusConflict)
		return
	}

	fmt.Println("[INFO] : Trying adding anime with id " + strconv.Itoa(body.AnimeID) + " for user " + strconv.Itoa(userId))

	err = services.AddAnimeToUserList(userId, body.AnimeID)
	if err != nil {
		http.Error(w, "Failed to add anime", http.StatusInternalServerError)
		return
	}

	fmt.Println("[INFO] : Anime " + strconv.Itoa(body.AnimeID) + " successfuly added for user " + strconv.Itoa(userId))

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Anime added",
	})

}

func RemoveAnimeFromList(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	// fmt.Println("vars:", vars)
	// fmt.Println("id string:", vars["id"])

	idStr, ok := vars["id"]
	if !ok || idStr == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}

	animeId, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	userId, ok := r.Context().Value(contextkeys.UserID).(int)

	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	exists := services.IsAnimeInUserList(userId, animeId)
	if !exists {
		http.Error(w, "Anime not in list", http.StatusConflict)
		return
	}

	err = services.RemoveAnimeFromUserList(userId, animeId)
	if err != nil {
		http.Error(w, "Failed to remove anime from list", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
