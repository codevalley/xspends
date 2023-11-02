package handlers

import (
	"net/http"
	"xspends/models" // Importing our data models

	"github.com/gin-gonic/gin"
)

// GetUserProfile is a handler function to fetch the profile details of the authenticated user.
func GetUserProfile(c *gin.Context) {
	// Retrieve the userID from the request's context.
	// This userID is expected to have been set in the context by a previous middleware,
	// usually the one handling JWT authentication.
	userID, exists := c.Get("userID")

	// If the userID doesn't exist in the context, it suggests the user is not authenticated.
	// In such cases, an unauthorized error is returned.
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	// Convert the interface type to int.
	intUserID, ok := userID.(int64)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	// Use the userID to fetch the user's details from the database using the `GetUserByID` function from the models package.
	user, err := models.GetUserByID(c, intUserID)

	// If there's an error fetching the user, it's assumed the user does not exist in the database.
	// Therefore, a not found error is returned.
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	// If everything goes smoothly, the user's details are returned with a status of OK.
	c.JSON(http.StatusOK, user)
}

// UpdateUserProfile is a handler function to update the details of the authenticated user.
func UpdateUserProfile(c *gin.Context) {
	// Retrieve the userID from the request's context.
	userID := c.MustGet("userID").(int64)

	var updatedUser models.User
	if err := c.ShouldBindJSON(&updatedUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set the ID of the updatedUser to the userID we retrieved from the context.
	updatedUser.ID = userID

	if err := models.UpdateUser(c, &updatedUser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to update user"})
		return
	}

	c.JSON(http.StatusOK, updatedUser)
}

func DeleteUser(c *gin.Context) {
	userID := c.MustGet("userID").(int64)

	if err := models.DeleteUser(c, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user deleted successfully"})
}
