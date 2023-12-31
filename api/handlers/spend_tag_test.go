package handlers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
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

// TestListTags tests the ListTags handler
func TestListTagsOld(t *testing.T) {
	mockTagModel := initTagTest(t)
	defer mockTagModel.AssertExpectations(t)

	isContext := mock.MatchedBy(func(ctx context.Context) bool { return true })

	tests := []struct {
		name           string
		setupMock      func()
		userID         string
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Successful retrieval",
			setupMock: func() {
				mockTagModel.On("GetAllTags", isContext, int64(1), interfaces.PaginationParams{}, mock.AnythingOfType("[]*sql.Tx")).Return([]interfaces.Tag{{ID: 1, Name: "Tag 1"}}, nil).Once()
			},
			userID:         "1",
			expectedStatus: http.StatusOK,
			expectedBody:   `[{"ID":1,"Name":"Tag 1"}]`,
		},
		{
			name: "Error fetching tags",
			setupMock: func() {
				mockTagModel.On("GetAllTags", isContext, int64(1), interfaces.PaginationParams{}, mock.AnythingOfType("[]*sql.Tx")).Return(nil, errors.New("database error")).Once()
			},
			userID:         "1",
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":"unable to fetch tags"}`,
		},
		{
			name: "Unauthorized access",
			setupMock: func() {
				// No mock setup needed as the handler should return error before reaching the model
			},
			userID:         "", // Assuming empty string indicates unauthorized or missing user
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "", // Adapt based on your actual unauthorized response
		},
		{
			name: "No tags found",
			setupMock: func() {
				mockTagModel.On("GetAllTags", isContext, int64(1), interfaces.PaginationParams{}, mock.AnythingOfType("[]*sql.Tx")).Return([]interfaces.Tag{}, nil).Once()
			},
			userID:         "1",
			expectedStatus: http.StatusOK,
			expectedBody:   `{"message":"no tags found"}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/dummy-url", nil)
			c, _ := gin.CreateTestContext(w)
			c.Request = r
			c.Params = gin.Params{gin.Param{Key: "user_id", Value: tc.userID}}

			if tc.userID != "" {
				userID, _ := strconv.ParseInt(tc.userID, 10, 64)
				c.Set("userID", userID)
			}

			tc.setupMock()
			ListTags(c)
			assert.Equal(t, tc.expectedStatus, w.Code)
			assert.Contains(t, w.Body.String(), tc.expectedBody)
		})
	}
}

func TestListTags(t *testing.T) {
	mockTagModel := initTagTest(t)
	defer mockTagModel.AssertExpectations(t)

	isContext := mock.MatchedBy(func(ctx context.Context) bool { return true })

	tests := []struct {
		name           string
		setupMock      func()
		userID         string
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
				mockTagModel.On("GetAllTags", isContext, int64(1), interfaces.PaginationParams{Limit: limit, Offset: offset}, mock.AnythingOfType("[]*sql.Tx")).Return([]interfaces.Tag{{ID: 1, Name: "Tag 1"}}, nil).Once()
			},
			userID:         "1",
			limit:          "10",
			offset:         "0",
			expectedStatus: http.StatusOK,
			expectedBody:   `[{"id":1,"user_id":0,"name":"Tag 1","created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z"}]`,
		},
		{
			name: "No tags found",
			setupMock: func() {
				mockTagModel.On("GetAllTags", isContext, int64(1), interfaces.PaginationParams{Limit: 10, Offset: 0}, mock.AnythingOfType("[]*sql.Tx")).Return([]interfaces.Tag{}, nil).Once()
			},
			userID:         "1",
			limit:          "10",
			offset:         "0",
			expectedStatus: http.StatusOK,
			expectedBody:   `{"message":"no tags found"}`,
		},
		{
			name: "Internal server error",
			setupMock: func() {
				mockTagModel.On("GetAllTags", isContext, int64(1), interfaces.PaginationParams{Limit: 10, Offset: 0}, mock.AnythingOfType("[]*sql.Tx")).Return([]interfaces.Tag{}, errors.New("internal server error")).Once()
			},
			userID:         "1",
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

			if tc.userID != "" {
				userID, _ := strconv.ParseInt(tc.userID, 10, 64)
				c.Set("userID", userID)
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

// Similar test functions for GetTag, CreateTag, UpdateTag, DeleteTag follow the same pattern,
// adjusting for each handler's specific logic and requirements, including setting up mock responses,
// creating request bodies, and asserting the expected results.

// Please adapt the test setup and assertions according to the actual implementation and expected behavior of your tag
