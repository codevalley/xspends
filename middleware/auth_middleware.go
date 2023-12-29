/*
MIT License

# Copyright (c) 2023 Narayan Babu

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

package middleware

import (
	"net/http"
	"strings"

	"xspends/api/handlers"
	"xspends/kvstore"
	"xspends/models/impl"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/volatiletech/authboss/v3"
)

// Initialize AuthBoss.
var ab *authboss.Authboss

const userIDKey = "userID"
const authKey = "Authorization"

func AuthMiddleware(ab *authboss.Authboss) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from the Authorization header
		authorizationHeader := c.GetHeader(authKey)
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

		c.Set(userIDKey, claims.UserID)

		// Continue with the request
		c.Next()
	}
}

func EnsureUserID() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, exists := c.Get(userIDKey)
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			c.Abort() // This prevents the handler from being executed if the check fails
			return
		}
		c.Next()
	}
}

func SetupAuthBoss(router *gin.Engine, kvClient kvstore.RawKVClientInterface) *authboss.Authboss {
	// ... other setup
	ab = authboss.New()
	// Set up AuthBoss storage with your custom implementations
	ab.Config.Storage.Server = impl.NewUserStorer()
	ab.Config.Storage.SessionState = impl.NewSessionStorer(kvClient)
	ab.Config.Storage.CookieState = impl.NewCookieStorer(kvClient)

	// ... finish setup
	return ab
}
