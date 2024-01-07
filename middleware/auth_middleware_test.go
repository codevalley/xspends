package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"xspends/api/handlers"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestEnsureUserID(t *testing.T) {
	// Create a gin router
	router := gin.New()
	router.Use(AuthMiddleware(ab)) // First apply AuthMiddleware
	router.Use(EnsureUserID())     // Then apply EnsureUserID
	router.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "Passed")
	})

	// Generate a valid token using the existing method
	userID := int64(123)
	sessionID := "session123"                                             // Use an appropriate session ID or equivalent
	validToken, _ := handlers.GenerateTokenWithTTL(userID, sessionID, 30) // 30 mins or appropriate duration

	// Test the success scenario
	t.Run("Success", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Bearer "+validToken) // Set the valid token in Authorization header
		router.ServeHTTP(w, req)

		// Assertions
		assert.Equal(t, http.StatusOK, w.Code, "Expected to pass with valid token and userID set")
	})

	// Test the failure scenario
	t.Run("Failure", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		c, _ := gin.CreateTestContext(w)
		// Note: Not setting userID here to simulate failure
		c.Request = req // Associate the request with the context
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusUnauthorized, w.Code, "Expected to fail without userID")
	})
}

func TestAuthMiddlewareOld(t *testing.T) {
	// Create a gin router
	router := gin.New()
	router.Use(AuthMiddleware(ab)) // assuming 'ab' is your initialized *authboss.Authboss
	router.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "Passed")
	})

	// Generate a valid token for testing (you'll need to replace this with your actual token generation)
	validToken, _ := handlers.GenerateTokenWithTTL(123, "session123", 30) // Adjust to use your actual function

	// Test the success scenario
	t.Run("Success", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Bearer "+validToken)
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code, "Expected to pass with valid token")
	})

	// Test failure scenarios here, e.g., missing Authorization header, invalid format, invalid token, etc.
	// Similar structure to the above test cases
	// ...
}

func TestAuthMiddleware(t *testing.T) {
	// Create a gin router
	router := gin.New()
	router.Use(AuthMiddleware(ab)) // assuming 'ab' is your initialized *authboss.Authboss
	router.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "Passed")
	})

	// Generate a valid token for testing
	validToken, _ := handlers.GenerateTokenWithTTL(123, "session123", 30) // Adjust with your actual function

	// Test the success scenario
	t.Run("Success", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Bearer "+validToken)
		c, _ := gin.CreateTestContext(w)
		c.Request = req // Associate the request with the context
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code, "Expected to pass with valid token")
	})

	// Add more tests for failure scenarios as needed
	// ...
}
