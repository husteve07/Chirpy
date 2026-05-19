package main
import (
	"net/http"
	"encoding/json"
)

func handlerValidateChirp(w http.ResponseWriter, r *http.Request)  {
	type Chirp struct {
		Body string `json:"body"`
	}

	type returnVal struct {
		CleanedBody string `json:"cleaned_body"`
	}

	var chirp Chirp
	err := json.NewDecoder(r.Body).Decode(&chirp)

	defer r.Body.Close()

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return 
	}

	const maxChirpLength = 140
	if len(chirp.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	cleanedBody := filterBody(chirp.Body)

	respondWithJSON(w, http.StatusOK, returnVal{CleanedBody: cleanedBody})

}