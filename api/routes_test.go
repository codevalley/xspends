package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// The health check endpoint should return a 200 status code with a JSON response containing the "status" field set to "UP".
func TestHealthCheckEndpoint(t *testing.T) {
	// Create a mock gin.Engine
	router := gin.Default()

	// Call the code under test
	setupHealthEndpoint(router)

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
