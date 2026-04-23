package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"main/contextkeys"
	"main/models"
	"main/services"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"golang.org/x/time/rate"
)

var jikan_url = "https://api.jikan.moe/v4"
var jikanLimiter = rate.NewLimiter(3, 1)

func FetchAnime(id int) ([]byte, error) {

	err := jikanLimiter.Wait(context.Background())
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/anime/%d", jikan_url, id)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("erreur requete api Jikan: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("jikan error: status %d", resp.StatusCode)
	}

	//fmt.Println("STATUS:", resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("erreur lecture reponse Jikan: %w", err)
	}

	return body, nil
}

// page anime
func GetAnimePage(w http.ResponseWriter, r *http.Request) {

	//variable pour l'affichage des boutons
	var inList bool = false
	var inFavorite bool = false

	if r.Method != http.MethodGet {
		http.Error(w, "Non autorisé", http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)

	animeIdStr := vars["anime_id"]

	//fmt.Println(animeIdStr)

	if animeIdStr == "" {
		http.Error(w, "ID manquant", http.StatusBadRequest)
		return
	}

	animeId, err := strconv.Atoi(animeIdStr)

	if err != nil {
		http.Error(w, "id non conforme", http.StatusBadRequest)
		return
	}

	animedata, err := FetchAnime(animeId)

	//fmt.Println(string(animedata))

	//on transforme la reponse de l'API faciliter la lecture
	var resp models.AnimeJikanResponse
	err = json.Unmarshal(animedata, &resp)
	if err != nil {
		fmt.Println("Erreur JSON:", err)
		http.Error(w, "Erreur décodage API", http.StatusInternalServerError)
		return
	}

	anime := resp.Data

	//on met l'anime dans le cache si il est pas deja dans le cache
	services.PutAnimeInCache(animeId, anime.Title, anime.Images.JPG.ImageURL, anime.Episodes)

	userId := -1

	if v := r.Context().Value(contextkeys.UserID); v != nil {
		if id, ok := v.(int); ok {
			userId = id
		}
	}
	fmt.Printf("userId: %d\n", userId)

	inList = services.IsAnimeInUserList(userId, animeId)

	//on génere la réponse
	data := models.AnimeResponse{
		Anime:      animedata,
		IsInList:   inList,
		IsFavorite: inFavorite,
	}

	//on la transforme en JSON
	response, err := json.Marshal(data)

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)

}
