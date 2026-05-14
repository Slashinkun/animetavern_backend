package utils

import (
	"encoding/json"
	"errors"
	"main/contextkeys"
	"main/models"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

func ValidateAddAnimeBody(body models.RequestAddAnimeBody) error {

	if body.AnimeID <= 0 {
		return errors.New("invalid anime ID")
	}

	return nil
}

func ValidateReviewBody(body models.RequestAddReviewBody) error {

	if body.Rating < 0 || body.Rating > 10 {
		return errors.New("invalid rating")
	}

	if strings.TrimSpace(body.Content) == "" {
		return errors.New("review cannot be empty")
	}

	if len(body.Content) > 2000 {
		return errors.New("review too long")
	}

	return nil
}

func ValidateUpdateEpisodeBody(body models.RequestUpdateEpisode) error {

	if body.Delta == 0 {
		return errors.New("delta cannot be 0")
	}

	if body.Delta > 1 || body.Delta < -1 {
		return errors.New("delta must be -1 or +1")
	}

	return nil
}

func ValidateAnimeStatusBody(body models.RequestUpdateAnimeStatus) error {

	if body.Status != "PLANNING" &&
		body.Status != "WATCHING" &&
		body.Status != "FINISHED" {
		return errors.New("invalid status")
	}

	return nil
}

func GetIntParam(r *http.Request, key string) (int, error) {
	vars := mux.Vars(r)

	value, ok := vars[key]
	if !ok || value == "" {
		return 0, errors.New("missing or invalid param")
	}

	return strconv.Atoi(value)
}

func GetUserID(r *http.Request) (int, bool) {
	userID, ok := r.Context().Value(contextkeys.UserID).(int)
	return userID, ok
}

func DecodeJSON(r *http.Request, dst any) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	return decoder.Decode(dst)
}
