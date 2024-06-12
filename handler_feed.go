package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/TerminatorShri/RSS_Aggregator/internal/database"
	"github.com/google/uuid"
)

func (apiCfg *apiConfig) handlerCreateFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		Name string `json:"name"`
		URL string `'json:"url"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{} 
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Erros Parsing JSON: %v", err))
		return
	}

	feed, err := apiCfg.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		ID: uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name: params.Name,
		Url: params.URL,
		UserID: user.ID,
	})

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Unable to Create User: %v", err))
		return
	}

	respondWithJSON(w, 201, databaseFeedToFeed(feed))
}
func (apiCfg *apiConfig) handlerGetFeeds(w http.ResponseWriter, r *http.Request) {
	feeds, err := apiCfg.DB.GetFeed(r.Context())
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't Get Feeds: %v", err))
		return
	}
	respondWithJSON(w, 201, databaseFeedsToFeeds(feeds))
}

