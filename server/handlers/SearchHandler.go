package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func Search(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	query := r.URL.Query().Get("query")
	if query == "" {
		http.Error(w, "missing query parameter", http.StatusBadRequest)
		return
	}

	apiURL := "https://api.jikan.moe/v4/anime?q=" + url.QueryEscape(query) + "&limit=10&sfw"

	resp, err := http.Get(apiURL)

	if err != nil {
		fmt.Println("Jikan request failed:", err)
		http.Error(w, "Bad gateway", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == 429 {
		http.Error(w, "Rate limited by Jikan", http.StatusTooManyRequests)
		return
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Jikan status:", resp.StatusCode)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadGateway)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Unable to call Jikan API",
		})
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Error reading Jikan API", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}
