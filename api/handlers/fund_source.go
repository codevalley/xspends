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
	scopeInfo, ok := c.Get("scopeInfo")
	if !ok {
		log.Printf("[ListSources] Error: %v", "Missing user or scope information")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing user or scope information"})
		return
	}
	userInfo, ok := scopeInfo.(ScopeInfo)
	if !ok {
		log.Printf("[ListSources] Error: %v", "Failed to typecast scope information")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to typecast scope information"})
		return
	}
	useScope := []int64{userInfo.OwnerScope}
	if userInfo.GroupID != 0 {
		useScope = append(useScope, userInfo.GroupScope)
	}
	sources, err := impl.GetModelsService().SourceModel.GetSources(c, useScope)
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
	scopeInfo, ok := c.Get("scopeInfo")
	if !ok {
		log.Printf("[ListSources] Error: %v", "Missing user or scope information")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing user or scope information"})
		return
	}
	userInfo, ok := scopeInfo.(ScopeInfo)
	if !ok {
		log.Printf("[ListSources] Error: %v", "Failed to typecast scope information")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to typecast scope information"})
		return
	}
	useScope := []int64{userInfo.OwnerScope}
	if userInfo.GroupID != 0 {
		useScope = append(useScope, userInfo.GroupScope)
	}
	sourceID, ok := getSourceID(c)
	if !ok {
		log.Printf("[ListSources]: Missing source ID")
		return
	}
	source, err := impl.GetModelsService().SourceModel.GetSourceByID(c, sourceID, useScope)
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
	scopeInfo, ok := c.Get("scopeInfo")
	if !ok {
		log.Printf("[ListSources] Error: %v", "Missing user or scope information")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing user or scope information"})
		return
	}
	userInfo, ok := scopeInfo.(ScopeInfo)
	if !ok {
		log.Printf("[ListSources] Error: %v", "Failed to typecast scope information")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to typecast scope information"})
		return
	}

	var newSource interfaces.Source
	if err := c.ShouldBindJSON(&newSource); err != nil {
		log.Printf("[CreateSource] Error: %v", "Invalid source data")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid source data"})
		return
	}
	newSource.UserID = userInfo.UserID
	newSource.ScopeID = userInfo.UseScope
	if err := impl.GetModelsService().SourceModel.InsertSource(c, &newSource); err != nil {
		log.Printf("[CreateSource] Error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create source"})
		return
	}

	c.JSON(http.StatusCreated, newSource)
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
	scopeInfo, ok := c.Get("scopeInfo")
	if !ok {
		log.Printf("[ListSources] Error: %v", "Missing user or scope information")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing user or scope information"})
		return
	}
	userInfo, ok := scopeInfo.(ScopeInfo)
	if !ok {
		log.Printf("[ListSources] Error: %v", "Failed to typecast scope information")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to typecast scope information"})
		return
	}

	var updatedSource interfaces.Source
	if err := c.ShouldBindJSON(&updatedSource); err != nil {
		log.Printf("[UpdateSource] Error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid source data"})
		return
	}

	// model verifies if the sourceID matches the scope ID and user ID, if not updation fails
	updatedSource.UserID = userInfo.UserID
	updatedSource.ScopeID = userInfo.UseScope
	if err := impl.GetModelsService().SourceModel.UpdateSource(c, &updatedSource); err != nil {
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
	scopeInfo, ok := c.Get("scopeInfo")
	if !ok {
		log.Printf("[ListSources] Error: %v", "Missing user or scope information")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing user or scope information"})
		return
	}
	userInfo, ok := scopeInfo.(ScopeInfo)
	if !ok {
		log.Printf("[ListSources] Error: %v", "Failed to typecast scope information")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to typecast scope information"})
		return
	}
	useScope := []int64{userInfo.OwnerScope}
	if userInfo.GroupID != 0 {
		useScope = append(useScope, userInfo.GroupScope)
	}
	sourceID, ok := getSourceID(c)
	if !ok {
		log.Printf("[DeleteSource]: Missing source ID")
		return
	}
	if err := impl.GetModelsService().SourceModel.DeleteSource(c, sourceID, useScope); err != nil {
		log.Printf("[DeleteSource] Error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete source"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Source deleted successfully"})
}

func getSourceID(c *gin.Context) (int64, bool) {
	sourceIDStr := c.Param("id")
	if sourceIDStr == "" {
		log.Printf("[getSourceID] Error: source ID is required")
		c.JSON(http.StatusBadRequest, gin.H{"error": "source ID is required"})
		return 0, false
	}

	sourceID, err := strconv.ParseInt(sourceIDStr, 10, 64)
	if err != nil {
		log.Printf("[getSourceID] Error: invalid source ID format")
		c.JSON(http.StatusNotFound, gin.H{"error": "invalid source ID format"})
		return 0, false
	}

	return sourceID, true
}
