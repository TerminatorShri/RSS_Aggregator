package main

import (
	"fmt"
	"net/http"

	"github.com/TerminatorShri/RSS_Aggregator/internal/auth"
	"github.com/TerminatorShri/RSS_Aggregator/internal/database"
)

type authedHandler func(http.ResponseWriter, *http.Request, database.User) 

func (apiCfg *apiConfig) middlewareAuth(handler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, 403, fmt.Sprintf("Auth Error: %v", err))
	}

	user, err := apiCfg.DB.GetUserByApiKey(r.Context(), apiKey)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't Get User: %v", err))
		return
	}
	handler(w, r, user)
	}
}