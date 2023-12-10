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

// TestGetCategoryID tests the getCategoryID function
func TestGetCategoryID(t *testing.T) {
	// Set Gin to Test Mode
	gin.SetMode(gin.TestMode)

	// Define test cases
	tests := []struct {
		name        string
		categoryID  string
		expectedID  int64
		expectError bool
	}{
		// ... [your test cases here]
	}

	// Run test cases
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Set up a response recorder and context
			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)
			ctx.Params = gin.Params{
				gin.Param{Key: "id", Value: tc.categoryID},
			}

			// Call the function
			id, err := getCategoryID(ctx)

			// Assert expectations
			if tc.expectError {
				assert.False(t, err, "Expected an error but didn't get one")
			} else {
				assert.True(t, err, "Expected no error but got one")
				assert.Equal(t, tc.expectedID, id)
			}
		})
	}
}

// TestListCategories tests the ListCategories handler
func TestListCategories(t *testing.T) {
	gin.SetMode(gin.TestMode)

	_, _, _, _, tearDown := testutils.SetupModelTestEnvironment(t)
	defer tearDown()

	mockCategoryModel := testutils.MockCategoryModel
	defer mockCategoryModel.AssertExpectations(t)
	isContext := mock.MatchedBy(func(ctx context.Context) bool { return true })

	tests := []struct {
		name           string
		setupMock      func()
		userID         int64
		page           string
		itemsPerPage   string
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Successful retrieval",
			setupMock: func() {
				page, _ := strconv.Atoi("1")
				itemsPerPage, _ := strconv.Atoi("10")
				mockCategoryModel.On("GetPagedCategories", isContext, page, itemsPerPage, int64(1), mock.AnythingOfType("[]*sql.Tx")).Return([]interfaces.Category{{ID: 1, Name: "Category 1"}}, nil).Once()
			},
			userID:         1,
			page:           "1",
			itemsPerPage:   "10",
			expectedStatus: http.StatusOK,
			expectedBody:   "Category 1",
		},
		{
			name: "No categories found",
			setupMock: func() {
				mockCategoryModel.On("GetPagedCategories", isContext, 1, 10, int64(1), mock.AnythingOfType("[]*sql.Tx")).Return([]interfaces.Category{}, nil).Once()
			},
			userID:         1,
			page:           "1",
			itemsPerPage:   "10",
			expectedStatus: http.StatusOK,
			expectedBody:   "no categories found",
		}, {
			name: "Internal server error",
			setupMock: func() {
				mockCategoryModel.On("GetPagedCategories", isContext, 1, 10, int64(1), mock.AnythingOfType("[]*sql.Tx")).
					Return([]interfaces.Category{}, errors.New("internal server error")).Once()
			},
			userID:         1,
			page:           "1",
			itemsPerPage:   "10",
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "unable to fetch categories",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/dummy-url", nil)
			c, _ := gin.CreateTestContext(w)
			c.Request = r
			c.Params = gin.Params{gin.Param{Key: "user_id", Value: strconv.FormatInt(tc.userID, 10)}}
			c.Set("userID", tc.userID)

			q := c.Request.URL.Query()
			q.Add("page", tc.page)
			q.Add("items_per_page", tc.itemsPerPage)
			c.Request.URL.RawQuery = q.Encode()

			tc.setupMock()
			ListCategories(c)
			assert.Equal(t, tc.expectedStatus, w.Code)
			assert.Contains(t, w.Body.String(), tc.expectedBody)
		})
	}
}
