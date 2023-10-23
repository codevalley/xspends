package handlers

import (
	"net/http"
	"xspends/models"

	"github.com/gin-gonic/gin"
)

// ListSources retrieves all sources for the authenticated user.
func ListSources(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	sources, err := models.GetSourcesByUserID(userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to fetch sources"})
		return
	}

	c.JSON(http.StatusOK, sources)
}

// GetSource retrieves details of a specific source by its ID.
func GetSource(c *gin.Context) {
	sourceID := c.Param("id")
	source, err := models.GetSourceByID(sourceID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "source not found"})
		return
	}

	c.JSON(http.StatusOK, source)
}

// CreateSource adds a new source for the authenticated user.
func CreateSource(c *gin.Context) {
	var newSource models.Source
	if err := c.ShouldBindJSON(&newSource); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := models.InsertSource(newSource); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to create source"})
		return
	}

	c.JSON(http.StatusOK, newSource)
}

// UpdateSource modifies details of an existing source.
func UpdateSource(c *gin.Context) {
	var updatedSource models.Source
	if err := c.ShouldBindJSON(&updatedSource); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := models.UpdateSource(updatedSource); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to update source"})
		return
	}

	c.JSON(http.StatusOK, updatedSource)
}

// DeleteSource removes a specific source by its ID.
func DeleteSource(c *gin.Context) {
	sourceID := c.Param("id")
	if err := models.DeleteSource(sourceID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to delete source"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "source deleted successfully"})
}
