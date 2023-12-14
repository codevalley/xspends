package handlers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"xspends/models/interfaces"
	"xspends/testutils"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestListSources tests the ListSources handler
func TestListSources(t *testing.T) {
	gin.SetMode(gin.TestMode)

	_, _, _, _, tearDown := testutils.SetupModelTestEnvironment(t)
	defer tearDown()

	mockSourceModel := testutils.MockSourceModel
	defer mockSourceModel.AssertExpectations(t)
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
				mockSourceModel.On("GetSources", isContext, int64(1), mock.AnythingOfType("[]*sql.Tx")).Return([]interfaces.Source{{ID: 1, Name: "Source 1"}}, nil).Once()
			},
			userID:         "1",
			expectedStatus: http.StatusOK,
			expectedBody:   "Source 1",
		},
		{
			name: "Error fetching sources",
			setupMock: func() {
				mockSourceModel.On("GetSources", isContext, int64(1), mock.AnythingOfType("[]*sql.Tx")).Return([]interfaces.Source{}, errors.New("database error")).Once()
			},
			userID:         "1",
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "Unable to fetch sources",
		},
		{
			name: "Unauthorized access",
			setupMock: func() {
				// No mock setup needed as the handler should return error before reaching the model
			},
			userID:         "", // Assuming empty string indicates unauthorized or missing user
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "user not authenticated", // or whatever your actual unauthorized response is
		},
		{
			name: "No sources found",
			setupMock: func() {
				mockSourceModel.On("GetSources", isContext, int64(1), mock.AnythingOfType("[]*sql.Tx")).Return([]interfaces.Source{}, nil).Once()
			},
			userID:         "1",
			expectedStatus: http.StatusOK,
			expectedBody:   "[]", // Assuming an empty JSON array is returned for no sources
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
			ListSources(c)
			assert.Equal(t, tc.expectedStatus, w.Code)
			assert.Contains(t, w.Body.String(), tc.expectedBody)
		})
	}
}

// TestGetSource tests the GetSource handler
