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
	"xspends/models/impl"
	"xspends/models/interfaces"

	"github.com/gin-gonic/gin"
)

type GroupObject struct {
	GroupName   string
	Description string
	UserRoles   map[int64]string // Map of userID to role
}

func CreateGroup(c *gin.Context) {
	userID, ok := getUser(c)
	if !ok {
		log.Printf("[CreateGroup] Error: %v", "Missing user information")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing user information"})
		return
	}
	var request GroupObject
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Create a new scope for the group
	scopeID, err := impl.GetModelsService().ScopeModel.CreateScope(c, impl.ScopeTypeGroup)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create scope"})
		return
	}

	// Create the group
	group := interfaces.Group{
		OwnerID:     userID, // Assuming a function to extract userID from context
		ScopeID:     scopeID,
		GroupName:   request.GroupName,
		Description: request.Description,
	}
	if err := impl.GetModelsService().GroupModel.CreateGroup(c, &group, nil); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create group"})
		return
	}

	// Assign roles to users including the owner
	if err := impl.GetModelsService().UserScopeModel.UpsertUserScope(c, userID, scopeID, impl.RoleOwner); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assign roles"})
		return
	}
	for user, role := range request.UserRoles {
		//additional check to ensure user is not assigned the same role twice (or role overwritten wrongly)
		if user == userID {
			log.Printf("[CreateGroup] Warning: %v", "Owner cannot be assigned another role")
			continue
		}
		//if invalid role string, skip
		if role != impl.RoleView && role != impl.RoleWrite {
			log.Printf("[CreateGroup] Warning: %v", "Role can only be view or write")
			c.JSON(http.StatusNotAcceptable, gin.H{"error": "Invalid role: " + role})
			return
		}
		if err := impl.GetModelsService().UserScopeModel.UpsertUserScope(c, user, scopeID, role); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assign roles"})
			return
		}
	}
	//TODO: Evaluate if we should pass scopeID
	c.JSON(http.StatusCreated, group)
}
