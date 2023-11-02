package handlers

import (
	"net/http"
	"os"
	"time"
	"xspends/models"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type Claims struct {
	UserID int64 `json:"user_id"`
	jwt.StandardClaims
}

var JwtKey = []byte(os.Getenv("JWT_KEY"))

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

func Register(c *gin.Context) {
	var newUser models.User
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data"})
		return
	}

	// Check if username or email already exists
	exists, err := models.UserExists(c, newUser.Username, newUser.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if exists {
		c.JSON(http.StatusConflict, gin.H{"error": "Username or email already exists"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), 12)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
		return
	}
	newUser.Password = string(hashedPassword)

	// ID generation is now handled in the model, so we can safely remove the UUID generation here

	err = models.InsertUser(c, &newUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error inserting user into database"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}

func Login(c *gin.Context) {
	var creds models.User
	if err := c.ShouldBindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data"})
		return
	}

	user, err := models.GetUserByUsername(c, creds.Username)
	if err != nil {
		if err == models.ErrUserNotFound {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving user"})
		return
	}

	// if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password)); err != nil {
	// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
	// 	return
	// }

	token, err := generateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
