package handlers
import (
	"net/http"
	"github.com/google/uuid"
	"encoding/json"
	"github.com/husteve07/Chirpy/internal/database"
	"github.com/husteve07/Chirpy/internal/auth"
	"time"
)

func (cfg *APIConfig) HandlerAuthenticateUser(w http.ResponseWriter, r *http.Request) {
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

	users, err := cfg.db.GetUserByEmail(r.Context(), req.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get user", err)
		return
	}

	if len(users) == 0 {
		respondWithError(w, http.StatusUnauthorized, "Invalid email or password", nil)
		return
	}

	user := users[0]

	match, err := auth.ComparePasswordAndHash(req.Password, user.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't compare password and hash", err)
		return
	}

	if !match {
		respondWithError(w, http.StatusUnauthorized, "Invalid email or password", nil)
		return
	}

	token, err := auth.MakeJWT(user.ID, cfg.secretKey, cfg.jwtExpiresInSeconds)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create JWT", err)
		return
	}

	refreshToken, err := cfg.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token: auth.MakeRefreshToken(),
		UserID: user.ID,
		ExpiresAt: time.Now().Add(cfg.refreshTokenExpiresInDays),
	})


	respondWithJSON(w, http.StatusOK, struct {
		ID uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email string `json:"email"`
		Token string `json:"token"`
		RefreshToken string `json:"refresh_token"`
		IsChirpyRed bool `json:"is_chirpy_red"`
	}{
		ID: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
		Token: token,
		RefreshToken: refreshToken.Token,
		IsChirpyRed: user.IsChirpyRed,
	})

}

func (cfg *APIConfig) HandlerRefresh(w http.ResponseWriter, r *http.Request) {

	refreshTokenStr, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid Authorization header", err)
		return
	}

	refreshToken, err := cfg.db.GetRefreshToken(r.Context(), refreshTokenStr)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid refresh token", err)
		return
	}

	if refreshToken.ExpiresAt.Before(time.Now()) || refreshToken.RevokedAt.Valid {
		respondWithError(w, http.StatusUnauthorized, "Refresh token expired or revoked", nil)
		return
	}

	newAccessToken, err := auth.MakeJWT(refreshToken.UserID, cfg.secretKey, cfg.jwtExpiresInSeconds)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create new access token", err)
		return
	}

	respondWithJSON(w, http.StatusOK, struct {
		Token string `json:"token"`
	}{
		Token: newAccessToken,
	})
}

func (cfg *APIConfig) HandlerRevokeRefreshToken(w http.ResponseWriter, r *http.Request) {

	tokenStr, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid Authorization header", err)
		return
	}
	_, err = cfg.db.RevokeRefreshToken(r.Context(), tokenStr)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't revoke refresh token", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}