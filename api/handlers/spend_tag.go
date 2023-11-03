package handlers

import (
	"log"
	"net/http"
	"strconv"
	"xspends/models"

	"github.com/gin-gonic/gin"
)

const defaultLimit = 10

func getTagID(c *gin.Context) (int64, bool) {
	tagIDStr := c.Param("id")
	if tagIDStr == "" {
		log.Printf("[getTagID] Error: tag ID is required")
		c.JSON(http.StatusBadRequest, gin.H{"error": "tag ID is required"})
		return 0, false
	}

	tagID, err := strconv.ParseInt(tagIDStr, 10, 64)
	if err != nil {
		log.Printf("[getTagID] Error: invalid tag ID format")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid tag ID format"})
		return 0, false
	}

	return tagID, true
}

func ListTags(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", strconv.Itoa(defaultLimit)))
	offset, _ := strconv.Atoi(c.Query("offset"))

	tags, err := models.GetAllTags(c, userID, models.PaginationParams{Limit: limit, Offset: offset})
	if err != nil {
		log.Printf("[ListTags] Error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to fetch tags"})
		return
	}

	if len(tags) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "no tags found"})
		return
	}

	c.JSON(http.StatusOK, tags)
}

func GetTag(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	tagID, ok := getTagID(c)
	if !ok {
		return
	}

	tag, err := models.GetTagByID(c, tagID, userID)
	if err != nil {
		log.Printf("[GetTag] Error: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "tag not found"})
		return
	}
	c.JSON(http.StatusOK, tag)
}

func CreateTag(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	var newTag models.Tag
	if err := c.ShouldBindJSON(&newTag); err != nil {
		log.Printf("[CreateTag] Error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	newTag.UserID = userID
	if err := models.InsertTag(c, &newTag); err != nil {
		log.Printf("[CreateTag] Error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to create tag"})
		return
	}
	c.JSON(http.StatusOK, newTag)
}

func UpdateTag(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	var updatedTag models.Tag
	if err := c.ShouldBindJSON(&updatedTag); err != nil {
		log.Printf("[UpdateTag] Error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	updatedTag.UserID = userID
	if err := models.UpdateTag(c, &updatedTag); err != nil {
		log.Printf("[UpdateTag] Error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to update tag"})
		return
	}
	c.JSON(http.StatusOK, updatedTag)
}

func DeleteTag(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	tagID, ok := getTagID(c)
	if !ok {
		return
	}

	if err := models.DeleteTag(c, tagID, userID); err != nil {
		log.Printf("[DeleteTag] Error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to delete tag"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "tag deleted successfully"})
}
