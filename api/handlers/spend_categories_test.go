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
		name           string
		categoryID     string
		expectedID     int64
		expectError    bool
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Valid category ID",
			categoryID:     "123",
			expectedID:     123,
			expectError:    false,
			expectedStatus: 200,
			expectedBody:   "",
		},
		{
			name:           "Invalid category ID format",
			categoryID:     "abc",
			expectError:    true,
			expectedStatus: http.StatusNotFound,
			expectedBody:   "invalid category ID format",
		},
		{
			name:           "Missing category ID",
			categoryID:     "",
			expectError:    true,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "category ID is required",
		},
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
				assert.Equal(t, tc.expectedStatus, w.Code)
				assert.Contains(t, w.Body.String(), tc.expectedBody)
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

func TestGetCategory(t *testing.T) {
	gin.SetMode(gin.TestMode)

	_, _, _, _, tearDown := testutils.SetupModelTestEnvironment(t)
	defer tearDown()

	mockCategoryModel := testutils.MockCategoryModel
	defer mockCategoryModel.AssertExpectations(t)

	tests := []struct {
		name           string
		setupMock      func()
		userID         string
		categoryID     string
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Successful retrieval",
			setupMock: func() {
				mockCategoryModel.On("GetCategoryByID", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("int64"), mock.AnythingOfType("int64"), mock.AnythingOfType("[]*sql.Tx")).
					Return(&interfaces.Category{ID: 1, Name: "Category 1"}, nil).Once()
			},
			userID:         "1",
			categoryID:     "1",
			expectedStatus: http.StatusOK,
			expectedBody:   "Category 1",
		},
		{
			name: "Category not found",
			setupMock: func() {
				mockCategoryModel.On("GetCategoryByID", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("int64"), mock.AnythingOfType("int64"), mock.AnythingOfType("[]*sql.Tx")).
					Return((*interfaces.Category)(nil), errors.New("category not found")).Once()
			},
			userID:         "1",
			categoryID:     "1",
			expectedStatus: http.StatusNotFound,
			expectedBody:   "category not found",
		},
		{
			name: "Invalid category ID",
			setupMock: func() {
				// No mock setup needed as the handler should return error before reaching the model
			},
			userID:         "1",
			categoryID:     "NaN", // Assuming this is an invalid ID
			expectedStatus: http.StatusNotFound,
			expectedBody:   "invalid category ID format",
		},
		{
			name: "category ID is required",
			setupMock: func() {
				// No mock setup needed as the handler should return error before reaching the model
			},
			userID:         "1",
			categoryID:     "", // Assuming this is an invalid ID
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "category ID is required",
		},
		{
			name: "Unauthorized access",
			setupMock: func() {
				// No mock setup needed as the handler should return error before reaching the model
			},
			userID:         "", // Assuming 0 indicates unauthorized or missing user
			categoryID:     "1",
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "user not authenticated",
		},

		// Additional test cases...
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", fmt.Sprintf("/categories/%v", tc.categoryID), nil)
			c, _ := gin.CreateTestContext(w)
			c.Request = r

			// Set user_id and category_id in Params
			c.Params = gin.Params{
				gin.Param{Key: "user_id", Value: tc.userID},
				gin.Param{Key: "id", Value: tc.categoryID},
			}

			// Set userID in the context
			if tc.userID != "" {
				userID, _ := strconv.ParseInt(tc.userID, 10, 64)
				c.Set("userID", userID)
			}

			tc.setupMock()
			GetCategory(c)
			assert.Equal(t, tc.expectedStatus, w.Code)
			assert.Contains(t, w.Body.String(), tc.expectedBody)
		})
	}
}

// Create a category with valid input, expect status 201 and the new category object in the response body.
func TestCreateCategoryWithValidInput(t *testing.T) {
	// Set up test environment
	gin.SetMode(gin.TestMode)
	_, _, _, _, tearDown := testutils.SetupModelTestEnvironment(t)
	defer tearDown()

	// Set up mock dependencies
	mockCategoryModel := testutils.MockCategoryModel
	defer mockCategoryModel.AssertExpectations(t)

	// Define test cases
	tests := []struct {
		name           string
		setupMock      func()
		userID         string
		requestBody    string
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Successful creation",
			setupMock: func() {
				mockCategoryModel.On("InsertCategory", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("*interfaces.Category"), mock.AnythingOfType("[]*sql.Tx")).
					Return(nil).Once()
			},
			requestBody:    `{"name": "New Category"}`,
			userID:         "1",
			expectedStatus: http.StatusCreated,
			expectedBody:   `{"created_at":"0001-01-01T00:00:00Z", "description":"", "icon":"", "id":0, "name":"New Category", "updated_at":"0001-01-01T00:00:00Z", "user_id":1}`,
		},
		{
			setupMock: func() {
				// No mock setup needed as the handler should return error before reaching the model
			},
			name:           "Invalid JSON",
			userID:         "1",
			requestBody:    `{"name": "New Category",}`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"invalid JSON"}`,
		},
		{
			name: "Unauthorized access",
			setupMock: func() {
				// No mock setup needed as the handler should return error before reaching the model
			},
			userID:         "", // Assuming 0 indicates unauthorized or missing user
			requestBody:    `{"name": "New Category"}`,
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"error":"user not authenticated"}`,
		},
		{
			name: "Empty body",
			setupMock: func() {
				// No mock setup needed as the handler should return error before reaching the model
			},
			userID:         "1", // Assuming 0 indicates unauthorized or missing user
			requestBody:    "",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error": "invalid JSON"}`,
		},
		//non-existing user, missing-required fields, invalid field values.
	}

	// Run test cases
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest("POST", "/categories", strings.NewReader(tc.requestBody))
			c, _ := gin.CreateTestContext(w)
			c.Request = r
			// Set userID in the context
			if tc.userID != "" {
				userID, _ := strconv.ParseInt(tc.userID, 10, 64)
				c.Set("userID", userID)
			}
			// Invoke code under test
			tc.setupMock()
			CreateCategory(c)

			// Assert response
			assert.Equal(t, tc.expectedStatus, w.Code)
			assert.JSONEq(t, tc.expectedBody, w.Body.String())
		})
	}
}
