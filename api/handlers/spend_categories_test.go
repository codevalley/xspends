package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"xspends/models/impl"
	"xspends/models/interfaces"
	xmock "xspends/models/mock"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var mockCategoryModel *xmock.MockCategoryModel

func setupModelService() {
	// Initialize DBService if necessary
	dbService := &impl.DBService{} // Replace with actual initialization
	mockCategoryModel = new(xmock.MockCategoryModel)
	// Initialize ModelsServiceContainer if not already initialized
	if impl.ModelsService == nil {
		impl.ModelsService = &impl.ModelsServiceContainer{
			DBService:     dbService,
			CategoryModel: mockCategoryModel,
			// Initialize other services as necessary
		}
	}
}
func TestListCategories(t *testing.T) {

	setupModelService()
	// Initialize mock
	// Setup Gin
	gin.SetMode(gin.TestMode)
	r := gin.New()

	// Define the test case
	t.Run("List Categories", func(t *testing.T) {
		// Setup the mock expectations
		mockCategoryModel.On("GetPagedCategories", mock.Anything, 1, 10, int64(1), nil).Return([]interfaces.Category{}, nil)

		// Perform the test
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/categories?page=1&items_per_page=10", nil)
		r.GET("/categories", ListCategories)
		r.ServeHTTP(w, req)

		// Assert expectations
		assert.Equal(t, http.StatusUnauthorized, w.Code)
		mockCategoryModel.AssertExpectations(t)
	})

	// Additional test cases for other handlers...
}

// Implement similar tests for GetCategory, CreateCategory, UpdateCategory, and DeleteCategory
