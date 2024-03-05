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
/*
SetupRoutes configures all the routes for the application.
It takes a gin Engine, a sql DB, and a kvstore RawKVClientInterface as parameters.

Inputs:
- r: A pointer to a gin.Engine instance representing the Gin router.
- db: A pointer to a sql.DB instance representing the database connection.
- kvClient: An interface representing the key-value store client.

Flow:
1. The function sets up the routes for the application.
2. It initializes the authentication boss (ab) using the SetupAuthBoss function.
3. It defines a health check endpoint.
4. It sets up authentication routes for user registration, login, token refresh, and logout using custom JWT handlers.
5. It sets up routes for managing sources, categories, tags, and transactions.
6. Each route is associated with a specific HTTP method and URL path, and is handled by a corresponding handler function.

Outputs:
- The routes are configured and ready to handle incoming requests.
*/

// Package api contains all the routes for the application.
// It sets up the routes for authentication, sources, categories, tags, and transactions.
package api

import (
	"encoding/json"
	"net/http"
	"os"
	"xspends/api/handlers"
	"xspends/kvstore"
	"xspends/middleware"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// SetupRoutes sets up all the routes for the application
// @description This function will set all routes
func SetupRoutes(r *gin.Engine, kvClient kvstore.RawKVClientInterface) {
	ab := middleware.SetupAuthBoss(r, kvClient)

	setupHealthEndpoint(r)
	setupSwaggerHandler(r)

	auth := r.Group("/auth")
	{
		auth.POST("/register", handlers.JWTRegisterHandler(ab)) // Register a new user
		auth.POST("/login", handlers.JWTLoginHandler(ab))       // Login an existing user
		auth.POST("/refresh", handlers.JWTRefreshHandler(ab))   // Refresh JWT token
		auth.POST("/logout", handlers.JWTLogoutHandler(ab))     // Logout a user
	}

	// All other routes should be protected by the AuthMiddleware
	// These routes are used for managing sources, categories, tags, and transactions.
	apiRoutes := r.Group("/")
	apiRoutes.Use(middleware.AuthMiddleware(ab), middleware.ScopeMiddleware(), middleware.EnsureUserID(), middleware.EnsureScopeID())
	// Source routes
	// These routes are used for managing sources.
	sources := apiRoutes.Group("/sources")
	{
		sources.GET("", handlers.ListSources)
		sources.POST("", handlers.CreateSource)
		sources.GET("/:id", handlers.GetSource)
		sources.PUT("/:id", handlers.UpdateSource)
		sources.DELETE("/:id", handlers.DeleteSource)
	}
	// Category routes
	// These routes are used for managing categories.
	categories := apiRoutes.Group("/categories")
	{
		categories.GET("", handlers.ListCategories)
		categories.POST("", handlers.CreateCategory)
		categories.GET("/:id", handlers.GetCategory)
		categories.PUT("/:id", handlers.UpdateCategory)
		categories.DELETE("/:id", handlers.DeleteCategory)
	}
	// Tag routes
	// These routes are used for managing tags.
	tags := apiRoutes.Group("/tags")
	{
		tags.GET("", handlers.ListTags)
		tags.POST("", handlers.CreateTag)
		tags.GET("/:id", handlers.GetTag)
		tags.PUT("/:id", handlers.UpdateTag)
		tags.DELETE("/:id", handlers.DeleteTag)
	}
	// Transaction routes
	// These routes are used for managing transactions and transaction tags.
	transactions := apiRoutes.Group("/transactions")
	{
		transactions.GET("", handlers.ListTransactions)
		transactions.POST("", handlers.CreateTransaction)
		transactions.GET("/:id", handlers.GetTransaction)
		transactions.PUT("/:id", handlers.UpdateTransaction)
		transactions.DELETE("/:id", handlers.DeleteTransaction)
		transactions.GET("/:id/tags", handlers.ListTransactionTags)
		transactions.POST("/:id/tags", handlers.AddTagToTransaction)
		transactions.DELETE("/:id/tags/:tagID", handlers.RemoveTagFromTransaction)
	}
}

// @Summary Health check
// @Description Check the health status of the application
// @ID get-health
// @Produce  json
// @Success 200 {object} map[string]string "Health status of the application"
// @Router /health [get]
func setupHealthEndpoint(r *gin.Engine) {
	// Health check endpoint
	// This endpoint is used to check the health status of the application.
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "UP"})
	})
}

func setupSwaggerHandler(r *gin.Engine) {
	r.GET("/swagger/*any", func(c *gin.Context) {
		path := c.Param("any")
		if path == "/doc.json" {
			swaggerFilePath := os.Getenv("SWAGGER_JSON_PATH")
			swaggerJSON, err := setSwaggerHost(swaggerFilePath)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load swagger file: " + err.Error()})
				return
			}
			c.Data(http.StatusOK, "application/json", swaggerJSON)
			return
		}

		// Serve Swagger UI for any other path under "/swagger/"
		ginSwagger.CustomWrapHandler(&ginSwagger.Config{
			URL: "http://" + c.Request.Host + "/swagger/doc.json",
		}, swaggerFiles.Handler)(c)
	})
}

func setSwaggerHost(filePath string) ([]byte, error) {
	// Read the original swagger.json file
	jsonFile, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var swagger map[string]interface{}
	err = json.Unmarshal(jsonFile, &swagger)
	if err != nil {
		return nil, err
	}

	// Use environment variables to determine the host and port
	host := os.Getenv("SWAGGER_HOST") // Example: "api.example.com"
	port := os.Getenv("SWAGGER_PORT") // Example: "443" for HTTPS
	if host == "" {
		host = "127.0.0.1" // Default host
	}
	if port == "" {
		port = "8080" // Default port
	}
	host = host + ":" + port
	swagger["host"] = host

	// Marshal the modified swagger back to JSON
	modifiedSwaggerJSON, err := json.Marshal(swagger)
	if err != nil {
		return nil, err
	}

	return modifiedSwaggerJSON, nil
}
