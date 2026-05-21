package auth
import (
	"net/http"
)

func GetAPIKey(headers http.Header) string {

	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return ""
	}
	const prefix = "ApiKey "
	if len(authHeader) <= len(prefix) || authHeader[:len(prefix)] != prefix {
		return ""
	}
	return authHeader[len(prefix):]

}