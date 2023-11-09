/*
MIT License

Copyright (c) 2023 Narayan Babu

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

// Package handlers contains the HTTP handler functions for the application.
// This file specifically contains the handlers related to user authentication.
// Deprecated: This file is deprecated and will be removed in a future version.

package handlers

import (
	"errors"
	"log"
	"net/http"
	"os"
	"time"
	"xspends/models"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// Claims struct is used for JWT claims
type Claims struct {
	UserID int64 `json:"user_id"`
	jwt.StandardClaims
}

// Error variables for common error messages
var (
	ErrInvalidInputData = errors.New("invalid input data")
	ErrUserExists       = errors.New("username or email already exists")
	ErrHashingPassword  = errors.New("error hashing password")
	ErrInsertingUser    = errors.New("error inserting user into database")
	ErrGeneratingToken  = errors.New("error generating token")
)

// JwtKey is the key used for signing JWTs
var JwtKey = getJwtKey()

// getJwtKey retrieves the JWT key from environment variables or uses a default key
func getJwtKey() []byte {
	key := os.Getenv("JWT_KEY")
	if key == "" {
		// Fallback to a default key
		key = "uNauz8OMH3UzF6wum99OD6dsm1wSdMquDGkWznT6JrQ="
	}
	return []byte(key)
}

// generateToken generates a JWT for a given user ID
func generateToken(userID int64) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JwtKey)
}

// Register is the handler for user registration
func Register(c *gin.Context) {
	var newUser models.User
	if err := c.ShouldBindJSON(&newUser); err != nil {
		log.Printf("[Register] Error binding JSON: %v", err)
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": ErrInvalidInputData.Error()})
		return
	}

	exists, err := models.UserExists(c, newUser.Username, newUser.Email)
	if err != nil {
		log.Printf("[Register] Error checking user existence: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if exists {
		log.Printf("[Register] User already exists: %v", err)
		c.JSON(http.StatusConflict, gin.H{"error": ErrUserExists.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), 12)
	if err != nil {
		log.Printf("[Register] Error hashing password: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": ErrHashingPassword.Error()})
		return
	}
	newUser.Password = string(hashedPassword)

	err = models.InsertUser(c, &newUser)
	if err != nil {
		log.Printf("[Register] Error inserting user into database: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": ErrInsertingUser.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}

// Login is the handler for user login
func Login(c *gin.Context) {
	var creds models.User
	if err := c.ShouldBindJSON(&creds); err != nil {
		log.Printf("[Login] Error binding JSON: %v", err)
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": ErrInvalidInputData.Error()})
		return
	}

	user, err := models.GetUserByUsername(c, creds.Username)
	if err != nil {
		log.Printf("[Login] Error retrieving user: %v", err)
		if err == models.ErrUserNotFound {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving user"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password)); err != nil {
		log.Printf("[Login] Invalid credentials: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token, err := generateToken(user.ID)
	if err != nil {
		log.Printf("[Login] Error generating token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": ErrGeneratingToken.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
