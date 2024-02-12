package handlers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"xspends/models/interfaces"
	xmock "xspends/models/mock"
	"xspends/testutils"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func initTagTest(t *testing.T) *xmock.MockTagModel {
	gin.SetMode(gin.TestMode)
	_, modelsService, _, _, tearDown := testutils.SetupModelTestEnvironment(t)
	t.Cleanup(tearDown)

	mockTagModel := new(xmock.MockTagModel)
	modelsService.TagModel = mockTagModel
	return mockTagModel
}

func TestGetTagID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		tagID          string
		expectedID     int64
		expectError    bool
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Valid tag ID",
			tagID:          "123",
			expectedID:     123,
			expectError:    false,
			expectedStatus: 200,
			expectedBody:   "",
		},
		{
			name:           "Invalid tag ID format",
			tagID:          "abc",
			expectError:    true,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "invalid tag ID format",
		},
		{
			name:           "Missing tag ID",
			tagID:          "",
			expectError:    true,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "tag ID is required",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)
			ctx.Params = gin.Params{gin.Param{Key: "id", Value: tc.tagID}}

			id, ok := getTagID(ctx)

			if tc.expectError {
				assert.False(t, ok)
				assert.Equal(t, tc.expectedStatus, w.Code)
				assert.Contains(t, w.Body.String(), tc.expectedBody)
			} else {
				assert.True(t, ok)
				assert.Equal(t, tc.expectedID, id)
			}
		})
	}
}

// More tests follow the same pattern as in TestGetTagID,
// adapted for each handler's logic and requirements.

func TestListTags(t *testing.T) {
	mockTagModel := initTagTest(t)
	defer mockTagModel.AssertExpectations(t)

	isContext := mock.MatchedBy(func(ctx context.Context) bool { return true })

	tests := []struct {
		name           string
		setupMock      func()
		userID         string
		scopeID        string
		limit          string
		offset         string
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Successful retrieval",
			setupMock: func() {
				limit, _ := strconv.Atoi("10")
				offset, _ := strconv.Atoi("0")
				mockTagModel.On("GetScopedTags", isContext, []int64{1}, interfaces.PaginationParams{Limit: limit, Offset: offset}, mock.AnythingOfType("[]*sql.Tx")).Return([]interfaces.Tag{{ID: 1, ScopeID: 1, Name: "Tag 1"}}, nil).Once()
			},
			userID:         "1",
			scopeID:        "1",
			limit:          "10",
			offset:         "0",
			expectedStatus: http.StatusOK,
			expectedBody:   `[{"tag_id":1,"user_id":0,"name":"Tag 1","scope_id":1,"created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z"}]`,
		},
		{
			name: "No tags found",
			setupMock: func() {
				mockTagModel.On("GetScopedTags", isContext, []int64{1}, interfaces.PaginationParams{Limit: 10, Offset: 0}, mock.AnythingOfType("[]*sql.Tx")).Return([]interfaces.Tag{}, nil).Once()
			},
			userID:         "1",
			scopeID:        "1",
			limit:          "10",
			offset:         "0",
			expectedStatus: http.StatusOK,
			expectedBody:   `{"message":"no tags found"}`,
		},
		{
			name: "Internal server error",
			setupMock: func() {
				mockTagModel.On("GetScopedTags", isContext, []int64{1}, interfaces.PaginationParams{Limit: 10, Offset: 0}, mock.AnythingOfType("[]*sql.Tx")).Return([]interfaces.Tag{}, errors.New("internal server error")).Once()
			},
			userID:         "1",
			scopeID:        "1",
			limit:          "10",
			offset:         "0",
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":"unable to fetch tags"}`,
		},
		{
			name: "Unauthorized access",
			setupMock: func() {
				// No setup needed as the mock isn't called when unauthorized
			},
			userID:         "", // An empty userID to simulate unauthorized access
			scopeID:        "",
			limit:          "10",
			offset:         "0",
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"error":"user not authenticated"}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/dummy-url", nil)
			c, _ := gin.CreateTestContext(w)
			c.Request = r
			c.Params = gin.Params{gin.Param{Key: "user_id", Value: tc.userID}}
			c.Params = gin.Params{gin.Param{Key: "scope_id", Value: tc.userID}}

			if tc.userID != "" {
				userID, _ := strconv.ParseInt(tc.userID, 10, 64)
				scopeID, _ := strconv.ParseInt(tc.scopeID, 10, 64)
				c.Set("userID", userID)
				c.Set("scopeID", scopeID)
			}

			// Set query parameters
			q := c.Request.URL.Query()
			q.Add("limit", tc.limit)
			q.Add("offset", tc.offset)
			c.Request.URL.RawQuery = q.Encode()

			tc.setupMock()
			ListTags(c)
			assert.Equal(t, tc.expectedStatus, w.Code)
			assert.Contains(t, w.Body.String(), tc.expectedBody)
		})
	}
}

func TestGetTag(t *testing.T) {
	mockTagModel := initTagTest(t) // Initialize your test setup here, including mock
	defer mockTagModel.AssertExpectations(t)

	tests := []struct {
		name           string
		setupMock      func()
		userID         string
		scopeID        string
		tagID          string
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Successful retrieval",
			setupMock: func() {
				mockTagModel.On("GetTagByID", mock.AnythingOfType("*gin.Context"), int64(1), []int64{1}, mock.AnythingOfType("[]*sql.Tx")).Return(&interfaces.Tag{ID: 1, Name: "Sample Tag"}, nil).Once()
			},
			userID:         "1",
			scopeID:        "1",
			tagID:          "1",
			expectedStatus: http.StatusOK,
			expectedBody:   `{"scope_id":0,"tag_id":1,"name":"Sample Tag","created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z","user_id":0}`, // Adjusted to include all fields
		},
		{
			name: "Tag not found",
			setupMock: func() {
				mockTagModel.On("GetTagByID", mock.Anything, int64(1), []int64{1}, mock.AnythingOfType("[]*sql.Tx")).Return(&interfaces.Tag{}, errors.New("tag not found")).Once()
			},
			userID:         "1",
			scopeID:        "1",
			tagID:          "1",
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"error":"tag not found"}`,
		},
		{
			name: "Unauthorized access",
			setupMock: func() {
				// No setup needed as the mock isn't called when unauthorized
			},
			userID:         "", // an empty userID to simulate unauthorized access
			scopeID:        "1",
			tagID:          "1",
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"error":"user not authenticated"}`,
		},
		// add other test scenarios as necessary
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create a new HTTP request with the tag ID
			req, _ := http.NewRequest("GET", fmt.Sprintf("/tags/%s", tc.tagID), nil)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = req

			c.Params = gin.Params{
				gin.Param{Key: "id", Value: tc.tagID},
			}

			if tc.userID != "" {
				userID, _ := strconv.ParseInt(tc.userID, 10, 64)
				scopeID, _ := strconv.ParseInt(tc.scopeID, 10, 64)
				c.Set("userID", userID)
				c.Set("scopeID", scopeID)
			}

			tc.setupMock()
			GetTag(c) // Call your handler function

			assert.Equal(t, tc.expectedStatus, w.Code)
			assert.JSONEq(t, tc.expectedBody, w.Body.String())
		})
	}
}

func TestCreateTag(t *testing.T) {
	mockTagModel := initTagTest(t) // Initialize the mock and test environment
	defer mockTagModel.AssertExpectations(t)

	tests := []struct {
		name           string
		setupMock      func()
		userID         string
		scopeID        string
		requestBody    string
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Successful creation",
			setupMock: func() {
				newTag := &interfaces.Tag{UserID: 1, ScopeID: 1, Name: "New Tag"} // Adjust to match expected tag structure
				mockTagModel.On("InsertTag", mock.AnythingOfType("*gin.Context"), newTag, mock.AnythingOfType("[]*sql.Tx")).Return(nil).Once()
			},
			userID:         "1",
			scopeID:        "1",
			requestBody:    `{"name":"New Tag"}`,
			expectedStatus: http.StatusCreated,
			expectedBody:   `{"created_at":"0001-01-01T00:00:00Z", "name":"New Tag", "scope_id":1, "tag_id":0, "updated_at":"0001-01-01T00:00:00Z", "user_id":1}`, // Adapt based on actual response structure
		}, {
			name: "Tag creation fails",
			setupMock: func() {
				newTag := &interfaces.Tag{UserID: 1, ScopeID: 1, Name: "New Tag"}
				mockTagModel.On("InsertTag", mock.AnythingOfType("*gin.Context"), newTag, mock.AnythingOfType("[]*sql.Tx")).Return(errors.New("failed to create tag")).Once()
			},
			userID:         "1",
			scopeID:        "1",
			requestBody:    `{"name":"New Tag"}`,
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error": "unable to create tag"}`,
		},
		{
			name: "Invalid JSON",
			setupMock: func() {
				// No mock setup needed as the handler should return error before reaching the model
			},
			userID:         "1",
			scopeID:        "1",
			requestBody:    `{"name": "Invalid JSON",}`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error": "Invalid tag data JSON"}`,
		},
		{
			name: "Unauthorized access",
			setupMock: func() {
				// No mock setup needed as the handler should return error before reaching the model
			},
			userID:         "", // Assuming empty indicates unauthorized or missing user
			scopeID:        "1",
			requestBody:    `{"name":"New Tag"}`,
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"error": "user not authenticated"}`,
		},
		{
			name: "Empty body",
			setupMock: func() {
				// No mock setup needed as the handler should return error before reaching the model
			},
			userID:         "1",
			scopeID:        "1",
			requestBody:    ``,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error": "Invalid tag data JSON"}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Setup the HTTP request with the provided requestBody
			w := httptest.NewRecorder()
			r := httptest.NewRequest("POST", "/tags", strings.NewReader(tc.requestBody))
			c, _ := gin.CreateTestContext(w)
			c.Request = r

			// Set userID in the context if available
			if tc.userID != "" {
				userID, _ := strconv.ParseInt(tc.userID, 10, 64)
				scopeID, _ := strconv.ParseInt(tc.scopeID, 10, 64)
				c.Set("userID", userID)
				c.Set("scopeID", scopeID)
			}

			// Mock setup
			if tc.setupMock != nil {
				tc.setupMock()
			}

			// Invoke the handler
			CreateTag(c)

			// Assert expectations: Verify that the response status code and body are as expected
			assert.Equal(t, tc.expectedStatus, w.Code)
			assert.JSONEq(t, tc.expectedBody, w.Body.String())
		})
	}
}

func TestUpdateTag(t *testing.T) {
	mockTagModel := initTagTest(t) // Initialize the mock and test environment
	defer mockTagModel.AssertExpectations(t)

	tests := []struct {
		name           string
		setupMock      func()
		userID         string
		scopeID        string
		tagID          string
		requestBody    string
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Tag update fails",
			setupMock: func() {
				// Set up the expected tag and ensure it matches the call in your handler.
				updatedTag := &interfaces.Tag{ID: 0, UserID: 1, ScopeID: 1, Name: "Updated Tag"}
				// Setup the mock expectation with the correct parameters.
				mockTagModel.On("UpdateTag", mock.AnythingOfType("*gin.Context"), updatedTag, mock.AnythingOfType("[]*sql.Tx")).Return(errors.New("failed to update tag")).Once() // simulate failure
			},
			userID:         "1",
			scopeID:        "1",                                // assuming this is the user making the update
			requestBody:    `{"id":1, "name":"Updated Tag"}`,   // make sure this matches what your handler expects
			expectedStatus: http.StatusInternalServerError,     // or whatever is appropriate for a failure in update
			expectedBody:   `{"error":"unable to update tag"}`, // expected error message
		},
		{
			name: "Successful update",
			setupMock: func() {
				// Assuming the ID and UserID are known and correct as 1
				updatedTag := &interfaces.Tag{ID: 0, UserID: 1, ScopeID: 1, Name: "Updated Tag"}
				mockTagModel.On("UpdateTag", mock.AnythingOfType("*gin.Context"), updatedTag, mock.AnythingOfType("[]*sql.Tx")).Return(nil).Once()
			},
			userID:         "1",
			scopeID:        "1",
			requestBody:    `{"name":"Updated Tag"}`,
			expectedStatus: http.StatusOK,
			expectedBody: `{
				"scope_id": 1,
				"tag_id": 0,
				"name": "Updated Tag",
				"created_at": "0001-01-01T00:00:00Z",
				"updated_at": "0001-01-01T00:00:00Z",
				"user_id": 1
			}`, // Adjust based on actual response structure
		},
		{
			name: "Invalid JSON",
			setupMock: func() {
				// No mock setup needed as the handler should return error before reaching the model
			},
			userID:         "1",
			scopeID:        "1",
			tagID:          "1",
			requestBody:    `{"name": "Invalid JSON",}`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error": "Invalid tag data JSON"}`,
		},
		{
			name: "Unauthorized access",
			setupMock: func() {
				// No mock setup needed as the handler should return error before reaching the model
			},
			userID:         "", // Assuming empty indicates unauthorized or missing user
			scopeID:        "1",
			tagID:          "1",
			requestBody:    `{"name":"Updated Tag"}`,
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"error": "user not authenticated"}`,
		},

		// ... potentially more test cases for specific error conditions or edge cases ...
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest("PUT", fmt.Sprintf("/tags/%s", tc.tagID), strings.NewReader(tc.requestBody))
			c, _ := gin.CreateTestContext(w)
			c.Request = r

			// Set userID in the context if available
			if tc.userID != "" && tc.scopeID != "" {
				userID, _ := strconv.ParseInt(tc.userID, 10, 64)
				scopeID, _ := strconv.ParseInt(tc.scopeID, 10, 64)
				c.Set("userID", userID)
				c.Set("scopeID", scopeID)
			}

			// Mock setup
			if tc.setupMock != nil {
				tc.setupMock()
			}

			// Invoke the handler
			UpdateTag(c)

			// Assert expectations: Verify that the response status code and body are as expected
			assert.Equal(t, tc.expectedStatus, w.Code)
			assert.JSONEq(t, tc.expectedBody, w.Body.String())
		})
	}
}

func TestDeleteTag(t *testing.T) {
	mockTagModel := initTagTest(t) // Initialize the mock and test environment
	defer mockTagModel.AssertExpectations(t)

	tests := []struct {
		name           string
		setupMock      func()
		userID         string
		scopeID        string
		tagID          string
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Successful deletion",
			setupMock: func() {
				tagID := int64(1) // Assuming tag ID 1 for deletion
				scopeID := int64(1)
				mockTagModel.On("DeleteTag", mock.AnythingOfType("*gin.Context"), tagID, []int64{scopeID}, mock.AnythingOfType("[]*sql.Tx")).Return(nil).Once()
			},
			userID:         "1",
			tagID:          "1",
			scopeID:        "1",
			expectedStatus: http.StatusOK,
			expectedBody:   `{"message": "tag deleted successfully"}`,
		},
		{
			name: "Tag not found or error in deletion",
			setupMock: func() {
				tagID := int64(1)
				scopeID := int64(1)
				// Simulate an error during deletion
				mockTagModel.On("DeleteTag", mock.AnythingOfType("*gin.Context"), tagID, []int64{scopeID}, mock.AnythingOfType("[]*sql.Tx")).
					Return(errors.New("failed to delete tag")).Once()
			},
			userID:         "1",
			tagID:          "1",
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error": "unable to delete tag"}`,
		},
		// ... additional test cases ...
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Initialize Response Recorder (httptest) and Context (gin)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			// Create a request to pass to our handler.
			req := httptest.NewRequest("DELETE", fmt.Sprintf("/tags/%s", tc.tagID), nil)

			// Assign user ID and request to context
			c.Request = req
			if tc.userID != "" {
				userID, _ := strconv.ParseInt(tc.userID, 10, 64)
				scopeID, _ := strconv.ParseInt(tc.userID, 10, 64)
				c.Set("userID", userID)
				c.Set("scopeID", scopeID)
			}

			c.Params = []gin.Param{{Key: "id", Value: tc.tagID}}

			// Setup mocks
			if tc.setupMock != nil {
				tc.setupMock()
			}

			// Invoke the handler
			DeleteTag(c)

			// Assert expectations
			assert.Equal(t, tc.expectedStatus, w.Code)
			assert.JSONEq(t, tc.expectedBody, w.Body.String())
		})
	}
}

// Similar test functions for DeleteTag follow the same pattern,
// adjusting for each handler's specific logic and requirements, including setting up mock responses,
// creating request bodies, and asserting the expected results.

// Please adapt the test setup and assertions according to the actual implementation and expected behavior of your tag
