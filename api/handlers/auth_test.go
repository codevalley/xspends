package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
	ymock "xspends/kvstore/mock"
	"xspends/models/impl"
	"xspends/models/interfaces"
	xmock "xspends/models/mock"
	"xspends/testutils"
	"xspends/util"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/volatiletech/authboss/v3"
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
	scopeID := int64(123)
	sessionID := "session123"
	expiryMins := 30

	expectedKey := string(getJwtKey())
	os.Setenv("JWT_KEY", expectedKey)
	defer os.Unsetenv("JWT_KEY")

	token, err := GenerateTokenWithTTL(userID, scopeID, sessionID, expiryMins)
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

func initAuthTest(t *testing.T) (*xmock.MockUserModel, *impl.UserStorer, *impl.SessionStorer, ymock.MockRawKVClientInterface, func()) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockKVClient := ymock.NewMockRawKVClientInterface(ctrl)
	util.InitializeSnowflake()
	_, modelsService, _, _, tearDown := testutils.SetupModelTestEnvironment(t)

	MockUserModel := new(xmock.MockUserModel)
	modelsService.UserModel = MockUserModel

	MockUserStorer := impl.NewUserStorer() // Assuming you have a mock for UserStorer
	sessionStorer := impl.NewSessionStorer(mockKVClient)
	// // Set up expected behavior for UserStorer.Create

	return MockUserModel, MockUserStorer, sessionStorer, *mockKVClient, tearDown
}

func TestJWTRegisterHandler(t *testing.T) {
	// Initialize and configure your mocks
	mockUserModel, mockUserStorer, sessionStorer, mockKV, tearDown := initAuthTest(t)
	defer tearDown()
	defer mockUserModel.AssertExpectations(t)

	// Mock the UserExists method
	// Use gomock.Any() for the *gin.Context parameter and specific values for other parameters
	// Set up the expected behavior of UserExists method
	mockUserModel.On(
		"UserExists",
		mock.Anything,             // Context, use Anything if the exact value doesn't matter
		"newuser",                 // Username to match the input
		"newuser@example.com",     // Email to match the input
		[]*sql.Tx{(*sql.Tx)(nil)}, // Transaction, assuming it's nil in the call
	).Return(false, nil) // What UserExists should return in this scenario

	mockUserModel.On(
		"InsertUser",
		mock.Anything,                           // To match any context
		mock.AnythingOfType("*interfaces.User"), // To match any *interfaces.User
		[]*sql.Tx{(*sql.Tx)(nil)},               // To match the nil transaction slice
	).Return(nil).Once()

	mockKV.EXPECT().
		Put(
			context.Background(),
			gomock.Any(), // matches any []byte for sessionID
			gomock.Any(), // matches any []byte for state
		).Return(nil)

	// Mock dependencies
	ab := &authboss.Authboss{} // Populate with necessary mock implementation
	ab.Config.Storage.Server = mockUserStorer
	ab.Config.Storage.SessionState = sessionStorer
	// Create test request with a new user's details
	body := `{"username":"newuser","email":"newuser@example.com","password":"password"}`
	req := httptest.NewRequest("POST", "/auth/register", strings.NewReader(body))
	w := httptest.NewRecorder()

	// Create a Gin context from the request
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Call the handler
	handler := JWTRegisterHandler(ab)
	handler(c)

	// Check response
	assert.Equal(t, http.StatusOK, w.Code, "Expected successful registration")
	// Further assertions can be made on the response body, headers, etc.

	// Additional tests cases for user already exists, invalid input, etc.
}

func TestJWTLoginHandler(t *testing.T) {
	// Initialize and configure your mocks
	mockUserModel, mockUserStorer, sessionStorer, mockKV, tearDown := initAuthTest(t)
	defer tearDown()
	defer mockUserModel.AssertExpectations(t) // Ensure all expectations are met

	// Mock the GetUserByUsername method
	// Use mock.Anything() for the context parameter
	// Set up the expected behavior of GetUserByUsername method
	mockUserModel.On(
		"GetUserByUsername",
		mock.Anything,             // To match any context
		"",                        // The email or username used to log in
		[]*sql.Tx{(*sql.Tx)(nil)}, // Transaction, assuming it's nil in the call
	).Return(&interfaces.User{ // Return a mock user
		// Populate the mock user fields as necessary
		Username: "existinguser@example.com",
		Password: "$2a$12$bf4KQvsZflGhJmEMMM3hSu/J0yvqAosHpakT1FbHp0WA1LXdV4crC", // Assuming this matches the hash of "password"
	}, nil).Once()

	mockKV.EXPECT().
		Put(
			context.Background(),
			gomock.Any(), // matches any []byte for sessionID
			gomock.Any(), // matches any []byte for state
		).Return(nil)
	// Mock dependencies
	ab := &authboss.Authboss{} // Populate with necessary mock implementation
	ab.Config.Storage.Server = mockUserStorer
	ab.Config.Storage.SessionState = sessionStorer

	// Create test request with a user's login details
	body := `{"email":"existinguser@example.com","password":"password"}`
	req := httptest.NewRequest("POST", "/auth/login", strings.NewReader(body))
	w := httptest.NewRecorder()

	// Create a Gin context from the request
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Call the handler
	handler := JWTLoginHandler(ab)
	handler(c)

	// Check response
	assert.Equal(t, http.StatusOK, w.Code, "Expected successful login")
	// Further assertions can be made on the response body, headers, etc.

	// Additional test cases for invalid credentials, non-existent user, etc.
}

func TestJWTRefreshHandler(t *testing.T) {
	_, mockUserStorer, sessionStorer, mockKV, tearDown := initAuthTest(t)
	defer tearDown()
	userID := int64(123)      // Example user ID, adjust as necessary
	scopeID := int64(123)     // Example scope ID, adjust as necessary
	sessionID := "session123" // Example session ID, adjust as necessary
	refreshToken, _ := GenerateTokenWithTTL(userID, scopeID, sessionID, refreshTokenExpiryMins)
	// Mock dependencies
	mockKV.EXPECT().
		Get(
			context.Background(),
			gomock.Any(), // To match any []byte for sessionID
		).Return([]byte(refreshToken), nil) // Returning the refresh token

	// Set up the expected behavior of Delete method for KvClient
	mockKV.EXPECT().
		Delete(
			context.Background(),
			gomock.Any(), // To match any []byte for sessionID
		).Return(nil) // Simulate successful deletion
	mockKV.EXPECT().
		Put(
			context.Background(),
			gomock.Any(), // matches any []byte for sessionID
			gomock.Any(), // matches any []byte for state
		).Return(nil)

	ab := &authboss.Authboss{} // Populate with necessary mock implementation
	ab.Config.Storage.Server = mockUserStorer
	ab.Config.Storage.SessionState = sessionStorer
	// Create test request with a valid refresh token
	body := fmt.Sprintf(`{"refresh_token":"%s"}`, refreshToken)
	req := httptest.NewRequest("POST", "/auth/refresh", strings.NewReader(body))
	w := httptest.NewRecorder()

	// Create a Gin context from the request
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Call the handler
	handler := JWTRefreshHandler(ab)
	handler(c)

	// Check response
	assert.Equal(t, http.StatusOK, w.Code, "Expected successful token refresh")
	// Further assertions can be made on the response body, headers, etc.

	// Additional test cases for invalid token, expired token, etc.
}

func TestJWTLogoutHandler(t *testing.T) {
	_, mockUserStorer, sessionStorer, mockKV, tearDown := initAuthTest(t)
	defer tearDown()
	userID := int64(123)      // Example user ID, adjust as necessary
	scopeID := int64(123)     // Example scope ID, adjust as necessary
	sessionID := "session123" // Example session ID, adjust as necessary
	refreshToken, _ := GenerateTokenWithTTL(userID, scopeID, sessionID, refreshTokenExpiryMins)
	// Mock dependencies
	ab := &authboss.Authboss{} // Populate with necessary mock implementation
	ab.Config.Storage.Server = mockUserStorer
	ab.Config.Storage.SessionState = sessionStorer
	// Create test request with a valid refresh token
	body := fmt.Sprintf(`{"refresh_token":"%s"}`, refreshToken)
	req := httptest.NewRequest("POST", "/auth/refresh", strings.NewReader(body))
	w := httptest.NewRecorder()

	mockKV.EXPECT().
		Delete(
			context.Background(),
			gomock.Any(), // To match any []byte for sessionID
		).Return(nil) // Simulate successful deletion
	// Create a Gin context from the request
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Call the handler
	handler := JWTLogoutHandler(ab)
	handler(c)

	// Check response
	assert.Equal(t, http.StatusOK, w.Code, "Expected successful logout")
	// Further assertions can be made on the response body, headers, etc.

	// Additional test cases for invalid token, expired token, etc.
}
