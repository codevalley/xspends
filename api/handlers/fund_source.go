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

	"github.com/gin-gonic/gin"
)

// ListSources retrieves all sources for the authenticated user.
// @Summary List all sources
// @Description Get a list of all sources
// @ID list-sources
// @Accept  json
// @Produce  json
// @Success 200 {array} impl.Source
// @Failure 500 {object} map[string]string "Unable to fetch sources"
// @Router /sources [get]
func ListSources(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	sources, err := impl.GetSources(c, userID, impl.GetDBService())
	if err != nil {
		log.Printf("[ListSources] Error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to fetch sources"})
		return
	}

	c.JSON(http.StatusOK, sources)
}

// @Summary Get a specific source
// @Description Get a specific source by its ID
// @ID get-source
// @Accept  json
// @Produce  json
// @Param id path int true "Source ID"
// @Success 200 {object} impl.Source
// @Failure 400 {object} map[string]string "Invalid source ID"
// @Failure 404 {object} map[string]string "Source not found"
// @Router /sources/{id} [get]
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

	source, err := impl.GetSourceByID(c, sourceID, userID, impl.GetDBService())
	if err != nil {
		log.Printf("[GetSource] Error: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Source not found"})
		return
	}

	c.JSON(http.StatusOK, source)
}

// @Summary Create a new source
// @Description Create a new source with the provided information
// @ID create-source
// @Accept  json
// @Produce  json
// @Param source body impl.Source true "Source info for creation"
// @Success 200 {object} impl.Source
// @Failure 400 {object} map[string]string "Invalid source data"
// @Failure 500 {object} map[string]string "Failed to create source"
// @Router /sources [post]
func CreateSource(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}
	var newSource impl.Source
	if err := c.ShouldBindJSON(&newSource); err != nil {
		log.Printf("[CreateSource] Error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid source data"})
		return
	}
	newSource.UserID = userID
	if err := impl.InsertSource(c, &newSource, impl.GetDBService()); err != nil {
		log.Printf("[CreateSource] Error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create source"})
		return
	}

	c.JSON(http.StatusOK, newSource)
}

// @Summary Update a specific source
// @Description Update a specific source by its ID
// @ID update-source
// @Accept  json
// @Produce  json
// @Param id path int true "Source ID"
// @Param source body impl.Source true "Source info for update"
// @Success 200 {object} impl.Source
// @Failure 400 {object} map[string]string "Invalid source data"
// @Failure 500 {object} map[string]string "Failed to update source"
// @Router /sources/{id} [put]
func UpdateSource(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}
	var updatedSource impl.Source
	if err := c.ShouldBindJSON(&updatedSource); err != nil {
		log.Printf("[UpdateSource] Error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid source data"})
		return
	}
	updatedSource.UserID = userID
	if err := impl.UpdateSource(c, &updatedSource, impl.GetDBService()); err != nil {
		log.Printf("[UpdateSource] Error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update source"})
		return
	}

	c.JSON(http.StatusOK, updatedSource)
}

// DeleteSource
// @Summary Delete a specific source
// @Description Delete a specific source by its ID
// @ID delete-source
// @Accept  json
// @Produce  json
// @Param id path int true "Source ID"
// @Success 200 {object} map[string]string "message: Source deleted successfully"
// @Failure 400 {object} map[string]string "Invalid source ID"
// @Failure 500 {object} map[string]string "Failed to delete source"
// @Router /sources/{id} [delete]
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

	if err := impl.DeleteSource(c, sourceID, userID, impl.GetDBService()); err != nil {
		log.Printf("[DeleteSource] Error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete source"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Source deleted successfully"})
}
