//write testing cases for the password.go file
package auth
import (
	"testing"
	"github.com/google/uuid"
	"time"
)

func TestHashAndComparePassword(t *testing.T) {
	password := "my_secure_password"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	match, err := ComparePasswordAndHash(password, hash)
	if err != nil {
		t.Fatalf("Failed to compare password and hash: %v", err)
	}
	if !match {
		t.Fatalf("Expected password to match hash, but it did not")
	}

	wrongPassword := "wrong_password"
	match, err = ComparePasswordAndHash(wrongPassword, hash)
	if err != nil {
		t.Fatalf("Failed to compare wrong password and hash: %v", err)
	}
	if match {
		t.Fatalf("Expected wrong password not to match hash, but it did")
	}
}

//testing jwt.go file


func TestMakeAndValidateJWT(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "my_secret_key"
	expiresIn := time.Hour
	token, err := MakeJWT(userID, tokenSecret, expiresIn)
	if err != nil {
		t.Fatalf("Failed to make JWT: %v", err)
	}
	validatedUserID, err := ValidateJWT(token, tokenSecret)
	if err != nil {
		t.Fatalf("Failed to validate JWT: %v", err)
	}
	if validatedUserID != userID {
		t.Fatalf("Expected validated user ID to match original user ID, but it did not")
	}
}

