package handlers

import (
	"log"
	"net/http"
	"xspends/models" // Importing our data models

	"github.com/gin-gonic/gin"
)

func getUserID(c *gin.Context) (int64, bool) {
	userID, exists := c.Get("userID")
	if !exists {
		log.Printf("[getUserID] Error: user not authenticated")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return 0, false
	}

	intUserID, ok := userID.(int64)
	if !ok {
		log.Printf("[getUserID] Error: failed to convert userID to int64")
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

	user, err := models.GetUserByID(c, userID)
	if err != nil {
		log.Printf("[GetUserProfile] Error: %v", err)
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

	var updatedUser models.User
	if err := c.ShouldBindJSON(&updatedUser); err != nil {
		log.Printf("[UpdateUserProfile] Error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedUser.ID = userID

	if err := models.UpdateUser(c, &updatedUser); err != nil {
		log.Printf("[UpdateUserProfile] Error: %v", err)
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

	if err := models.DeleteUser(c, userID); err != nil {
		log.Printf("[DeleteUser] Error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user deleted successfully"})
}
