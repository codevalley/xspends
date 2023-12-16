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

const defaultLimit = 10
const maxTagNameLength = 255

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

// ListTags
// @Summary List all tags
// @Description Get a list of all tags
// @ID list-tags
// @Accept  json
// @Produce  json
// @Param limit query int false "Limit number of tags returned"
// @Param offset query int false "Offset for tags returned"
// @Success 200 {array} impl.Tag
// @Failure 500 {object} map[string]string "Unable to fetch tags"
// @Router /tags [get]
func ListTags(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", strconv.Itoa(defaultLimit)))
	offset, _ := strconv.Atoi(c.Query("offset"))

	tags, err := impl.GetModelsService().TagModel.GetAllTags(c, userID, interfaces.PaginationParams{Limit: limit, Offset: offset}, nil)
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

// GetTag
// @Summary Get a specific tag
// @Description Get a specific tag by its ID
// @ID get-tag
// @Accept  json
// @Produce  json
// @Param id path int true "Tag ID"
// @Success 200 {object} impl.Tag
// @Failure 404 {object} map[string]string "Tag not found"
// @Router /tags/{id} [get]
func GetTag(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	tagID, ok := getTagID(c)
	if !ok {
		return
	}

	tag, err := impl.GetModelsService().TagModel.GetTagByID(c, tagID, userID, nil)
	if err != nil {
		log.Printf("[GetTag] Error: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "tag not found"})
		return
	}
	c.JSON(http.StatusOK, tag)
}

// CreateTag
// @Summary Create a new tag
// @Description Create a new tag with the provided information
// @ID create-tag
// @Accept  json
// @Produce  json
// @Param tag body impl.Tag true "Tag info for creation"
// @Success 201 {object} impl.Tag
// @Failure 400 {object} map[string]string "Invalid tag data"
// @Failure 500 {object} map[string]string "Unable to create tag"
// @Router /tags [post]
func CreateTag(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	var newTag interfaces.Tag
	if err := c.ShouldBindJSON(&newTag); err != nil {
		log.Printf("[CreateTag] Error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	newTag.UserID = userID
	if len(newTag.Name) > maxTagNameLength {
		log.Printf("[CreateTag] Error: tag name exceeds maximum length")
		c.JSON(http.StatusBadRequest, gin.H{"error": "tag name exceeds maximum length"})
		return
	}
	if err := impl.GetModelsService().TagModel.InsertTag(c, &newTag, nil); err != nil {
		log.Printf("[CreateTag] Error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to create tag"})
		return
	}
	c.JSON(http.StatusOK, newTag)
}

// UpdateTag
// @Summary Update a specific tag
// @Description Update a specific tag by its ID
// @ID update-tag
// @Accept  json
// @Produce  json
// @Param id path int true "Tag ID"
// @Param tag body impl.Tag true "Tag info for update"
// @Success 200 {object} impl.Tag
// @Failure 400 {object} map[string]string "Invalid tag data"
// @Failure 500 {object} map[string]string "Unable to update tag"
// @Router /tags/{id} [put]
func UpdateTag(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	var updatedTag interfaces.Tag
	if err := c.ShouldBindJSON(&updatedTag); err != nil {
		log.Printf("[UpdateTag] Error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	updatedTag.UserID = userID
	if len(updatedTag.Name) > maxTagNameLength {
		log.Printf("[UpdateTag] Error: tag name exceeds maximum length")
		c.JSON(http.StatusBadRequest, gin.H{"error": "tag name exceeds maximum length"})
		return
	}
	if err := impl.GetModelsService().TagModel.UpdateTag(c, &updatedTag, nil); err != nil {
		log.Printf("[UpdateTag] Error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to update tag"})
		return
	}
	c.JSON(http.StatusOK, updatedTag)
}

// DeleteTag
// @Summary Delete a specific tag
// @Description Delete a specific tag by its ID
// @ID delete-tag
// @Accept  json
// @Produce  json
// @Param id path int true "Tag ID"
// @Success 200 {object} map[string]string "Message: Tag deleted successfully"
// @Failure 500 {object} map[string]string "Unable to delete tag"
// @Router /tags/{id} [delete]
func DeleteTag(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	tagID, ok := getTagID(c)
	if !ok {
		return
	}

	if err := impl.GetModelsService().TagModel.DeleteTag(c, tagID, userID, nil); err != nil {
		log.Printf("[DeleteTag] Error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to delete tag"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "tag deleted successfully"})
}
