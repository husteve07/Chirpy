package handlers

import (
	"net/http"
)

func (cfg *APIConfig) HandlerReset(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits.Store(0)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0"))

	cfg.handlerDeleteAllUsers(w, r)
}
