package handlers

import (
	"os"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func TestGetJwtKey(t *testing.T) {
	// Test case 1: JWT_KEY environment variable is set
	expectedKey := "test_key"
	os.Setenv("JWT_KEY", expectedKey)
	defer os.Unsetenv("JWT_KEY")

	key := getJwtKey()
	if string(key) != expectedKey {
		t.Errorf("Expected key: %s, got: %s", expectedKey, string(key))
	}

	// Test case 2: JWT_KEY environment variable is not set
	os.Unsetenv("JWT_KEY")

	defaultKey := "uNauz8OMH3UzF6wum99OD6dsm1wSdMquDGkWznT6JrQ="
	key = getJwtKey()
	if string(key) != defaultKey {
		t.Errorf("Expected default key: %s, got: %s", defaultKey, string(key))
	}
}
func TestGenerateTokenWithTTL(t *testing.T) {
	userID := int64(123)
	sessionID := "session123"
	expiryMins := 30

	expectedKey := string(getJwtKey())
	os.Setenv("JWT_KEY", expectedKey)
	defer os.Unsetenv("JWT_KEY")

	token, err := generateTokenWithTTL(userID, sessionID, expiryMins)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Verify token
	parsedToken, err := jwt.ParseWithClaims(token, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(expectedKey), nil
	})
	if err != nil {
		t.Errorf("Unexpected error while parsing token: %v", err)
	}

	claims, ok := parsedToken.Claims.(*JWTClaims)
	if !ok {
		t.Errorf("Failed to parse token claims")
	}

	if claims.UserID != userID {
		t.Errorf("Expected UserID: %d, got: %d", userID, claims.UserID)
	}

	if claims.SessionID != sessionID {
		t.Errorf("Expected SessionID: %s, got: %s", sessionID, claims.SessionID)
	}

	// Test case with expired token
	expiredToken := generateExpiredToken(expectedKey)
	parsedExpiredToken, err := jwt.ParseWithClaims(expiredToken, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(expectedKey), nil
	})
	if err == nil {
		t.Errorf("Expected error while parsing expired token, but got no error")
	}

	if parsedExpiredToken.Valid {
		t.Errorf("Expected expired token to fail parsing")
	}
}

func generateExpiredToken(key string) string {
	expirationTime := time.Now().Add(-time.Minute)
	claims := &JWTClaims{
		UserID:    123,
		SessionID: "expired_session",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, _ := token.SignedString([]byte(key))
	return signedToken
}
