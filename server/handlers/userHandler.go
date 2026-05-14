package handlers

import (
	"database/sql"
	"encoding/json"
	"main/contextkeys"
	"main/models"
	"main/services"
	"main/utils"
	"net/http"
)

// page utilisateur
func UserPageHandler(w http.ResponseWriter, r *http.Request) {

	//on recuperer l'id dans l'url
	userId, err := utils.GetIntParam(r, "id")
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

	//fmt.Printf("page id : %d  currentid : %d \n", userId, currentUserID)

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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

}

// page des reviews ecrit par l'utilisateur
func UserReviewsHandler(w http.ResponseWriter, r *http.Request) {

	//on recuperer l'id dans l'url

	userId, err := utils.GetIntParam(r, "id")
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
	var body models.RequestAddAnimeBody

	//on recupere l'id depuis les cookies
	userId, ok := utils.GetUserID(r)

	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	err := utils.DecodeJSON(r, &body)
	if err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}

	err = utils.ValidateAddAnimeBody(body)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//verifie si l'utilisateur a deja l'anime dans sa liste
	exists := services.IsAnimeInUserList(userId, body.AnimeID)
	if exists {
		http.Error(w, "Anime already in the list", http.StatusConflict)
		return
	}

	err = services.AddAnimeToUserList(userId, body.AnimeID)
	if err != nil {
		http.Error(w, "Failed to add anime", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Anime added",
	})

}

func RemoveAnimeFromList(w http.ResponseWriter, r *http.Request) {

	animeId, err := utils.GetIntParam(r, "id")
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	userId, ok := utils.GetUserID(r)

	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	exists := services.IsAnimeInUserList(userId, animeId)
	if !exists {
		http.Error(w, "Anime not in list", http.StatusNotFound)
		return
	}

	err = services.RemoveAnimeFromUserList(userId, animeId)
	if err != nil {
		http.Error(w, "Failed to remove anime from list", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func UpdateEpisodes(w http.ResponseWriter, r *http.Request) {
	var body models.RequestUpdateEpisode

	animeId, err := utils.GetIntParam(r, "id")

	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	userId, ok := utils.GetUserID(r)

	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	exists := services.IsAnimeInUserList(userId, animeId)
	if !exists {
		http.Error(w, "Anime not in list", http.StatusNotFound)
		return
	}

	err = utils.DecodeJSON(r, &body)
	if err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	err = utils.ValidateUpdateEpisodeBody(body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	updated, err := services.UpdateEpisodes(userId, animeId, body.Delta)

	if err != nil {
		http.Error(w, "Cannot update episode counter", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{
		"viewed_episodes": updated,
	})
}

func UpdateAnimeStatus(w http.ResponseWriter, r *http.Request) {

	var body models.RequestUpdateAnimeStatus

	animeId, err := utils.GetIntParam(r, "id")

	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	userId, ok := utils.GetUserID(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	exists := services.IsAnimeInUserList(userId, animeId)
	if !exists {
		http.Error(w, "Anime not in list", http.StatusNotFound)
		return
	}

	err = utils.DecodeJSON(r, &body)
	if err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	err = utils.ValidateAnimeStatusBody(body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newStatus, err := services.UpdateAnimeStatus(userId, animeId, body.Status)
	if err != nil {
		http.Error(w, "Cannot update status", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status": newStatus,
	})
}
