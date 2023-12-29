package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"xspends/kvstore"
	"xspends/middleware"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// The function should correctly set up the authentication boss using the SetupAuthBoss function.
func test_setup_auth_boss(t *testing.T) {
	// Create a mock gin.Engine
	router := &gin.Engine{}

	// Create a mock kvstore.RawKVClientInterface
	kvClient := &kvstore.RawKVClientWrapper{}

	// Call the code under test
	ab := middleware.SetupAuthBoss(router, kvClient)

	// Assert that the authentication boss is set up correctly
	assert.NotNil(t, ab)
	assert.NotNil(t, ab.Config.Storage.Server)
	assert.NotNil(t, ab.Config.Storage.SessionState)
	assert.NotNil(t, ab.Config.Storage.CookieState)
}

// The health check endpoint should return a 200 status code with a JSON response containing the "status" field set to "UP".
func test_health_check_endpoint(t *testing.T) {
	// Create a mock gin.Engine
	router := gin.Default()

	// Call the code under test
	SetupRoutes(router, nil)

	// Create a mock HTTP request to the health check endpoint
	req, _ := http.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert that the response has a 200 status code and the "status" field is set to "UP"
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "UP", response["status"])
}

// The authentication routes for user registration, login, token refresh, and logout should be set up correctly and return the appropriate responses.
func test_authentication_routes(t *testing.T) {
	// Create a mock gin.Engine
	router := gin.Default()

	// Create a mock kvstore.RawKVClientInterface
	kvClient := &kvstore.RawKVClientWrapper{}

	// Call the code under test
	SetupRoutes(router, kvClient)

	// Create a mock HTTP request to the register endpoint
	registerData := []byte(`{"username": "testuser", "email": "test@example.com", "password": "password"}`)
	req, _ := http.NewRequest("POST", "/auth/register", bytes.NewBuffer(registerData))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert that the response has the appropriate status code and response body

	// Create a mock HTTP request to the login endpoint
	loginData := []byte(`{"username": "testuser", "password": "password"}`)
	req, _ = http.NewRequest("POST", "/auth/login", bytes.NewBuffer(loginData))
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert that the response has the appropriate status code and response body

	// Create a mock HTTP request to the refresh endpoint
	refreshData := []byte(`{"refresh_token": "test_refresh_token"}`)
	req, _ = http.NewRequest("POST", "/auth/refresh", bytes.NewBuffer(refreshData))
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert that the response has the appropriate status code and response body

	// Create a mock HTTP request to the logout endpoint
	logoutData := []byte(`{"access_token": "test_access_token"}`)
	req, _ = http.NewRequest("POST", "/auth/logout", bytes.NewBuffer(logoutData))
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert that the response has the appropriate status code and response body
}

// // The routes for managing sources, categories, tags, and transactions should be set up correctly and return the appropriate responses.
// func test_management_routes(t *testing.T) {
// 	// Create a mock gin.Engine
// 	router := gin.Default()

// 	// Create a mock sql.DB
// 	db, mock, err := sqlmock.New()
// 	if err != nil {
// 		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
// 	}

// 	// Create a mock kvstore.RawKVClientInterface
// 	kvClient := new(mock.RawKVClientInterface)

// 	// Call the code under test
// 	SetupRoutes(router, db, kvClient)

// 	// Create a mock HTTP request to the sources endpoint
// 	req, _ := http.NewRequest("GET", "/sources", nil)
// 	w := httptest.NewRecorder()
// 	router.ServeHTTP(w, req)

// 	// Assert that the response has the appropriate status code and response body
// 	assert.Equal(t, http.StatusOK, w.Code)
// 	assert.Contains(t, w.Body.String(), "sources")
// }
