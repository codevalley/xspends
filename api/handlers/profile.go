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

	// Use the userID to fetch the user's details from the database using the `GetUserByID` function from the models package.
	// Note: userID is an interface type when retrieved from the context, so it's type-asserted to string.
	user, err := models.GetUserByID(userID.(string))

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
	// The userID is expected to be injected into the context by a previous middleware (usually the JWT authentication middleware).
	userID, exists := c.Get("userID")

	// If the userID doesn't exist in the context, it means the user is not authenticated.
	// So, return an unauthorized error.
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	// Declare a variable of type model.User to hold the user's updated details.
	var updatedUser models.User

	// Try to bind (or map) the incoming JSON request body to the updatedUser variable.
	// If there's an error in parsing the JSON or if it doesn't match the expected structure, return a bad request error.
	if err := c.ShouldBindJSON(&updatedUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set the ID of the updatedUser to the userID we retrieved from the context.
	// This ensures that the user can only update their own profile and not someone else's.
	updatedUser.ID = userID.(string)

	// Call the model's UpdateUser function to update the user's details in the database.
	// If there's an error during the update, return an internal server error.
	if err := models.UpdateUser(&updatedUser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to update user"})
		return
	}

	// If everything goes well, return the updated user details with a status OK.
	c.JSON(http.StatusOK, updatedUser)
}

// DeleteUser deletes the authenticated user.
func DeleteUser(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	if err := models.DeleteUser(userID.(string)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user deleted successfully"})
}
