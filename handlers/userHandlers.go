package handlers
import (
	"net/http"
	"github.com/google/uuid"
	"time"
	"encoding/json"
	"github.com/husteve07/Chirpy/internal/database"
	"github.com/husteve07/Chirpy/internal/auth"
)

type User struct {
	ID			 	uuid.UUID    `json:"id"`
	CreatedAt time.Time		 `json:"created_at"`
	UpdatedAt time.Time		 `json:"updated_at"`
	Email			string		 	 `json:"email"`
}


func (cfg *APIConfig) HandlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type reqContent struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}

	type returnVal struct {
		ID uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email string `json:"email"`
		IsChirpyRed bool `json:"is_chirpy_red"`
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

	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't hash password", err)
		return
	}

	user, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{
		Email: req.Email,
		HashedPassword: hashedPassword,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, returnVal{
		ID: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
		IsChirpyRed: user.IsChirpyRed,
	})


}



func (cfg *APIConfig) HandlerUpdateUser(w http.ResponseWriter, r *http.Request) {
	type reqContent struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}

	var req reqContent
	err := json.NewDecoder(r.Body).Decode(&req)
	defer r.Body.Close()

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return 
	}

	tokenStr, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid Authorization header", err)
		return
	}

	UserID, err := auth.ValidateJWT(tokenStr, cfg.secretKey)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid token", err)
		return
	}

	updatedUser, err := cfg.db.UpdateUser(r.Context(), database.UpdateUserParams{
		ID: UserID,
		Email: req.Email,
		HashedPassword: req.Password,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't update user", err)
		return
	}

	respondWithJSON(w, http.StatusOK, struct {
	ID uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email string `json:"email"`
	IsChirpyRed bool `json:"is_chirpy_red"`
}{
	ID: updatedUser.ID,
	CreatedAt: updatedUser.CreatedAt,
	UpdatedAt: updatedUser.UpdatedAt,
	Email: updatedUser.Email,
	IsChirpyRed: updatedUser.IsChirpyRed,
}	)
}

func (cfg *APIConfig) handlerDeleteAllUsers(w http.ResponseWriter, r *http.Request) {
	err := cfg.db.DeleteAllUsers(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't delete users", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("All users deleted"))
}
