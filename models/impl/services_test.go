package impl

import (
	"testing"

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
	firstInstance := GetModelsService()
	// Reinitialize or simulate app restart
	_ = InitModelsService()
	secondInstance := GetModelsService()

	assert.Equal(t, firstInstance, secondInstance)
}

func TestNewModelsServiceContainer(t *testing.T) {
	tearDown := setUp(t)
	defer tearDown()
	container := NewModelsServiceContainer()
	assert.NotNil(t, container)
	assert.NotNil(t, container.DBService)
	assert.NotNil(t, container.CategoryModel)
	assert.NotNil(t, container.SourceModel)
	// ... continue for all other services
}
