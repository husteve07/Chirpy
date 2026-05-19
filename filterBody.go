package main
import (
	"strings"
)

func filterBody(body string) string {
	bodySlice := strings.Fields(body)

	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert": {},
		"fornax": {},
	}
	for i, word := range bodySlice {
		if _, found := badWords[strings.ToLower(word)]; found {
			bodySlice[i] = "****"
		}
	}
	return strings.Join(bodySlice, " ")
}