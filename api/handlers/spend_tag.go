package handlers

import (
	"net/http"
	"strconv"
	"xspends/models"

	"github.com/gin-gonic/gin"
)

const defaultLimit = 10

// ListTags retrieves all available tags with pagination.
func ListTags(c *gin.Context) {
	// Retrieve the user ID from JWT (set in a middleware).
	userID := c.MustGet("userID").(int64)

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", strconv.Itoa(defaultLimit)))
	offset, _ := strconv.Atoi(c.Query("offset"))

	tags, err := models.GetAllTags(c, userID, models.PaginationParams{Limit: limit, Offset: offset})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to fetch tags"})
		return
	}

	if len(tags) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "no tags found"})
		return
	}

	c.JSON(http.StatusOK, tags)
}

// GetTag fetches details of a specific tag by its ID.
func GetTag(c *gin.Context) {
	userID := c.MustGet("userID").(int64)
	tagIDStr := c.Param("id")
	tagID, err := strconv.ParseInt(tagIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid tag ID"})
		return
	}

	tag, err := models.GetTagByID(c, tagID, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "tag not found"})
		return
	}
	c.JSON(http.StatusOK, tag)
}

// CreateTag adds a new tag.
func CreateTag(c *gin.Context) {
	userID := c.MustGet("userID").(int64)
	var newTag models.Tag
	if err := c.ShouldBindJSON(&newTag); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	newTag.UserID = userID
	if err := models.InsertTag(c, &newTag); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to create tag"})
		return
	}
	c.JSON(http.StatusOK, newTag)
}

// UpdateTag modifies details of an existing tag.
func UpdateTag(c *gin.Context) {
	userID := c.MustGet("userID").(int64)
	var updatedTag models.Tag
	if err := c.ShouldBindJSON(&updatedTag); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	updatedTag.UserID = userID
	if err := models.UpdateTag(c, &updatedTag); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to update tag"})
		return
	}
	c.JSON(http.StatusOK, updatedTag)
}

// DeleteTag removes a specific tag by its ID.
func DeleteTag(c *gin.Context) {
	userID := c.MustGet("userID").(int64)
	tagIDStr := c.Param("id")
	tagID, err := strconv.ParseInt(tagIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid tag ID"})
		return
	}

	if err := models.DeleteTag(c, tagID, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to delete tag"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "tag deleted successfully"})
}
