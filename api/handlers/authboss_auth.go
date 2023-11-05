package handlers

import (
	"log"
	"net/http"
	"xspends/models"

	"github.com/gin-gonic/gin"
	"github.com/volatiletech/authboss/v3"
	"golang.org/x/crypto/bcrypt"
)

func JWTRegisterHandler(ab *authboss.Authboss) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Binding the incoming JSON to the newUser struct
		var newUser models.User
		if err := c.ShouldBindJSON(&newUser); err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Invalid input data"})
			return
		}
		//TODO: Here you should validate the newUser fields as per your application's requirements
		exists, err := models.UserExists(c, newUser.Username, newUser.Email)
		if err != nil {
			log.Printf("[JWTRegisterHandler] Error checking user existence: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if exists {
			log.Printf("[JWTRegisterHandler] User already exists: %v", err)
			c.JSON(http.StatusConflict, gin.H{"error": ErrUserExists.Error()})
			return
		}
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), 12)
		if err != nil {
			log.Printf("[JWTRegisterHandler] Error hashing password: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": ErrHashingPassword.Error()})
			return
		}
		newUser.Password = string(hashedPassword)
		// You must assert the type of the storer to the concrete type (*UserStorer) to access the Create method
		userStorer, ok := ab.Config.Storage.Server.(*models.UserStorer)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "User storage configuration error"})
			return
		}

		// Create the user using the UserStorer which is part of AuthBoss's user creation
		err = userStorer.Create(c.Request.Context(), &newUser)
		if err != nil {
			// Handle the error which may include user already exists etc.
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Generate the JWT token for the new user using the generateToken method
		token, err := generateToken(newUser.ID)
		if err != nil {
			// Handle the error in token generation
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating token"})
			return
		}

		// Return the JWT token to the client
		c.JSON(http.StatusOK, gin.H{"token": token})
	}
}
func JWTLoginHandler(ab *authboss.Authboss) gin.HandlerFunc {
	return func(c *gin.Context) {
		var creds models.User
		if err := c.ShouldBindJSON(&creds); err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": ErrInvalidInputData.Error()})
			return
		}

		// You need to assert the type of the storer to the concrete type (*UserStorer) to access the Load method
		userStorer, ok := ab.Config.Storage.Server.(*models.UserStorer)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "User storage configuration error"})
			return
		}

		// Use the userStorer to load the user
		userInterface, err := userStorer.Load(c.Request.Context(), creds.Username)
		if err != nil {
			if err == authboss.ErrUserNotFound {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			}
			return
		}

		// Assert the type of the user to the concrete type (*models.User)
		user, ok := userInterface.(*models.User)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "User retrieval error"})
			return
		}

		// Validate the password
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password)); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}

		// Generate the JWT token for the logged-in user
		token, err := generateToken(user.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": ErrGeneratingToken.Error()})
			return
		}

		// Return the JWT token to the client
		c.JSON(http.StatusOK, gin.H{"token": token})
	}
}
