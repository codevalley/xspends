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
	"net/http"
	"xspends/models/impl"
	"xspends/models/interfaces"

	// Importing our data models
	"github.com/gin-gonic/gin"
)

func getUserID(c *gin.Context) (int64, bool) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return 0, false
	}

	intUserID, ok := userID.(int64)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to convert userID to int64"})
		return 0, false
	}

	return intUserID, true
}

func GetUserProfile(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	user, err := impl.GetModelsService().UserModel.GetUserByID(c, userID, nil)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func UpdateUserProfile(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	var updatedUser interfaces.User
	if err := c.ShouldBindJSON(&updatedUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user json"})
		return
	}

	updatedUser.ID = userID

	if err := impl.GetModelsService().UserModel.UpdateUser(c, &updatedUser, nil); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to update user"})
		return
	}

	c.JSON(http.StatusOK, updatedUser)
}

func DeleteUser(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	if err := impl.GetModelsService().UserModel.DeleteUser(c, userID, nil); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user deleted successfully"})
}
