package main

import (
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"` // Note: In a real-world application, never store plain-text passwords.
}

var users = []User{{ID: "test", Username: "test", Password: "test"}}
var jwtKey = []byte("my_secret_key") // This should be a secret and stored securely.

// Claims struct will be used to customize standard claims
type Claims struct {
	UserID string `json:"user_id"`
	jwt.StandardClaims
}

func generateToken(user User) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)

	claims := &Claims{
		UserID: user.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func validateToken(tknStr string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if !token.Valid {
		return nil, err
	}

	return claims, nil
}

// Register route
// Register route
func register(c *gin.Context) {
	var newUser User
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Hash the password
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(newUser.Password), 8)
	newUser.Password = string(hashedPassword)

	// Generate UUID
	newUser.ID = uuid.New().String()

	users = append(users, newUser)
	c.JSON(http.StatusOK, gin.H{"data": newUser})
}

// Login route
func login(c *gin.Context) {
	var creds User
	if err := c.ShouldBindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var foundUser User
	for _, user := range users {
		log.Println("ID:", user.ID)
		log.Println("User:", user.Username)
		if user.Username == creds.Username {
			foundUser = user
			break
		}
	}

	// User not found
	if foundUser.ID == "0" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user"})
		return
	}
	log.Println("foundPass:", foundUser.Password)
	log.Println("credPass:", creds.Password)
	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(creds.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	// Generate JWT
	token, err := generateToken(foundUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
