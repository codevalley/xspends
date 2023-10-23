package handlers

import (
	"net/http"
	"strconv"
	"xspends/models"

	"github.com/gin-gonic/gin"
)

const defaultItemsPerPage = 10

// ListCategories retrieves all available categories with pagination.
func ListCategories(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	itemsPerPage, _ := strconv.Atoi(c.DefaultQuery("items_per_page", strconv.Itoa(defaultItemsPerPage)))

	categories, err := models.GetPagedCategories(page, itemsPerPage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to fetch categories"})
		return
	}

	if len(categories) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "no categories found"})
		return
	}

	c.JSON(http.StatusOK, categories)
}

// GetCategory fetches details of a specific category by its ID.
func GetCategory(c *gin.Context) {
	categoryID := c.Param("id")
	if categoryID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "category ID is required"})
		return
	}

	category, err := models.GetCategoryByID(categoryID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "category not found"})
		return
	}

	c.JSON(http.StatusOK, category)
}

// CreateCategory adds a new category.
func CreateCategory(c *gin.Context) {
	var newCategory models.Category
	if err := c.ShouldBindJSON(&newCategory); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := models.InsertCategory(&newCategory); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to create category"})
		return
	}

	c.JSON(http.StatusCreated, newCategory)
}

// UpdateCategory modifies details of an existing category.
func UpdateCategory(c *gin.Context) {
	var updatedCategory models.Category
	if err := c.ShouldBindJSON(&updatedCategory); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Ensure the category exists before updating
	existingCategory, err := models.GetCategoryByID(updatedCategory.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "category not found"})
		return
	}

	// Validate that the user is updating their own category (or has permission)
	if existingCategory.UserID != updatedCategory.UserID {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "permission denied"})
		return
	}

	if err := models.UpdateCategory(&updatedCategory); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to update category"})
		return
	}

	c.JSON(http.StatusOK, updatedCategory)
}

// DeleteCategory removes a specific category by its ID.
func DeleteCategory(c *gin.Context) {
	categoryID := c.Param("id")
	if categoryID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "category ID is required"})
		return
	}

	// Check if category exists before deletion
	_, err := models.GetCategoryByID(categoryID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "category not found"})
		return
	}

	if err := models.DeleteCategory(categoryID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to delete category"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "category deleted successfully"})
}
