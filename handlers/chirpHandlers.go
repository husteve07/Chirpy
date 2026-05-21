package handlers
import (
	"net/http"
	"github.com/google/uuid"
	"encoding/json"
	"github.com/husteve07/Chirpy/internal/database"
	"github.com/husteve07/Chirpy/internal/auth"
)



func (cfg *APIConfig) HandlerCreateChirp(w http.ResponseWriter, r *http.Request)  {

	type reqContent struct {
		Body string `json:"body"`
	}

	type returnVal struct {
		ID uuid.UUID `json:"id"`
		Body string `json:"body"`
		UserID  uuid.UUID `json:"user_id"`
	}

	var req reqContent
	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return 
	}

	if req.Body == "" {
		respondWithError(w, http.StatusBadRequest, "Body is required", nil)
		return
	}

	cleanedBody := filterBody(req.Body)

	defer r.Body.Close()

	tokenStr, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Missing or invalid Authorization header", err)
		return
	}
	UserID, err := auth.ValidateJWT(tokenStr, cfg.secretKey)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid token", err)
		return
	}

	chirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body: cleanedBody,
		UserID:  UserID,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create chirp", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, returnVal{
		ID: chirp.ID,
		Body: chirp.Body,
		UserID: chirp.UserID,
	})


}

func (cfg *APIConfig) HandlerGetAllChirps(w http.ResponseWriter, r *http.Request)  {

	type Chirp struct {
		ID uuid.UUID `json:"id"`
		Body string				 `json:"body"` 
		UserID  uuid.UUID `json:"user_id"`
	}

	var returnChirps []Chirp

	queryParams := r.URL.Query()
	authorID, ok := queryParams["author_id"]
	sortOrder, sortOrderExists := queryParams["sort"]

	if ok && len(authorID) > 0 {
		var parsedAuthorID uuid.UUID
		parsedAuthorID, err := uuid.Parse(authorID[0])
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid author_id parameter", err)
			return
		}
		chirps, err := cfg.db.GetChirpsByAuthorIDs(r.Context(), parsedAuthorID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't get chirps", err)
			return
		}

		for _, chirp := range chirps {
			returnChirps = append(returnChirps, Chirp{
				ID: chirp.ID,
				Body: chirp.Body,
				UserID: chirp.UserID,
			})
		}
		respondWithJSON(w, http.StatusOK, returnChirps)
		return
	}

	if len(sortOrder) > 0 {
		if !sortOrderExists || sortOrder[0] == "asc"  {
			chirps, err := cfg.db.GetChirpsSortedByCreatedAtAsc(r.Context())
			if err != nil {
				respondWithError(w, http.StatusInternalServerError, "Couldn't get chirps", err)
				return
			}
			for _, chirp := range chirps {
				returnChirps = append(returnChirps, Chirp{
					ID: chirp.ID,
					Body: chirp.Body,
					UserID: chirp.UserID,
				})
			}
			respondWithJSON(w, http.StatusOK, returnChirps)
			return
		} else if sortOrder[0] == "desc" {
			chirps, err := cfg.db.GetChirpsSortedByCreatedAtDesc(r.Context())
			if err != nil {
				respondWithError(w, http.StatusInternalServerError, "Couldn't get chirps", err)
				return
			}
			for _, chirp := range chirps {
				returnChirps = append(returnChirps, Chirp{
					ID: chirp.ID,
					Body: chirp.Body,
					UserID: chirp.UserID,
				})
			}
			respondWithJSON(w, http.StatusOK, returnChirps)
			return
		}
	}

	chirps, err := cfg.db.GetAllChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get chirps", err)
		return
	}
	for _, chirp := range chirps {	
		returnChirps = append(returnChirps, Chirp{
			ID: chirp.ID,
			Body: chirp.Body,
			UserID: chirp.UserID,
		})
	}

	respondWithJSON(w, http.StatusOK, returnChirps)
}

func (cfg *APIConfig) HandlerGetChirpByID(w http.ResponseWriter, r *http.Request)  {

	type Chirp struct {
		ID uuid.UUID `json:"id"`
		Body string				 `json:"body"` 
		UserID  uuid.UUID `json:"user_id"`
	}

	chirpIDStr := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(chirpIDStr)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID", err)
		return
	}

	chirps, err := cfg.db.GetAllChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get chirps", err)
		return
	}
	for _, chirp := range chirps {
		if chirp.ID == chirpID {
			respondWithJSON(w, http.StatusOK, Chirp{
				ID: chirp.ID,
				Body: chirp.Body,
				UserID: chirp.UserID,
			})
			return
		}
	}

	respondWithError(w, http.StatusNotFound, "No chirps found for this ID", nil)
}

func (cfg *APIConfig) HandlerDeleteChirp(w http.ResponseWriter, r *http.Request)  {

	tokenStr, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Missing or invalid Authorization header", err)
		return
	}
	userID, err := auth.ValidateJWT(tokenStr, cfg.secretKey)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid token", err)
		return
	}

	chirpIDStr := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(chirpIDStr)

	chirp, err := cfg.db.GetChirpByID(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Chirp not found", err)
		return
	}

	if chirp.UserID != userID {
		respondWithError(w, http.StatusForbidden, "You don't have permission to delete this chirp", nil)
		return
	}

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID", err)
		return
	}

	err = cfg.db.DeleteChirp(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't delete chirp", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)

}
