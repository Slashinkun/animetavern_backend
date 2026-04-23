package handlers

import (
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

	var currentUserID int
	currentUserID = -1
	if idCtx := r.Context().Value(contextkeys.UserID); idCtx != nil {
		currentUserID = idCtx.(int)
	}
	isUser := (currentUserID == userId)

	fmt.Printf("page id : %d  currentid : %d \n", userId, currentUserID)

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

func AddAnimeToList(w http.ResponseWriter, r *http.Request) {
	var body models.RequestBody

	fmt.Println("REQUEST COOKIES:", r.Cookies())

	userId, ok := r.Context().Value(contextkeys.UserID).(int)

	fmt.Println(userId)

	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}

	exists := services.IsAnimeInUserList(userId, body.AnimeID)
	if exists {
		http.Error(w, "Anime déjà dans la liste", http.StatusConflict)
		return
	}

	err = services.AddAnimeToUser(userId, body.AnimeID)
	if err != nil {
		http.Error(w, "Erreur ajout", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Anime ajouté",
	})

}
