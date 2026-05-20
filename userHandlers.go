package main
import (
	"net/http"
	"github.com/google/uuid"
	"time"
	"encoding/json"
)

type User struct {
	ID			 	uuid.UUID    `json:"id"`
	CreatedAt time.Time		 `json:"created_at"`
	UpdatedAt time.Time		 `json:"updated_at"`
	Email			string		 	 `json:"email"`
}


func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type reqContent struct {
		Email string `json:"email"`
	}

	type returnVal struct {
		ID uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email string `json:"email"`
	}

	var req reqContent
	err := json.NewDecoder(r.Body).Decode(&req)
	defer r.Body.Close()

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return 
	}

	if req.Email == "" {
		respondWithError(w, http.StatusBadRequest, "Email is required", nil)
		return
	}

	user, err := cfg.db.CreateUser(r.Context(), req.Email)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, returnVal{
		ID: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
	})


}

func (cfg *apiConfig) handlerDeleteAllUsers(w http.ResponseWriter, r *http.Request) {
	err := cfg.db.DeleteAllUsers(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't delete users", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("All users deleted"))
}
