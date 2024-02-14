/*
MIT License

Copyright (c) 2022 Narayan Babu

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
// This file specifically contains the handlers related to user authentication using JWT.

package handlers

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
	"xspends/models/impl"
	"xspends/models/interfaces"
	"xspends/util"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/volatiletech/authboss/v3"
	"golang.org/x/crypto/bcrypt"
)

// JWTClaims struct is used for JWT claims
type JWTClaims struct {
	UserID    int64  `json:"user_id"`
	SessionID string `json:"session_id"`
	ScopeID   int64  `json:"scope_id"`
	jwt.StandardClaims
}

const (
	tokenExpiryMins        = 30
	refreshTokenExpiryMins = 1440
)

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

// GenerateTokenWithTTL generates a JWT with a specific time-to-live (TTL)
func GenerateTokenWithTTL(userID int64, scopeID int64, sessionID string, expiryMins int) (string, error) {
	expirationTime := time.Now().Add(time.Duration(expiryMins) * time.Minute)
	claims := &JWTClaims{
		UserID:    userID,
		SessionID: sessionID,
		ScopeID:   scopeID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JwtKey)
}

// RefreshTokenHandler handles the refreshing of JWT tokens
func RefreshTokenHandler(ctx context.Context, oldRefreshToken string, ab *authboss.Authboss) (string, string, error) {
	claims := &JWTClaims{}
	fmt.Println("Token:", oldRefreshToken)
	tkn, err := jwt.ParseWithClaims(oldRefreshToken, claims, func(token *jwt.Token) (interface{}, error) {
		return JwtKey, nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return "", "", errors.Wrap(err, "[RefreshTokenHandler] invalid refresh token signature")
		}
		return "", "", errors.Wrap(err, "[RefreshTokenHandler] could not parse refresh token")
	}

	if !tkn.Valid {
		return "", "", errors.New("[RefreshTokenHandler] invalid refresh token")
	}

	sessionStorer, ok := ab.Config.Storage.SessionState.(*impl.SessionStorer)
	if !ok {
		return "", "", errors.New("[RefreshTokenHandler] session storage configuration error")
	}
	// Fetch the old refresh token from the session store
	storedRefreshToken, err := sessionStorer.Load(ctx, claims.SessionID)
	if err != nil {
		return "", "", errors.Wrap(err, "[RefreshTokenHandler] could not fetch old refresh token")
	}

	// Verify the old refresh token
	if storedRefreshToken != oldRefreshToken {
		return "", "", errors.New("[RefreshTokenHandler] provided refresh token does not match stored refresh token")
	}

	// Delete old session
	err = sessionStorer.Delete(ctx, claims.SessionID)
	if err != nil {
		return "", "", errors.Wrap(err, "[RefreshTokenHandler] could not delete old session")
	}

	// Create new session
	newSessionID, err := util.GenerateSnowflakeID()
	if err != nil {
		return "", "", errors.Wrap(err, "[RefreshTokenHandler] could not generate new session ID")
	}
	err = sessionStorer.Save(ctx, strconv.FormatInt(newSessionID, 10), strconv.FormatInt(claims.UserID, 10), 24*time.Hour)
	if err != nil {
		return "", "", errors.Wrap(err, "[RefreshTokenHandler] could not save new session")
	}

	// Generate new access token
	newAccessToken, err := GenerateTokenWithTTL(claims.UserID, claims.ScopeID, strconv.FormatInt(newSessionID, 10), tokenExpiryMins)
	if err != nil {
		return "", "", errors.Wrap(err, "[RefreshTokenHandler] could not generate new access token")
	}

	// Generate new refresh token
	newRefreshToken, err := GenerateTokenWithTTL(claims.UserID, claims.ScopeID, strconv.FormatInt(newSessionID, 10), refreshTokenExpiryMins)
	if err != nil {
		return "", "", errors.Wrap(err, "[RefreshTokenHandler] could not generate new refresh token")
	}

	return newAccessToken, newRefreshToken, nil
}

// @Summary Register a new user
// @Description Register a new user with email and password
// @ID register-user
// @Accept  json
// @Produce  json
// @Param user body impl.User true "User info for registration"
// @Success 200  {object}  map[string]interface{}  "User registered successfully"
// @Failure 400  {object}  map[string]string  "Invalid input data"
// @Failure 500  {object}  map[string]string  "Internal Server Error"
// @Router /auth/register [post]

func JWTRegisterHandler(ab *authboss.Authboss) gin.HandlerFunc {
	return func(c *gin.Context) {
		var newUser interfaces.User
		if err := c.ShouldBindJSON(&newUser); err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Invalid input data"})
			return
		}

		exists, err := impl.GetModelsService().UserModel.UserExists(c, newUser.Username, newUser.Email, nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": errors.Wrap(err, "[JWTRegisterHandler] Error checking user existence").Error()})
			return
		}
		if exists {
			c.JSON(http.StatusConflict, gin.H{"error": ErrUserExists.Error()})
			return
		}

		hashedPassword, err := hashPassword(newUser.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		newUser.Password = hashedPassword

		userStorer, ok := ab.Config.Storage.Server.(*impl.UserStorer)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "[JWTRegisterHandler] User storage configuration error"})
			return
		}

		err = userStorer.Create(c.Request.Context(), &newUser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": errors.Wrap(err, "[JWTRegisterHandler] Error creating user").Error()})
			return
		}

		newSessionID, err := util.GenerateSnowflakeID()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": errors.Wrap(err, "[JWTRegisterHandler] Error generating session ID").Error()})
			return
		}

		sessionID := strconv.FormatInt(newSessionID, 10)

		accessToken, err := GenerateTokenWithTTL(newUser.ID, newUser.Scope, sessionID, tokenExpiryMins)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": errors.Wrap(err, "[JWTRegisterHandler] Error generating access token").Error()})
			return
		}

		refreshToken, err := GenerateTokenWithTTL(newUser.ID, newUser.Scope, sessionID, refreshTokenExpiryMins)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": errors.Wrap(err, "[JWTRegisterHandler] Error generating refresh token").Error()})
			return
		}

		sessionStorer, ok := ab.Config.Storage.SessionState.(*impl.SessionStorer)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "[JWTRegisterHandler] Session storage configuration error"})
			return
		}

		err = sessionStorer.Save(c.Request.Context(), sessionID, refreshToken, refreshTokenExpiryMins*time.Minute)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": errors.Wrap(err, "[JWTRegisterHandler] Error storing refresh token").Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"access_token": accessToken, "refresh_token": refreshToken})
	}
}

// @Summary Login an existing user
// @Description Login an existing user with email and password
// @ID login-user
// @Accept  json
// @Produce  json
// @Param   email  body  string  true  "User Email"
// @Param   password  body  string  true  "User Password"
// @Success 200  {object}  map[string]string  "Access and refresh tokens"
// @Failure 400  {object}  map[string]string  "Invalid input data"
// @Failure 401  {object}  map[string]string  "Invalid credentials or username"
// @Failure 500  {object}  map[string]string  "Internal Server Error"
// @Router /auth/login [post]

func JWTLoginHandler(ab *authboss.Authboss) gin.HandlerFunc {
	return func(c *gin.Context) {
		var creds interfaces.User
		if err := c.ShouldBindJSON(&creds); err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Invalid input data"})
			return
		}

		userStorer, ok := ab.Config.Storage.Server.(*impl.UserStorer)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "[JWTLoginHandler] User storage configuration error"})
			return
		}

		userInterface, err := userStorer.Load(c.Request.Context(), creds.Username)
		if err != nil {
			if err == authboss.ErrUserNotFound {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": errors.Wrap(err, "[JWTLoginHandler] Error loading user").Error()})
			}
			return
		}

		user, ok := userInterface.(*interfaces.User)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "[JWTLoginHandler] User retrieval error"})
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password)); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}

		newSessionID, err := util.GenerateSnowflakeID()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": errors.Wrap(err, "[JWTLoginHandler] Error generating session ID").Error()})
			return
		}

		sessionID := strconv.FormatInt(newSessionID, 10)

		accessToken, err := GenerateTokenWithTTL(user.ID, user.Scope, sessionID, tokenExpiryMins)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": errors.Wrap(err, "[JWTLoginHandler] Error generating access token").Error()})
			return
		}

		refreshToken, err := GenerateTokenWithTTL(user.ID, user.Scope, sessionID, refreshTokenExpiryMins)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": errors.Wrap(err, "[JWTLoginHandler] Error generating refresh token").Error()})
			return
		}

		sessionStorer, ok := ab.Config.Storage.SessionState.(*impl.SessionStorer)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "[JWTLoginHandler] Session storage configuration error"})
			return
		}

		err = sessionStorer.Save(c.Request.Context(), sessionID, refreshToken, refreshTokenExpiryMins*time.Minute)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": errors.Wrap(err, "[JWTLoginHandler] Error storing refresh token").Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"access_token": accessToken, "refresh_token": refreshToken})
	}
}

// @Summary Refresh JWT token
// @Description Refresh JWT token for an authenticated user
// @ID refresh-token
// @Accept  json
// @Produce  json
// @Param   refresh_token  body  string  true  "Refresh token"
// @Success 200  {object}  map[string]string  "New access and refresh tokens"
// @Failure 400  {object}  map[string]string  "Invalid refresh token"
// @Failure 500  {object}  map[string]string  "Internal Server Error"
// @Router /auth/refresh [post]

func JWTRefreshHandler(ab *authboss.Authboss) gin.HandlerFunc {
	return func(c *gin.Context) {
		var body map[string]string
		if err := c.BindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}
		oldRefreshToken := body["refresh_token"]
		newAccessToken, newRefreshToken, err := RefreshTokenHandler(c.Request.Context(), oldRefreshToken, ab)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": errors.Wrap(err, "[JWTRefreshHandler] Error refreshing token").Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"access_token": newAccessToken, "refresh_token": newRefreshToken})
	}
}

// @Summary Logout the current user
// @Description Logout the user by invalidating their session
// @ID logout-user
// @Accept  json
// @Produce  json
// @Param   refresh_token  body  string  true  "Refresh Token"
// @Success 200  {object}  map[string]string  "message: Logged out successfully"
// @Failure 400  {object}  map[string]string  "Invalid request body"
// @Failure 500  {object}  map[string]string  "Internal Server Error"
// @Router /auth/logout [post]

func JWTLogoutHandler(ab *authboss.Authboss) gin.HandlerFunc {
	return func(c *gin.Context) {
		var body map[string]string
		if err := c.BindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		refreshToken := body["refresh_token"]
		claims := &JWTClaims{}
		_, err := jwt.ParseWithClaims(refreshToken, claims, func(token *jwt.Token) (interface{}, error) {
			return JwtKey, nil
		})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": errors.Wrap(err, "[JWTLogoutHandler] Error parsing refresh token").Error()})
			return
		}

		sessionStorer, ok := ab.Config.Storage.SessionState.(*impl.SessionStorer)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "[JWTLogoutHandler] Session storage configuration error"})
			return
		}

		err = sessionStorer.Delete(c.Request.Context(), claims.SessionID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": errors.Wrap(err, "[JWTLogoutHandler] Error deleting session").Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
	}
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	return string(bytes), err
}
