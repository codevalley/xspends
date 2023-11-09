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

// ListSources retrieves all sources for the authenticated user.
func ListSources(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	sources, err := models.GetSources(c, userID)
	if err != nil {
		log.Printf("[ListSources] Error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to fetch sources"})
		return
	}

	c.JSON(http.StatusOK, sources)
}

func GetSource(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	sourceIDStr := c.Param("id")
	sourceID, err := strconv.ParseInt(sourceIDStr, 10, 64)
	if err != nil {
		log.Printf("[GetSource] Error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid source ID"})
		return
	}

	source, err := models.GetSourceByID(c, sourceID, userID)
	if err != nil {
		log.Printf("[GetSource] Error: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Source not found"})
		return
	}

	c.JSON(http.StatusOK, source)
}

func CreateSource(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}
	var newSource models.Source
	if err := c.ShouldBindJSON(&newSource); err != nil {
		log.Printf("[CreateSource] Error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid source data"})
		return
	}
	newSource.UserID = userID
	if err := models.InsertSource(c, &newSource); err != nil {
		log.Printf("[CreateSource] Error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create source"})
		return
	}

	c.JSON(http.StatusOK, newSource)
}

func UpdateSource(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}
	var updatedSource models.Source
	if err := c.ShouldBindJSON(&updatedSource); err != nil {
		log.Printf("[UpdateSource] Error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid source data"})
		return
	}
	updatedSource.UserID = userID
	if err := models.UpdateSource(c, &updatedSource); err != nil {
		log.Printf("[UpdateSource] Error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update source"})
		return
	}

	c.JSON(http.StatusOK, updatedSource)
}

func DeleteSource(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}
	sourceIDStr := c.Param("id")
	sourceID, err := strconv.ParseInt(sourceIDStr, 10, 64)
	if err != nil {
		log.Printf("[DeleteSource] Error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid source ID"})
		return
	}

	if err := models.DeleteSource(c, sourceID, userID); err != nil {
		log.Printf("[DeleteSource] Error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete source"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Source deleted successfully"})
}
