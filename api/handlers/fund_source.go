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
