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
	userID := c.MustGet("userID").(int64)

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	itemsPerPage, _ := strconv.Atoi(c.DefaultQuery("items_per_page", strconv.Itoa(defaultItemsPerPage)))

	categories, err := models.GetPagedCategories(page, itemsPerPage, userID)
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

// GetCategory fetches details of a specific category by its ID.import "strconv"
func GetCategory(c *gin.Context) {
	userID := c.MustGet("userID").(int64)

	categoryIDStr := c.Param("id")
	if categoryIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "category ID is required"})
		return
	}

	// Convert string to int64
	categoryID, err := strconv.ParseInt(categoryIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category ID format"})
		return
	}

	category, err := models.GetCategoryByID(categoryID, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "category not found"})
		return
	}

	c.JSON(http.StatusOK, category)
}

// CreateCategory adds a new category.
func CreateCategory(c *gin.Context) {
	userID := c.MustGet("userID").(int64)
	var newCategory models.Category
	if err := c.ShouldBindJSON(&newCategory); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	newCategory.UserID = userID
	if err := models.InsertCategory(&newCategory); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to create category"})
		return
	}

	c.JSON(http.StatusCreated, newCategory)
}

// UpdateCategory modifies details of an existing category.
func UpdateCategory(c *gin.Context) {
	userID := c.MustGet("userID").(int64)
	var updatedCategory models.Category
	if err := c.ShouldBindJSON(&updatedCategory); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	updatedCategory.UserID = userID
	// Ensure the category exists before updating
	if err := models.UpdateCategory(&updatedCategory); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to update category"})
		return
	}

	c.JSON(http.StatusOK, updatedCategory)
}

// DeleteCategory removes a specific category by its ID.
func DeleteCategory(c *gin.Context) {
	userID := c.MustGet("userID").(int64)
	categoryIDStr := c.Param("id")
	if categoryIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "category ID is required"})
		return
	}

	// Convert string to int64
	categoryID, err := strconv.ParseInt(categoryIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category ID format"})
		return
	}

	if err := models.DeleteCategory(categoryID, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to delete category"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "category deleted successfully"})
}
