package models

import (
	"encoding/json"
)

type User struct {
	id       int
	username string
	password string
}

type AnimeResponse struct {
	Data       json.RawMessage `json:"data"`
	IsInList   bool            `json:"isInList,omitempty"`
	IsFavorite bool            `json:"isFavorite,omitempty"`
}
