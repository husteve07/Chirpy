package main

import (
	"net/http"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/husteve07/Chirpy/internal/auth"
	"os"
)

func (cfg *apiConfig) handlerPaymentWebhooks(w http.ResponseWriter, r *http.Request) {

	type reqContent struct {
		Event string `json:"event"`
		Data struct {
			UserID string `json:"user_id"`
		} `json:"data"`
	}

	var req reqContent
	err := json.NewDecoder(r.Body).Decode(&req)
	defer r.Body.Close()

	apiKey := auth.GetAPIKey(r.Header)
	if apiKey == "" {
		respondWithError(w, http.StatusUnauthorized, "Invalid API key", nil)
		return
	}
	if apiKey != os.Getenv("POLKA_KEY") {
		respondWithError(w, http.StatusUnauthorized, "Invalid API key", nil)
		return
	}

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return 
	}

	if req.Event == "user.upgraded" {
		_, err = cfg.db.UpgradeUserMembership(r.Context(), uuid.MustParse(req.Data.UserID))
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't upgrade user to Chirpy Red", err)
			return
		}
	}

	respondWithJSON(w, http.StatusNoContent, nil)

}

