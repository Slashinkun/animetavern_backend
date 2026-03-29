package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"main/models"
	"main/services"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

var jikan_url = "https://api.jikan.moe/v4"

func fetchAnime(id int) ([]byte, error) {
	url := fmt.Sprintf("%s/anime/%d", jikan_url, id)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("erreur requete api Jikan: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("erreur lecture reponse Jikan: %w", err)
	}

	return body, nil
}

func GetAnimePage(w http.ResponseWriter, r *http.Request) {
	var inList bool = false
	var inFavorite bool = false
	if r.Method != http.MethodGet {
		http.Error(w, "Non autorisé", http.StatusBadRequest)
		return
	}

	//temporaire pour eviter le pb du CORS
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 3 {
		http.Error(w, "ID manquant", http.StatusBadRequest)
		return
	}
	id := pathParts[2]
	idStr, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "id non conforme", http.StatusBadRequest)
		return
	}
	animedata, err := fetchAnime(idStr)

	//fmt.Println(string(animedata))

	var resp models.AnimeJikanResponse
	err = json.Unmarshal(animedata, &resp)
	if err != nil {
		fmt.Println("Erreur JSON:", err)
		http.Error(w, "Erreur décodage API", http.StatusInternalServerError)
		return
	}

	anime := resp.Data

	services.PutAnimeInCache(idStr, anime.Title, anime.Images.JPG.ImageURL, anime.Episodes)

	if err != nil {
		http.Error(w, "erreur api jikan", http.StatusNotFound)
		return
	}

	username := r.Context().Value("username")

	if username != nil {
		username := username.(string)
		inList = services.IsAnimeInUserList(username, idStr)
		inFavorite = services.IsAnimeUserFavorite(username, idStr)

	}

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

func Search(w http.ResponseWriter, r *http.Request) {

	//temporaire pour eviter le pb du CORS
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	if r.Method != http.MethodGet {
		http.Error(w, "Non autorisé", http.StatusMethodNotAllowed)
		return
	}

	query := r.URL.Query().Get("query")
	if query == "" {
		http.Error(w, "Paramètre query manquant", http.StatusBadRequest)
		return
	}

	resp, err := http.Get("https://api.jikan.moe/v4/anime?q=" + url.QueryEscape(query) + "&limit=10")
	if err != nil {
		http.Error(w, "Erreur API Jikan", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Erreur lecture API Jikan", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}

// POST /register
func RegisterHandler(w http.ResponseWriter, r *http.Request) {

	//temporaire pour eviter le pb du CORS
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

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

	err := services.Register(username, password)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.WriteHeader(http.StatusCreated)

}

// POST /login
func LoginHandler(w http.ResponseWriter, r *http.Request) {

	//temporaire pour eviter le pb du CORS
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	username := r.FormValue("username")
	password := r.FormValue("password")

	token, err := services.Login(username, password)
	if err != nil {
		http.Error(w, err.Error(), 401)
		return
	}

	expirationTime := time.Now().Add(1 * time.Hour)

	//ecriture du cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		Expires:  expirationTime,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})

}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {

	//temporaire pour eviter le pb du CORS
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	//on change le temps du cookie pour qu'il disparaisse immediatement
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Déconnecté"))
}
