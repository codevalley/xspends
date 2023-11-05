package middleware

import (
	"database/sql"
	"net/http"
	"strings"
	"xspends/api/handlers"
	"xspends/kvstore"
	"xspends/models"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/volatiletech/authboss/v3"
)

// Initialize AuthBoss.
var ab *authboss.Authboss

func AuthMiddleware(ab *authboss.Authboss) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from the Authorization header
		authorizationHeader := c.GetHeader("Authorization")
		if authorizationHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
			c.Abort()
			return
		}

		// Extract the actual token from the Bearer token format
		bearerToken := strings.Split(authorizationHeader, " ")
		if len(bearerToken) != 2 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			c.Abort()
			return
		}
		tokenStr := bearerToken[1]

		// Parse and validate the token
		claims := &handlers.Claims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return handlers.JwtKey, nil
		})
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// If the token is valid, store the user data (from the JWT claims) in the context
		c.Set("userID", claims.UserID)

		// Continue with the request
		c.Next()
	}
}

func EnsureUserID() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			c.Abort() // This prevents the handler from being executed if the check fails
			return
		}
		c.Next()
	}
}

func SetupAuthBoss(router *gin.Engine, db *sql.DB, kvClient kvstore.RawKVClientInterface) *authboss.Authboss {
	// ... other setup
	ab = authboss.New()
	// Set up AuthBoss storage with your custom implementations
	ab.Config.Storage.Server = models.NewUserStorer(db)
	ab.Config.Storage.SessionState = models.NewSessionStorer(kvClient)
	ab.Config.Storage.CookieState = models.NewCookieStorer(kvClient)

	// ... finish setup
	return ab
}
