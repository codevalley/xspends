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
	"xspends/models"

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
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category ID format"})
		return 0, false
	}

	return categoryID, true
}

func ListCategories(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	itemsPerPage, _ := strconv.Atoi(c.DefaultQuery("items_per_page", strconv.Itoa(defaultItemsPerPage)))

	categories, err := models.GetPagedCategories(c, page, itemsPerPage, userID)
	if err != nil {
		log.Printf("[ListCategories] Error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to fetch categories"})
		return
	}

	if len(categories) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "no categories found"})
		return
	}

	c.JSON(http.StatusOK, categories)
}

func GetCategory(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	categoryID, ok := getCategoryID(c)
	if !ok {
		return
	}

	category, err := models.GetCategoryByID(c, categoryID, userID)
	if err != nil {
		log.Printf("[GetCategory] Error: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "category not found"})
		return
	}

	c.JSON(http.StatusOK, category)
}

func CreateCategory(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	var newCategory models.Category
	if err := c.ShouldBindJSON(&newCategory); err != nil {
		log.Printf("[CreateCategory] Error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newCategory.UserID = userID
	if err := models.InsertCategory(c, &newCategory); err != nil {
		log.Printf("[CreateCategory] Error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to create category"})
		return
	}

	c.JSON(http.StatusCreated, newCategory)
}

func UpdateCategory(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	var updatedCategory models.Category
	if err := c.ShouldBindJSON(&updatedCategory); err != nil {
		log.Printf("[UpdateCategory] Error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedCategory.UserID = userID
	if err := models.UpdateCategory(c, &updatedCategory); err != nil {
		log.Printf("[UpdateCategory] Error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to update category"})
		return
	}

	c.JSON(http.StatusOK, updatedCategory)
}

func DeleteCategory(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	categoryID, ok := getCategoryID(c)
	if !ok {
		return
	}

	if err := models.DeleteCategory(c, categoryID, userID); err != nil {
		log.Printf("[DeleteCategory] Error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to delete category"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "category deleted successfully"})
}
