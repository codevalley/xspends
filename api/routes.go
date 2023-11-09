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

// Package api contains all the routes for the application.
// It sets up the routes for authentication, sources, categories, tags, and transactions.
package api

import (
	"database/sql"
	"xspends/api/handlers"
	"xspends/kvstore"
	"xspends/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all the routes for the application.
// It takes a gin Engine, a sql DB, and a kvstore RawKVClientInterface as parameters.
func SetupRoutes(r *gin.Engine, db *sql.DB, kvClient kvstore.RawKVClientInterface) {
	ab := middleware.SetupAuthBoss(r, db, kvClient)

	// Health check endpoint
	// This endpoint is used to check the health status of the application.
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "UP"})
	})

	// Authentication routes using custom JWT handlers
	// These routes are used for user registration, login, and token refresh.
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
	apiRoutes.Use(middleware.AuthMiddleware(ab), middleware.EnsureUserID())

	// Source routes
	// These routes are used for managing sources.
	sources := apiRoutes.Group("/sources")
	{
		sources.GET("", handlers.ListSources)         // List all sources
		sources.POST("", handlers.CreateSource)       // Create a new source
		sources.GET("/:id", handlers.GetSource)       // Get a specific source
		sources.PUT("/:id", handlers.UpdateSource)    // Update a specific source
		sources.DELETE("/:id", handlers.DeleteSource) // Delete a specific source
	}

	// Category routes
	// These routes are used for managing categories.
	categories := apiRoutes.Group("/categories")
	{
		categories.GET("", handlers.ListCategories)        // List all categories
		categories.POST("", handlers.CreateCategory)       // Create a new category
		categories.GET("/:id", handlers.GetCategory)       // Get a specific category
		categories.PUT("/:id", handlers.UpdateCategory)    // Update a specific category
		categories.DELETE("/:id", handlers.DeleteCategory) // Delete a specific category
	}

	// Tag routes
	// These routes are used for managing tags.
	tags := apiRoutes.Group("/tags")
	{
		tags.GET("", handlers.ListTags)         // List all tags
		tags.POST("", handlers.CreateTag)       // Create a new tag
		tags.GET("/:id", handlers.GetTag)       // Get a specific tag
		tags.PUT("/:id", handlers.UpdateTag)    // Update a specific tag
		tags.DELETE("/:id", handlers.DeleteTag) // Delete a specific tag
	}

	// Transaction routes
	// These routes are used for managing transactions and transaction tags.
	transactions := apiRoutes.Group("/transactions")
	{
		transactions.GET("", handlers.ListTransactions)         // List all transactions
		transactions.POST("", handlers.CreateTransaction)       // Create a new transaction
		transactions.GET("/:id", handlers.GetTransaction)       // Get a specific transaction
		transactions.PUT("/:id", handlers.UpdateTransaction)    // Update a specific transaction
		transactions.DELETE("/:id", handlers.DeleteTransaction) // Delete a specific transaction

		// Transaction Tags routes
		transactions.GET("/:id/tags", handlers.ListTransactionTags)                // List all tags for a specific transaction
		transactions.POST("/:id/tags", handlers.AddTagToTransaction)               // Add a tag to a specific transaction
		transactions.DELETE("/:id/tags/:tagID", handlers.RemoveTagFromTransaction) // Remove a tag from a specific transaction
	}
}
