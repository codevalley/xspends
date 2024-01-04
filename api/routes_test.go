package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"xspends/kvstore/mock"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
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

func TestSetupSwaggerHandler(t *testing.T) {
	// Set up environment variable for the test
	originalSwaggerPath := os.Getenv("SWAGGER_JSON_PATH")  // Backup original value if it exists
	os.Setenv("SWAGGER_JSON_PATH", "../docs/swagger.json") // Set to test-specific path
	defer func() {
		if originalSwaggerPath != "" {
			os.Setenv("SWAGGER_JSON_PATH", originalSwaggerPath) // Restore original value
		} else {
			os.Unsetenv("SWAGGER_JSON_PATH") // Clean up
		}
	}()

	gin.SetMode(gin.TestMode)
	router := gin.New()
	setupSwaggerHandler(router) // This now uses the SWAGGER_JSON_PATH env variable

	req, _ := http.NewRequest("GET", "/swagger/doc.json", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert the status code and content type
	assert.Equal(t, http.StatusOK, w.Code, "Expected the status code to be 200 OK")

	// Check if the Content-Type header starts with "application/json"
	contentType := w.Header().Get("Content-Type")
	assert.True(t, strings.HasPrefix(contentType, "application/json"), "Expected the content type to start with application/json")

}

func TestSetSwaggerHost(t *testing.T) {
	// Create a temporary file to mimic the Swagger JSON file
	tmpFile, err := ioutil.TempFile("", "swagger*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name()) // Clean up after the test

	// Write a minimal valid Swagger JSON to the temp file
	swaggerJSON := `{"swagger": "2.0", "info": {"title": "Test API", "version": "1.0.0"}}`
	if _, err := tmpFile.Write([]byte(swaggerJSON)); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	// Test case 1: Successful execution
	t.Run("Success", func(t *testing.T) {
		// Set environment variables for the test
		os.Setenv("SWAGGER_HOST", "testhost")
		os.Setenv("SWAGGER_PORT", "testport")
		defer os.Unsetenv("SWAGGER_HOST")
		defer os.Unsetenv("SWAGGER_PORT")

		// Call the function under test
		result, err := setSwaggerHost(tmpFile.Name())

		// Assertions
		assert.NoError(t, err, "Expected no error from setSwaggerHost")

		// Unmarshal the result and check the "host" field
		var modifiedSwagger map[string]interface{}
		json.Unmarshal(result, &modifiedSwagger)
		assert.Equal(t, "testhost:testport", modifiedSwagger["host"], "Expected host to be set to testhost:testport")
	})

	// Add more test cases as needed for error conditions, missing file, etc.
}

func TestSetupRoutes(t *testing.T) {
	// Create a new Gin engine instance
	r := gin.New()

	// Create a mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a mock kvClient
	mockKVClient := mock.NewMockRawKVClientInterface(ctrl)

	// Setup expected calls on the mock (if any), e.g., if your routes make any calls to the kvClient during setup

	// Call SetupRoutes with the test engine and mock client
	SetupRoutes(r, mockKVClient)

	// After setting up routes, you will want to check that the routes are correctly set up.
	// This involves checking if the paths, methods, and handlers are correctly configured.
	// However, directly comparing handler functions in Go isn't straightforward or advisable due to function pointer equality issues.

	// Instead, you might want to check for the existence of expected routes and their methods.
	// For a more detailed test, you would integrate with the handlers and test end-to-end functionality.

	// Example: Check if a specific route exists
	expectedRoutes := []string{"/auth/register", "/auth/login", "/auth/refresh", "/auth/logout", "/sources"}
	for _, route := range expectedRoutes {
		found := false
		for _, info := range r.Routes() {
			if info.Path == route {
				found = true
				break
			}
		}
		assert.True(t, found, "Expected to find route: "+route)
	}

	// Similarly, you might want to check the methods (GET, POST, etc.) of these routes
	// However, for detailed behavior, consider testing individual handlers with their own unit tests or integration tests.

	// Note: This test assumes that SetupRoutes is only setting up routes and not making any initial calls to the kvClient.
	// If SetupRoutes or any middleware/handlers it uses makes calls to the kvClient methods, you will need to set
	// up expectations and return values for those calls on the mockKVClient.
}
