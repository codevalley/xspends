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
	GroupName   string           `json:"group_name"`
	Description string           `json:"description"`
	UserRoles   map[int64]string `json:"user_roles"`
}

// AddToGroupRequest represents the request payload for adding a user to a group.
type AddToGroupRequest struct {
	GroupID int64  `json:"group_id"`
	UserID  int64  `json:"user_id"`
	Role    string `json:"role"`
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

func AddToGroup(c *gin.Context) {
	// Step 1: Authenticate and get current userID
	currentUserID, ok := getUser(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing user information"})
		return
	}

	// Step 2: Fetch the request payload
	var request struct {
		GroupID int64  `json:"groupID"`
		UserID  int64  `json:"userID"`
		Role    string `json:"role"`
	}
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Step 3: Verify if the current user is the owner of the requested GroupID
	group, err := impl.GetModelsService().GroupModel.GetGroupByID(c, request.GroupID, currentUserID)
	if err != nil || group.OwnerID != currentUserID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized to add members to this group"})
		return
	}

	// Step 4: Validate role type
	if request.Role != impl.RoleView && request.Role != impl.RoleWrite {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role specified"})
		return
	}

	// Step 5: Add the userID tuple to the userScope table
	if err := impl.GetModelsService().UserScopeModel.UpsertUserScope(c, request.UserID, group.ScopeID, request.Role); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add user to group"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User added to group successfully"})
}

func RemoveFromGroup(c *gin.Context) {
	// Step 1: Authenticate and get current userID
	currentUserID, ok := getUser(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing user information"})
		return
	}

	// Step 2: Fetch the request payload
	var request struct {
		GroupID int64 `json:"groupID"`
		UserID  int64 `json:"userID"`
	}
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Step 3: Verify if the current user is the owner of the requested GroupID
	group, err := impl.GetModelsService().GroupModel.GetGroupByID(c, request.GroupID, currentUserID)
	if err != nil || group.OwnerID != currentUserID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized to remove members from this group"})
		return
	}

	// Additional step: Ensure the user to be removed exists within the group
	_, err = impl.GetModelsService().UserScopeModel.GetUserScope(c, request.UserID, group.ScopeID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not part of the group"})
		return
	}

	// Step 4: Remove the userID tuple from the userScope table
	if err := impl.GetModelsService().UserScopeModel.DeleteUserScope(c, request.UserID, group.ScopeID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove user from group"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User removed from group successfully"})
}
