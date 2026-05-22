package handlers

import (
	"encoding/json"
	"main/models"
	"main/services"
	"main/utils"
	"net/http"
)

// page anime
func GetAnimePage(w http.ResponseWriter, r *http.Request) {

	//variable pour l'affichage des boutons
	var inList bool
	var inFavorite bool

	if r.Method != http.MethodGet {
		http.Error(w, "Unauthorized", http.StatusMethodNotAllowed)
		return
	}

	//récuperation de l'id dans l'url
	animeId, err := utils.GetIntParam(r, "anime_id")

	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	animedata, err := services.GetAnimeFromJikan(animeId)

	if err != nil {
		http.Error(w, "Jikan API error", http.StatusInternalServerError)
		return
	}

	//on transforme la reponse de l'API faciliter la lecture
	var resp models.AnimeJikanResponse
	err = json.Unmarshal(animedata, &resp)
	if err != nil {
		http.Error(w, "Jikan parsing error", http.StatusInternalServerError)
		return
	}

	anime := resp.Data

	//on met l'anime dans le cache si il est pas deja dans le cache
	services.PutAnimeInCache(animeId, anime.Title, anime.Images.JPG.ImageURL, anime.Episodes)

	userId, ok := utils.GetUserID(r)

	if !ok {
		inList = false
		inFavorite = false
	} else {
		inList = services.IsAnimeInUserList(userId, animeId)
		inFavorite = services.IsAnimeUserFavorite(userId, animeId)
	}

	reviews, err := services.GetAnimeReviews(animeId)

	if err != nil {
		http.Error(w, "Unable to get anime's reviews", http.StatusInternalServerError)
		return
	}

	//on construit la réponse
	data := models.AnimeResponse{
		Anime:      animedata,
		IsInList:   inList,
		IsFavorite: inFavorite,
		Reviews:    reviews,
	}

	//on la transforme en JSON
	response, err := json.Marshal(data)
	if err != nil {
		http.Error(w, "JSON error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)

}

func AddReview(w http.ResponseWriter, r *http.Request) {

	var body models.RequestAddReviewBody

	animeID, err := utils.GetIntParam(r, "id")
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	//vérification utilisateur connecté
	userId, ok := utils.GetUserID(r)

	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	err = utils.DecodeJSON(r, &body)

	if err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}

	err = utils.ValidateReviewBody(body)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	exists := services.AlreadyWroteReview(userId, animeID)
	if exists {
		http.Error(w, "Already wrote a review", http.StatusConflict)
		return
	}

	err = services.AddReview(userId, animeID, body)

	if err != nil {
		http.Error(w, "Error while adding review ", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Review added",
	})

}

func RemoveReview(w http.ResponseWriter, r *http.Request) {

	//récuperation de l'id dans l'url
	reviewID, err := utils.GetIntParam(r, "id")

	if err != nil {
		http.Error(w, "invalid review id", http.StatusBadRequest)
		return
	}

	//vérification utilisateur connecté
	userId, ok := utils.GetUserID(r)

	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var isUserReview bool
	isUserReview, err = services.IsReviewFromUser(userId, reviewID)

	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	if !isUserReview {
		http.Error(w, "Cannot delete others reviews", http.StatusForbidden)
		return
	}

	err = services.RemoveReview(reviewID, userId)

	if err != nil {
		http.Error(w, "Cannot remove the review", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Review deleted",
	})

}
