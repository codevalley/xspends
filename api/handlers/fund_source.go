package handlers

import (
	"net/http"
	"strconv"
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

	intUserID, ok := userID.(int)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	sources, err := models.GetSourcesByUserID(intUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to fetch sources"})
		return
	}

	c.JSON(http.StatusOK, sources)
}

func GetSource(c *gin.Context) {
	sourceIDStr := c.Param("id")
	sourceID, err := strconv.Atoi(sourceIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid source ID"})
		return
	}

	source, err := models.GetSourceByID(sourceID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "source not found"})
		return
	}

	c.JSON(http.StatusOK, source)
}

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

func DeleteSource(c *gin.Context) {
	sourceIDStr := c.Param("id")
	sourceID, err := strconv.Atoi(sourceIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid source ID"})
		return
	}

	if err := models.DeleteSource(sourceID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to delete source"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "source deleted successfully"})
}
