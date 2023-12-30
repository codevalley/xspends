package impl

import (
	"testing"
	xmock "xspends/models/mock"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

// func TestInitModelsService(t *testing.T) {
// 	tearDown := setUp(t)
// 	defer tearDown()
// 	err := InitModelsService()
// 	assert.NoError(t, err)
// 	assert.NotNil(t, ModelsService)
// 	assert.NotNil(t, ModelsService.DBService)
// 	// Add assertions for each service initialized in ModelsService
// 	assert.NotNil(t, ModelsService.CategoryModel)
// 	assert.NotNil(t, ModelsService.SourceModel)
// 	// ... continue for all other services
// }

func TestModelsServiceSingleton(t *testing.T) {
	tearDown := setUp(t)
	defer tearDown()

	// Create a mock configuration for testing
	mockConfig := &ModelsConfig{
		DBService: &DBService{}, // Mock DBService or real one as needed
		// other mock services...
	}

	// Initialize the ModelsService for the first time
	InitModelsService(mockConfig)
	firstInstance := GetModelsService()

	// Reinitialize or simulate app restart with potentially a new config
	anotherMockConfig := &ModelsConfig{
		DBService: &DBService{}, // Mock DBService or real one as needed
		// other mock services...
	}
	InitModelsService(anotherMockConfig) // This should not change the instance
	secondInstance := GetModelsService()

	assert.Equal(t, firstInstance, secondInstance) // Ensure it's still the same instance
}

func TestModelsServiceInitialization(t *testing.T) {
	tearDown := setUp(t)
	defer tearDown()

	// Create mock or real instances for each service
	mockDBService := xmock.NewMockDBExecutor(gomock.NewController(t))
	mockCategoryModel := new(xmock.MockCategoryModel)
	// ... and so on for all services

	// Create a configuration with the mock or real services
	config := &ModelsConfig{
		DBService:     &DBService{Executor: mockDBService},
		CategoryModel: mockCategoryModel,
		// ... rest of the mocks or real services
	}

	// Initialize the ModelsService with the mock configuration
	InitModelsService(config)

	// Fetch the initialized ModelsService
	initializedService := GetModelsService()

	// Assertions to ensure each service in ModelsService is correctly assigned
	// If using mocks, assert against the mock instances
	assert.IsType(t, &DBService{}, initializedService.DBService)
	assert.IsType(t, mockCategoryModel, initializedService.CategoryModel)
	// ... continue for all other services
}
