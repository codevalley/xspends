package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"
	"xspends/models"
	"xspends/util"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/volatiletech/authboss/v3"
	"golang.org/x/crypto/bcrypt"
)

type ABClaims struct {
	UserID    int64  `json:"user_id"`
	SessionID string `json:"session_id"`
	jwt.StandardClaims
}

const (
	tokenExpiryMins        = 30
	refreshTokenExpiryMins = 1440
)

///////////////////////////////////////////////////////////////////////////////

func generateTokenWithTTL(userID int64, sessionID string, expiryMins int) (string, error) {
	expirationTime := time.Now().Add(time.Duration(expiryMins) * time.Minute)
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

	sessionStorer, ok := ab.Config.Storage.SessionState.(*models.SessionStorer)
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
	newAccessToken, err := generateTokenWithTTL(claims.UserID, strconv.FormatInt(newSessionID, 10), tokenExpiryMins)
	if err != nil {
		return "", "", errors.Wrap(err, "[RefreshTokenHandler] could not generate new access token")
	}

	// Generate new refresh token
	newRefreshToken, err := generateTokenWithTTL(claims.UserID, strconv.FormatInt(newSessionID, 10), refreshTokenExpiryMins)
	if err != nil {
		return "", "", errors.Wrap(err, "[RefreshTokenHandler] could not generate new refresh token")
	}

	return newAccessToken, newRefreshToken, nil
}
func JWTRegisterHandler(ab *authboss.Authboss) gin.HandlerFunc {
	return func(c *gin.Context) {
		var newUser models.User
		if err := c.ShouldBindJSON(&newUser); err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Invalid input data"})
			return
		}

		exists, err := models.UserExists(c, newUser.Username, newUser.Email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": errors.Wrap(err, "[JWTRegisterHandler] Error checking user existence").Error()})
			return
		}
		if exists {
			c.JSON(http.StatusConflict, gin.H{"error": ErrUserExists.Error()})
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), 12)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": errors.Wrap(err, "[JWTRegisterHandler] Error hashing password").Error()})
			return
		}
		newUser.Password = string(hashedPassword)

		userStorer, ok := ab.Config.Storage.Server.(*models.UserStorer)
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

		accessToken, err := generateTokenWithTTL(newUser.ID, strconv.FormatInt(newSessionID, 10), tokenExpiryMins)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": errors.Wrap(err, "[JWTRegisterHandler] Error generating access token").Error()})
			return
		}

		refreshToken, err := generateTokenWithTTL(newUser.ID, strconv.FormatInt(newSessionID, 10), refreshTokenExpiryMins)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": errors.Wrap(err, "[JWTRegisterHandler] Error generating refresh token").Error()})
			return
		}

		sessionStorer, ok := ab.Config.Storage.SessionState.(*models.SessionStorer)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "[JWTRegisterHandler] Session storage configuration error"})
			return
		}

		err = sessionStorer.Save(c.Request.Context(), strconv.FormatInt(newSessionID, 10), refreshToken, refreshTokenExpiryMins*time.Minute)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": errors.Wrap(err, "[JWTRegisterHandler] Error storing refresh token").Error()})
			return
		}

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

		userStorer, ok := ab.Config.Storage.Server.(*models.UserStorer)
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

		user, ok := userInterface.(*models.User)
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

		accessToken, err := generateTokenWithTTL(user.ID, strconv.FormatInt(newSessionID, 10), tokenExpiryMins)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": errors.Wrap(err, "[JWTLoginHandler] Error generating access token").Error()})
			return
		}

		refreshToken, err := generateTokenWithTTL(user.ID, strconv.FormatInt(newSessionID, 10), refreshTokenExpiryMins)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": errors.Wrap(err, "[JWTLoginHandler] Error generating refresh token").Error()})
			return
		}

		sessionStorer, ok := ab.Config.Storage.SessionState.(*models.SessionStorer)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "[JWTLoginHandler] Session storage configuration error"})
			return
		}

		err = sessionStorer.Save(c.Request.Context(), strconv.FormatInt(newSessionID, 10), refreshToken, refreshTokenExpiryMins*time.Minute)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": errors.Wrap(err, "[JWTLoginHandler] Error storing refresh token").Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"access_token": accessToken, "refresh_token": refreshToken})
	}
}

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
