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

func initCategoryTest(t *testing.T) *xmock.MockCategoryModel {
	gin.SetMode(gin.TestMode)
	_, modelsService, _, _, tearDown := testutils.SetupModelTestEnvironment(t)
	defer tearDown()

	mockCategoryModel := new(xmock.MockCategoryModel)
	modelsService.CategoryModel = mockCategoryModel
	return mockCategoryModel
}

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
	mockCategoryModel := initCategoryTest(t)
	defer mockCategoryModel.AssertExpectations(t)

	isContext := mock.MatchedBy(func(ctx context.Context) bool { return true })

	tests := []struct {
		name           string
		setupMock      func()
		userID         string
		ScopeID        string
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
				mockCategoryModel.On("GetScopedCategories", isContext, page, itemsPerPage, []int64{1}, mock.AnythingOfType("[]*sql.Tx")).Return([]interfaces.Category{{ID: 1, ScopeID: 1, Name: "Category 1"}}, nil).Once()
			},
			userID:         "1",
			ScopeID:        "1",
			page:           "1",
			itemsPerPage:   "10",
			expectedStatus: http.StatusOK,
			expectedBody:   "Category 1",
		},
		{
			name: "No categories found",
			setupMock: func() {
				mockCategoryModel.On("GetScopedCategories", isContext, 1, 10, []int64{1}, mock.AnythingOfType("[]*sql.Tx")).Return([]interfaces.Category{}, nil).Once()
			},
			userID:         "1",
			ScopeID:        "1",
			page:           "1",
			itemsPerPage:   "10",
			expectedStatus: http.StatusOK,
			expectedBody:   "no categories found",
		}, {
			name: "Internal server error",
			setupMock: func() {
				mockCategoryModel.On("GetScopedCategories", isContext, 1, 10, []int64{1}, mock.AnythingOfType("[]*sql.Tx")).
					Return([]interfaces.Category{}, errors.New("internal server error")).Once()
			},
			userID:         "1",
			ScopeID:        "1",
			page:           "1",
			itemsPerPage:   "10",
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "unable to fetch categories",
		},
		{
			name: "Unauthorized access",
			setupMock: func() {
				// No mock setup needed as the handler should return error before reaching the model
			},
			userID:         "", // Assuming 0 indicates unauthorized or missing user
			ScopeID:        "",
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "user not authenticated",
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
				scopeID, _ := strconv.ParseInt(tc.ScopeID, 10, 64)
				c.Set("userID", userID)
				c.Set("scopeID", scopeID)
			}
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
	mockCategoryModel := initCategoryTest(t)
	defer mockCategoryModel.AssertExpectations(t)

	tests := []struct {
		name           string
		setupMock      func()
		userID         string
		scopeID        string
		categoryID     string
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Successful retrieval",
			setupMock: func() {
				mockCategoryModel.On("GetCategoryByID", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("int64"), mock.AnythingOfType("[]int64"), mock.AnythingOfType("[]*sql.Tx")).
					Return(&interfaces.Category{ID: 1, ScopeID: 1, Name: "Category 1"}, nil).Once()
			},
			userID:         "1",
			scopeID:        "1",
			categoryID:     "1",
			expectedStatus: http.StatusOK,
			expectedBody:   "Category 1",
		},
		{
			name: "Category not found",
			setupMock: func() {
				mockCategoryModel.On("GetCategoryByID", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("int64"), mock.AnythingOfType("[]int64"), mock.AnythingOfType("[]*sql.Tx")).
					Return((*interfaces.Category)(nil), errors.New("category not found")).Once()
			},
			userID:         "1",
			scopeID:        "1",
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
			scopeID:        "1",
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
			scopeID:        "1",
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
			scopeID:        "",
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
				gin.Param{Key: "scope_id", Value: tc.scopeID},
				gin.Param{Key: "id", Value: tc.categoryID},
			}

			// Set userID in the context
			if tc.userID != "" {
				userID, _ := strconv.ParseInt(tc.userID, 10, 64)
				scopeID, _ := strconv.ParseInt(tc.scopeID, 10, 64)
				c.Set("userID", userID)
				c.Set("scopeID", scopeID)
			}

			tc.setupMock()
			GetCategory(c)
			assert.Equal(t, tc.expectedStatus, w.Code)
			assert.Contains(t, w.Body.String(), tc.expectedBody)
		})
	}
}

// Create a category with valid input, expect status 201 and the new category object in the response body.
func TestCreateCategory(t *testing.T) {
	// Set up test environment
	mockCategoryModel := initCategoryTest(t)
	defer mockCategoryModel.AssertExpectations(t)
	// Define test cases
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
				mockCategoryModel.On("InsertCategory", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("*interfaces.Category"), mock.AnythingOfType("[]*sql.Tx")).
					Return(nil).Once()
			},
			requestBody:    `{"name": "New Category"}`,
			userID:         "1",
			scopeID:        "1",
			expectedStatus: http.StatusCreated,
			expectedBody:   `{"created_at":"0001-01-01T00:00:00Z", "description":"", "icon":"", "category_id":0, "name":"New Category", "scope_id":1, "updated_at":"0001-01-01T00:00:00Z", "user_id":1}`,
		},
		{
			setupMock: func() {
				// No mock setup needed as the handler should return error before reaching the model
			},
			name:           "Invalid JSON",
			userID:         "1",
			scopeID:        "1",
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
			scopeID:        "",
			requestBody:    `{"name": "New Category"}`,
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"error":"user not authenticated"}`,
		},
		{
			name: "Invalid category payload",
			setupMock: func() {
				mockCategoryModel.On("InsertCategory", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("*interfaces.Category"), mock.AnythingOfType("[]*sql.Tx")).
					Return(errors.New("category creation fails")).Once()
			},
			requestBody:    `{"name": "New Category"}`,
			userID:         "1",
			scopeID:        "1",
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":"unable to create category"}`,
		},
		{
			name: "Empty body",
			setupMock: func() {
				// No mock setup needed as the handler should return error before reaching the model
			},
			userID:         "1", // Assuming 0 indicates unauthorized or missing user
			scopeID:        "1",
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
				scopeID, _ := strconv.ParseInt(tc.scopeID, 10, 64)
				c.Set("userID", userID)
				c.Set("scopeID", scopeID)
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

// Successfully update a category with valid input
func TestUpdateCategory(t *testing.T) {
	// Set up test environment
	mockCategoryModel := initCategoryTest(t)
	defer mockCategoryModel.AssertExpectations(t)

	// Define test cases
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
			name: "Successful update",
			setupMock: func() {
				mockCategoryModel.On("UpdateCategory", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("*interfaces.Category"), mock.AnythingOfType("[]*sql.Tx")).
					Return(nil).Once()
			},
			requestBody:    `{"name": "Updated Category"}`,
			userID:         "1",
			scopeID:        "1",
			expectedStatus: http.StatusOK,
			expectedBody:   `{"created_at":"0001-01-01T00:00:00Z", "description":"", "icon":"", "category_id":0, "name":"Updated Category", "scope_id":1, "updated_at":"0001-01-01T00:00:00Z", "user_id":1}`,
		},
		{
			name: "Invalid request body",
			setupMock: func() {
				// No mock setup needed as the handler should return error before reaching the model
			},
			requestBody:    `{"name": "Updated Category",}`,
			userID:         "1",
			scopeID:        "1",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"invalid JSON"}`,
		},
		{
			name: "Unable to update category",
			setupMock: func() {
				mockCategoryModel.On("UpdateCategory", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("*interfaces.Category"), mock.AnythingOfType("[]*sql.Tx")).
					Return(errors.New("unable to update category")).Once()
			},
			requestBody:    `{"name": "Updated Category"}`,
			userID:         "1",
			scopeID:        "1",
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":"unable to update category"}`,
		},
		{
			name: "User ID not authenticated",
			setupMock: func() {
				// No mock setup needed as the handler should return error before reaching the model
			},
			requestBody:    `{"name": "Updated Category"}`,
			userID:         "",
			scopeID:        "",
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"error":"user not authenticated"}`,
		},
		{
			name: "Category update fails",
			setupMock: func() {
				mockCategoryModel.On("UpdateCategory", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("*interfaces.Category"), mock.AnythingOfType("[]*sql.Tx")).
					Return(errors.New("category update fails")).Once()
			},
			requestBody:    `{"name": "Updated Category"}`,
			userID:         "1",
			scopeID:        "1",
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":"unable to update category"}`,
		},
		//doesn't handle the scenario where the payload is invalid.
	}

	// Run test cases
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest("PUT", "/categories/1", strings.NewReader(tc.requestBody))
			c, _ := gin.CreateTestContext(w)
			c.Request = r
			// Set userID in the context
			if tc.userID != "" {
				userID, _ := strconv.ParseInt(tc.userID, 10, 64)
				scopeID, _ := strconv.ParseInt(tc.scopeID, 10, 64)
				c.Set("userID", userID)
				c.Set("scopeID", scopeID)
			}
			// Invoke code under test
			tc.setupMock()
			UpdateCategory(c)

			// Assert response
			assert.Equal(t, tc.expectedStatus, w.Code)
			assert.JSONEq(t, tc.expectedBody, w.Body.String())
		})
	}
}

// The function successfully deletes a category when valid user and category IDs are provided.
func TestDeleteCategory(t *testing.T) {
	// Set up test environment
	mockCategoryModel := initCategoryTest(t)
	defer mockCategoryModel.AssertExpectations(t)

	// Define test cases
	tests := []struct {
		name           string
		setupMock      func()
		userID         string
		scopeID        string
		categoryID     string
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Successful deletion",
			setupMock: func() {
				mockCategoryModel.On("DeleteCategory", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("int64"), mock.AnythingOfType("[]int64"), mock.AnythingOfType("[]*sql.Tx")).
					Return(nil).Once()
			},
			userID:         "1",
			scopeID:        "1",
			categoryID:     "1",
			expectedStatus: http.StatusOK,
			expectedBody:   `{"message": "category deleted successfully"}`,
		},
		{
			name: "Invalid user ID",
			setupMock: func() {
				// No mock setup needed as the handler should return error before reaching the model
			},
			userID:         "",
			scopeID:        "",
			categoryID:     "1",
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"error": "user not authenticated"}`,
		},
		{
			name: "Invalid category ID",
			setupMock: func() {
				// No mock setup needed as the handler should return error before reaching the model
			},
			userID:         "1",
			scopeID:        "1",
			categoryID:     "",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error": "category ID is required"}`,
		},
		{
			name: "Error during category deletion",
			setupMock: func() {
				mockCategoryModel.On("DeleteCategory", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("int64"), mock.AnythingOfType("[]int64"), mock.AnythingOfType("[]*sql.Tx")).
					Return(errors.New("unable to delete category")).Once()
			},
			userID:         "1",
			scopeID:        "1",
			categoryID:     "1",
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error": "unable to delete category"}`,
		},
	}

	// Run test cases
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest("DELETE", "/categories/"+tc.categoryID, nil)
			c, _ := gin.CreateTestContext(w)
			c.Request = r
			// Set userID in the context
			if tc.userID != "" {
				userID, _ := strconv.ParseInt(tc.userID, 10, 64)
				scopeID, _ := strconv.ParseInt(tc.scopeID, 10, 64)
				c.Set("userID", userID)
				c.Set("scopeID", scopeID)
			}
			// Set user_id and category_id in Params
			c.Params = gin.Params{
				gin.Param{Key: "user_id", Value: tc.userID},
				gin.Param{Key: "id", Value: tc.categoryID},
			}
			// Invoke code under test
			tc.setupMock()
			DeleteCategory(c)

			// Assert response
			assert.Equal(t, tc.expectedStatus, w.Code)
			assert.JSONEq(t, tc.expectedBody, w.Body.String())
		})
	}
}
