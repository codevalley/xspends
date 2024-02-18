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

	// Importing our data models
	"github.com/gin-gonic/gin"
)

func getUserAndScopes(c *gin.Context, role string) (int64, []int64, bool) {
	//TODO: In case of Group, we need to pass the group scope in array[0]
	userID, okUser := getUser(c)
	if !okUser {
		log.Printf("[getUserAndScopes] Error: %v", "Missing user information")
		return 0, nil, false
	}
	scopes, okScope := getScopes(c, userID, role)
	if !okScope {
		log.Printf("[getUserAndScopes] Error: %v", "Missing scope information")
		return userID, nil, false
	}

	return userID, scopes, true
}

func getUser(c *gin.Context) (int64, bool) {
	userID, exists := c.Get("userID")
	if !exists {
		log.Printf("[getUser] Error: %v", "user not authenticated")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return 0, false
	}

	intUserID, ok := userID.(int64)
	if !ok {
		log.Printf("[getUser] Error: %v", "failed to convert userID to int64")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to convert userID to int64"})
		return 0, false
	}

	return intUserID, true
}

func getScope(c *gin.Context) (int64, bool) {
	scopeID, exists := c.Get("scopeID")
	if !exists {
		log.Printf("[getScope] Error: %v", "missing scope parameter")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing scope parameter"})
		return 0, false
	}

	intScopeID, ok := scopeID.(int64)
	if !ok {
		log.Printf("[getScope] Error: %v", "failed to convert scopeID to int64")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to convert scopeID to int64"})
		return 0, false
	}

	return intScopeID, true
}

func getScopes(c *gin.Context, userID int64, role string) ([]int64, bool) {
	scopeID, exists := getScope(c)
	if !exists {
		log.Printf("[getScopes] Error: %v", "missing scope parameter")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing scope parameter"})
		return nil, false
	}
	scopes := []int64{scopeID}
	scopeList, err := impl.GetModelsService().UserScopeModel.GetUserScopesByRole(c, userID, role)
	if err != nil {
		log.Printf("[getScopes] Error: %v", "unable to fetch related scopes for user")
		c.JSON(http.StatusBadRequest, gin.H{"error": "unable to fetch related scopes for user"})
		return nil, false
	}
	for _, scope := range scopeList {
		scopes = append(scopes, scope.ScopeID)
	}
	return scopes, true
}

func GetUserProfile(c *gin.Context) {
	userID, ok := getUser(c)
	if !ok {
		log.Printf("[GetUserProfile] Error: %v", "Missing user information")
		return
	}

	user, err := impl.GetModelsService().UserModel.GetUserByID(c, userID, nil)
	if err != nil {
		log.Printf("[GetUserProfile] Error: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func UpdateUserProfile(c *gin.Context) {
	userID, ok := getUser(c)
	if !ok {
		log.Printf("[UpdateUserProfile] Error: %v", "Missing user information")
		return
	}

	var updatedUser interfaces.User
	if err := c.ShouldBindJSON(&updatedUser); err != nil {
		log.Printf("[UpdateUserProfile] Error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user json"})
		return
	}

	updatedUser.ID = userID

	if err := impl.GetModelsService().UserModel.UpdateUser(c, &updatedUser, nil); err != nil {
		log.Printf("[UpdateUserProfile] Error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to update user"})
		return
	}

	c.JSON(http.StatusOK, updatedUser)
}

func DeleteUser(c *gin.Context) {
	userID, ok := getUser(c)
	if !ok {
		log.Printf("[DeleteUser] Error: %v", "Missing user information")
		return
	}

	if err := impl.GetModelsService().UserModel.DeleteUser(c, userID, nil); err != nil {
		log.Printf("[DeleteUser] Error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user deleted successfully"})
}
