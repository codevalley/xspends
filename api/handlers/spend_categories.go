// Package handlers
/*
MIT License

# Copyright (c) 2023 Narayan Babu

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/
package handlers

import (
	"log"
	"net/http"
	"strconv"
	"xspends/models/impl"
	"xspends/models/interfaces"

	"github.com/gin-gonic/gin"
)

const defaultItemsPerPage = 10

func getCategoryID(c *gin.Context) (int64, bool) {
	categoryIDStr := c.Param("id")
	if categoryIDStr == "" {
		log.Printf("[getCategoryID] Error: category ID is required")
		c.JSON(http.StatusBadRequest, gin.H{"error": "category ID is required"})
		return 0, false
	}

	categoryID, err := strconv.ParseInt(categoryIDStr, 10, 64)
	if err != nil {
		log.Printf("[getCategoryID] Error: invalid category ID format")
		c.JSON(http.StatusNotFound, gin.H{"error": "invalid category ID format"})
		return 0, false
	}

	return categoryID, true
}

// ListCategories
// @Summary List all categories
// @Description Get a list of all categories with optional pagination
// @ID list-categories
// @Accept  json
// @Produce  json
// @Param page query int false "Page number"
// @Param items_per_page query int false "Items per page"
// @Success 200 {array} impl.Category
// @Failure 500 {object} map[string]string "Unable to fetch categories"
// @Router /categories [get]
func ListCategories(c *gin.Context) {
	userInfo, ok := GetScopeInfo(c, impl.RoleView)
	if !ok {
		log.Printf("[ListCategories] Error: %v", "Missing user or scope information")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing user or scope information"})
		return
	}
	useScope := userInfo.ownerScope
	if userInfo.groupID != 0 {
		useScope = userInfo.groupScope
	}
	//TODO: Extract literals like this to constants
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	//TODO: Extract literals like this to constants
	itemsPerPage, _ := strconv.Atoi(c.DefaultQuery("items_per_page", strconv.Itoa(defaultItemsPerPage)))
	categories, err := impl.GetModelsService().CategoryModel.GetScopedCategories(c, page, itemsPerPage, []int64{useScope}, nil)
	if err != nil {
		log.Printf("[ListCategories] Error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to fetch categories"})
		return
	}

	if len(categories) == 0 {
		log.Printf("[ListCategories] Error: %v", "no categories found")
		c.JSON(http.StatusOK, gin.H{"message": "no categories found"})
		return
	}

	c.JSON(http.StatusOK, categories)
}

// GetCategory
// @Summary Get a specific category
// @Description Get a specific category by its ID
// @ID get-category
// @Accept  json
// @Produce  json
// @Param id path int true "Category ID"
// @Success 200 {object} impl.Category
// @Failure 404 {object} map[string]string "Category not found"
// @Router /categories/{id} [get]
func GetCategory(c *gin.Context) {
	userInfo, ok := GetScopeInfo(c, impl.RoleView)
	if !ok {
		log.Printf("[GetCategories] Error: %v", "Missing user or scope information")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing user or scope information"})
		return
	}
	useScope := userInfo.ownerScope
	if userInfo.groupID != 0 {
		useScope = userInfo.groupScope
	}

	categoryID, ok := getCategoryID(c)
	if !ok {
		log.Printf("[GetCategory] Error: %v", "invalid category ID format")
		return
	}

	category, err := impl.GetModelsService().CategoryModel.GetCategoryByID(c, categoryID, []int64{useScope}, nil)
	if err != nil {
		log.Printf("[GetCategory] Error: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "category not found"})
		return
	}

	c.JSON(http.StatusOK, category)
}

// CreateCategory
// @Summary Create a new category
// @Description Create a new category with the provided information
// @ID create-category
// @Accept  json
// @Produce  json
// @Param category body impl.Category true "Category info for creation"
// @Success 201 {object} impl.Category
// @Failure 400 {object} map[string]string "Invalid category data"
// @Failure 500 {object} map[string]string "Unable to create category"
// @Router /categories [post]
func CreateCategory(c *gin.Context) {
	userInfo, ok := GetScopeInfo(c, impl.RoleView)
	if !ok {
		log.Printf("[CreateCategories] Error: %v", "Missing user or scope information")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing user or scope information"})
		return
	}
	useScope := userInfo.ownerScope
	if userInfo.groupID != 0 {
		useScope = userInfo.groupScope
	}

	var newCategory interfaces.Category
	if err := c.ShouldBindJSON(&newCategory); err != nil {
		log.Printf("[CreateCategory] Error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	newCategory.UserID = userInfo.userID
	newCategory.ScopeID = useScope
	if err := impl.GetModelsService().CategoryModel.InsertCategory(c, &newCategory, nil); err != nil {
		log.Printf("[CreateCategory] Error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to create category"})
		return
	}

	c.JSON(http.StatusCreated, newCategory)
}

// UpdateCategory
// @Summary Update a specific category
// @Description Update a specific category by its ID
// @ID update-category
// @Accept  json
// @Produce  json
// @Param id path int true "Category ID"
// @Param category body impl.Category true "Category info for update"
// @Success 200 {object} impl.Category
// @Failure 400 {object} map[string]string "Invalid category data"
// @Failure 500 {object} map[string]string "Unable to update category"
// @Router /categories/{id} [put]
func UpdateCategory(c *gin.Context) {
	userInfo, ok := GetScopeInfo(c, impl.RoleView)
	if !ok {
		log.Printf("[CreateCategories] Error: %v", "Missing user or scope information")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing user or scope information"})
		return
	}
	useScope := userInfo.ownerScope
	if userInfo.groupID != 0 {
		useScope = userInfo.groupScope
	}
	var updatedCategory interfaces.Category
	if err := c.ShouldBindJSON(&updatedCategory); err != nil {
		log.Printf("[UpdateCategory] Error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	// model verifies if the categoryID matches the scope ID and user ID, if not updation fails
	updatedCategory.UserID = userInfo.userID
	updatedCategory.ScopeID = useScope
	if err := impl.GetModelsService().CategoryModel.UpdateCategory(c, &updatedCategory, nil); err != nil {
		log.Printf("[UpdateCategory] Error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to update category"})
		return
	}

	c.JSON(http.StatusOK, updatedCategory)
}

// DeleteCategory
// @Summary Delete a specific category
// @Description Delete a specific category by its ID
// @ID delete-category
// @Accept  json
// @Produce  json
// @Param id path int true "Category ID"
// @Success 200 {object} map[string]string "Message: Category deleted successfully"
// @Failure 500 {object} map[string]string "Unable to delete category"
// @Router /categories/{id} [delete]
func DeleteCategory(c *gin.Context) {
	userInfo, ok := GetScopeInfo(c, impl.RoleView)
	if !ok {
		log.Printf("[CreateCategories] Error: %v", "Missing user or scope information")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing user or scope information"})
		return
	}
	useScope := userInfo.ownerScope
	if userInfo.groupID != 0 {
		useScope = userInfo.groupScope
	}

	categoryID, ok := getCategoryID(c)
	if !ok {
		log.Printf("[DeleteCategory] Error: %v", "invalid category ID format")
		return
	}

	if err := impl.GetModelsService().CategoryModel.DeleteCategory(c, categoryID, []int64{useScope}, nil); err != nil {
		log.Printf("[DeleteCategory] Error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to delete category"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "category deleted successfully"})
}
