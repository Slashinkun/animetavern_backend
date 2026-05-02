package handlers

import (
	"io"
	"net/http"
	"net/url"
)

func Search(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "Non autorisé", http.StatusMethodNotAllowed)
		return
	}

	query := r.URL.Query().Get("query")
	if query == "" {
		http.Error(w, "Paramètre query manquant", http.StatusBadRequest)
		return
	}

	//fmt.Println(query)

	url := "https://api.jikan.moe/v4/anime?q=" + url.QueryEscape(query) + "&limit=10&sfw"

	//fmt.Println(url)

	resp, err := http.Get(url)

	//fmt.Println(resp)
	if err != nil {
		http.Error(w, "Erreur API externe", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, "API Jikan indisponible", http.StatusBadGateway)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Erreur lecture API Jikan", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}
