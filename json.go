package main 
import (
	"encoding/json"
	"net/http"
	"log"
)

func respondWithJSON(w http.ResponseWriter, statusCode int, payload interface{}) {

	w.Header().Set("Content-Type", "application/json")
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(statusCode)
	w.Write(data)
}

func respondWithError(w http.ResponseWriter, statusCode int, msg string, err error) {
	if err != nil {
		log.Println(err)
	}

	if statusCode > 499 {
		log.Printf("Server error: %s", msg)
	}

	type errorResponse struct {
		Error string `json:"error"`
	}

	respondWithJSON(w, statusCode, errorResponse{Error: msg})
}