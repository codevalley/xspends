package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
	"xspends/models"
	"xspends/util"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/volatiletech/authboss/v3"
	"golang.org/x/crypto/bcrypt"
)

type ABClaims struct {
	UserID    int64  `json:"user_id"`
	SessionID string `json:"session_id"`
	jwt.StandardClaims
}

const (
	accessTokenExpiryMins  = 30
	refreshTokenExpiryMins = 1440
)

///////////////////////////////////////////////////////////////////////////////

func generateAccessToken(userID int64, sessionID string) (string, error) {
	expirationTime := time.Now().Add(accessTokenExpiryMins * time.Minute)
	claims := &ABClaims{
		UserID:    userID,
		SessionID: sessionID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JwtKey)
}

func generateRefreshToken(userID int64, sessionID string) (string, error) {
	expirationTime := time.Now().Add(refreshTokenExpiryMins * time.Minute)
	claims := &ABClaims{
		UserID:    userID,
		SessionID: sessionID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JwtKey)
}

func RefreshTokenHandler(ctx context.Context, oldRefreshToken string, ab *authboss.Authboss) (string, string, error) {
	claims := &ABClaims{}
	tkn, err := jwt.ParseWithClaims(oldRefreshToken, claims, func(token *jwt.Token) (interface{}, error) {
		return JwtKey, nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return "", "", fmt.Errorf("invalid refresh token signature")
		}
		return "", "", fmt.Errorf("could not parse refresh token")
	}

	if !tkn.Valid {
		return "", "", fmt.Errorf("invalid refresh token")
	}

	sessionStorer, ok := ab.Config.Storage.SessionState.(*models.SessionStorer)
	if !ok {
		return "", "", fmt.Errorf("session storage configuration error")
	}

	// Delete old session
	err = sessionStorer.Delete(ctx, claims.SessionID)
	if err != nil {
		return "", "", fmt.Errorf("could not delete old session")
	}

	// Create new session
	newSessionID, err := util.GenerateSnowflakeID()
	if err != nil {
		return "", "", fmt.Errorf("could not generate new session ID")
	}
	err = sessionStorer.Save(ctx, strconv.FormatInt(newSessionID, 10), strconv.FormatInt(claims.UserID, 10), 24*time.Hour)
	if err != nil {
		return "", "", fmt.Errorf("could not save new session")
	}

	// Generate new access token
	newAccessToken, err := generateAccessToken(claims.UserID, strconv.FormatInt(newSessionID, 10))
	if err != nil {
		return "", "", fmt.Errorf("could not generate new access token")
	}

	// Generate new refresh token
	newRefreshToken, err := generateRefreshToken(claims.UserID, strconv.FormatInt(newSessionID, 10))
	if err != nil {
		return "", "", fmt.Errorf("could not generate new refresh token")
	}

	return newAccessToken, newRefreshToken, nil
}
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
		newSessionID, err := util.GenerateSnowflakeID()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating session ID"})
			return
		}
		// Generate the JWT access token for the new user
		accessToken, err := generateAccessToken(newUser.ID, strconv.FormatInt(newSessionID, 10))
		if err != nil {
			// Handle the error in token generation
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating access token"})
			return
		}

		// Generate the JWT refresh token for the new user
		refreshToken, err := generateRefreshToken(newUser.ID, strconv.FormatInt(newSessionID, 10))
		if err != nil {
			// Handle the error in token generation
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating refresh token"})
			return
		}

		// Return the JWT access token and refresh token to the client
		c.JSON(http.StatusOK, gin.H{"access_token": accessToken, "refresh_token": refreshToken})
	}
}
func JWTLoginHandler(ab *authboss.Authboss) gin.HandlerFunc {
	return func(c *gin.Context) {
		var creds models.User
		if err := c.ShouldBindJSON(&creds); err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Invalid input data"})
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

		newSessionID, err := util.GenerateSnowflakeID()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Generate the JWT access token for the logged-in user
		accessToken, err := generateAccessToken(user.ID, strconv.FormatInt(newSessionID, 10))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating access token"})
			return
		}

		// Generate the JWT refresh token for the logged-in user
		refreshToken, err := generateRefreshToken(user.ID, strconv.FormatInt(newSessionID, 10))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating refresh token"})
			return
		}

		// Return the JWT access token and refresh token to the client
		c.JSON(http.StatusOK, gin.H{"access_token": accessToken, "refresh_token": refreshToken})
	}
}

func JWTRefreshHandler(ab *authboss.Authboss) gin.HandlerFunc {
	return func(c *gin.Context) {
		oldRefreshToken := c.PostForm("refresh_token")
		newAccessToken, newRefreshToken, err := RefreshTokenHandler(c.Request.Context(), oldRefreshToken, ab)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"access_token": newAccessToken, "refresh_token": newRefreshToken})
	}
}
