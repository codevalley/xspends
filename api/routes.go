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
	"database/sql"
	"xspends/api/handlers"
	"xspends/kvstore"
	"xspends/middleware"

	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"

	swaggerFiles "github.com/swaggo/files"
)

// SetupRoutes configures all the routes for the application.
// It takes a gin Engine, a sql DB, and a kvstore RawKVClientInterface as parameters.
// @title XSpends API
// @version 1.0
// @description This is the API for the XSpends application.
// @host localhost:8080
// @BasePath /
func SetupRoutes(r *gin.Engine, db *sql.DB, kvClient kvstore.RawKVClientInterface) {
	ab := middleware.SetupAuthBoss(r, db, kvClient)

	// Health check endpoint
	// This endpoint is used to check the health status of the application.
	// @Summary Health check
	// @Description Check the health status of the application
	// @ID get-health
	// @Produce  json
	// @Success 200 {object} HealthResponse
	// @Router /health [get]
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "UP"})
	})

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Authentication routes using custom JWT handlers
	// These routes are used for user registration, login, and token refresh.
	// @tags Authentication
	auth := r.Group("/auth")
	{
		// @Summary Register a new user
		// @Description Register a new user with email and password
		// @ID register-user
		// @Accept  json
		// @Produce  json
		// @Param user body User true "User info for registration"
		// @Success 200 {object} User
		// @Router /auth/register [post]
		auth.POST("/register", handlers.JWTRegisterHandler(ab)) // Register a new user

		// @Summary Login an existing user
		// @Description Login an existing user with email and password
		// @ID login-user
		// @Accept  json
		// @Produce  json
		// @Param user body User true "User info for login"
		// @Success 200 {object} User
		// @Router /auth/login [post]
		auth.POST("/login", handlers.JWTLoginHandler(ab)) // Login an existing user

		// @Summary Refresh JWT token
		// @Description Refresh JWT token for an authenticated user
		// @ID refresh-token
		// @Accept  json
		// @Produce  json
		// @Param refresh_token body string true "Refresh token"
		// @Success 200 {object} TokenResponse
		// @Router /auth/refresh [post]
		auth.POST("/refresh", handlers.JWTRefreshHandler(ab)) // Refresh JWT token
		// @Summary Logout a user
		// @Description Logout a user by invalidating their session
		// @ID logout-user
		// @Accept  json
		// @Produce  json
		// @Success 200 {object} LogoutResponse
		// @Router /auth/logout [post]
		auth.POST("/logout", handlers.JWTLogoutHandler(ab)) // Logout a user
	}

	// All other routes should be protected by the AuthMiddleware
	// These routes are used for managing sources, categories, tags, and transactions.
	apiRoutes := r.Group("/")
	apiRoutes.Use(middleware.AuthMiddleware(ab), middleware.EnsureUserID())

	// Source routes
	// These routes are used for managing sources.
	sources := apiRoutes.Group("/sources")
	{
		// @Summary List all sources
		// @Description Get a list of all sources
		// @ID list-sources
		// @Accept  json
		// @Produce  json
		// @Success 200 {array} Source
		// @Router /sources [get]
		sources.GET("", handlers.ListSources)

		// @Summary Create a new source
		// @Description Create a new source with the provided information
		// @ID create-source
		// @Accept  json
		// @Produce  json
		// @Param source body Source true "Source info for creation"
		// @Success 200 {object} Source
		// @Router /sources [post]
		sources.POST("", handlers.CreateSource)

		// @Summary Get a specific source
		// @Description Get a specific source by its ID
		// @ID get-source
		// @Accept  json
		// @Produce  json
		// @Param id path int true "Source ID"
		// @Success 200 {object} Source
		// @Router /sources/{id} [get]
		sources.GET("/:id", handlers.GetSource)

		// @Summary Update a specific source
		// @Description Update a specific source by its ID
		// @ID update-source
		// @Accept  json
		// @Produce  json
		// @Param id path int true "Source ID"
		// @Param source body Source true "Source info for update"
		// @Success 200 {object} Source
		// @Router /sources/{id} [put]
		sources.PUT("/:id", handlers.UpdateSource)

		// @Summary Delete a specific source
		// @Description Delete a specific source by its ID
		// @ID delete-source
		// @Accept  json
		// @Produce  json
		// @Param id path int true "Source ID"
		// @Success 200 {object} DeleteResponse
		// @Router /sources/{id} [delete]
		sources.DELETE("/:id", handlers.DeleteSource)
	}

	// Category routes
	// These routes are used for managing categories.
	categories := apiRoutes.Group("/categories")
	{
		// @Summary List all categories
		// @Description Get a list of all categories
		// @ID list-categories
		// @Accept  json
		// @Produce  json
		// @Success 200 {array} Category
		// @Router /categories [get]
		categories.GET("", handlers.ListCategories)
		// @Summary Create a new category
		// @Description Create a new category with the provided information
		// @ID create-category
		// @Accept  json
		// @Produce  json
		// @Param category body Category true "Category info for creation"
		// @Success 200 {object} Category
		// @Router /categories [post]
		categories.POST("", handlers.CreateCategory)

		// @Summary Get a specific category
		// @Description Get a specific category by its ID
		// @ID get-category
		// @Accept  json
		// @Produce  json
		// @Param id path int true "Category ID"
		// @Success 200 {object} Category
		// @Router /categories/{id} [get]
		categories.GET("/:id", handlers.GetCategory)

		// @Summary Update a specific category
		// @Description Update a specific category by its ID
		// @ID update-category
		// @Accept  json
		// @Produce  json
		// @Param id path int true "Category ID"
		// @Param category body Category true "Category info for update"
		// @Success 200 {object} Category
		// @Router /categories/{id} [put]
		categories.PUT("/:id", handlers.UpdateCategory)

		// @Summary Delete a specific category
		// @Description Delete a specific category by its ID
		// @ID delete-category
		// @Accept  json
		// @Produce  json
		// @Param id path int true "Category ID"
		// @Success 200 {object} DeleteResponse
		// @Router /categories/{id} [delete]
		categories.DELETE("/:id", handlers.DeleteCategory)
	}

	// Tag routes
	// These routes are used for managing tags.
	tags := apiRoutes.Group("/tags")
	{
		// @Summary List all tags
		// @Description Get a list of all tags
		// @ID list-tags
		// @Accept  json
		// @Produce  json
		// @Success 200 {array} Tag
		// @Router /tags [get]
		tags.GET("", handlers.ListTags)

		// @Summary Create a new tag
		// @Description Create a new tag with the provided information
		// @ID create-tag
		// @Accept  json
		// @Produce  json
		// @Param tag body Tag true "Tag info for creation"
		// @Success 200 {object} Tag
		// @Router /tags [post]
		tags.POST("", handlers.CreateTag)

		// @Summary Get a specific tag
		// @Description Get a specific tag by its ID
		// @ID get-tag
		// @Accept  json
		// @Produce  json
		// @Param id path int true "Tag ID"
		// @Success 200 {object} Tag
		// @Router /tags/{id} [get]
		tags.GET("/:id", handlers.GetTag)
		// @Summary Update a specific tag
		// @Description Update a specific tag by its ID
		// @ID update-tag
		// @Accept  json
		// @Produce  json
		// @Param id path int true "Tag ID"
		// @Param tag body Tag true "Tag info for update"
		// @Success 200 {object} Tag
		// @Router /tags/{id} [put]
		tags.PUT("/:id", handlers.UpdateTag)

		// @Summary Delete a specific tag
		// @Description Delete a specific tag by its ID
		// @ID delete-tag
		// @Accept  json
		// @Produce  json
		// @Param id path int true "Tag ID"
		// @Success 200 {object} DeleteResponse
		// @Router /tags/{id} [delete]
		tags.DELETE("/:id", handlers.DeleteTag)
	}

	// Transaction routes
	// These routes are used for managing transactions and transaction tags.
	transactions := apiRoutes.Group("/transactions")
	{
		// @Summary List all transactions
		// @Description Get a list of all transactions
		// @ID list-transactions
		// @Accept  json
		// @Produce  json
		// @Success 200 {array} Transaction
		// @Router /transactions [get]
		transactions.GET("", handlers.ListTransactions)

		// @Summary Create a new transaction
		// @Description Create a new transaction with the provided information
		// @ID create-transaction
		// @Accept  json
		// @Produce  json
		// @Param transaction body Transaction true "Transaction info for creation"
		// @Success 200 {object} Transaction
		// @Router /transactions [post]
		transactions.POST("", handlers.CreateTransaction)

		// @Summary Get a specific transaction
		// @Description Get a specific transaction by its ID
		// @ID get-transaction
		// @Accept  json
		// @Produce  json
		// @Param id path int true "Transaction ID"
		// @Success 200 {object} Transaction
		// @Router /transactions/{id} [get]
		transactions.GET("/:id", handlers.GetTransaction)

		// @Summary Update a specific transaction
		// @Description Update a specific transaction by its ID
		// @ID update-transaction
		// @Accept  json
		// @Produce  json
		// @Param id path int true "Transaction ID"
		// @Param transaction body Transaction true "Transaction info for update"
		// @Success 200 {object} Transaction
		// @Router /transactions/{id} [put]
		transactions.PUT("/:id", handlers.UpdateTransaction)

		// @Summary Delete a specific transaction
		// @Description Delete a specific transaction by its ID
		// @ID delete-transaction
		// @Accept  json
		// @Produce  json
		// @Param id path int true "Transaction ID"
		// @Success 200 {object} DeleteResponse
		// @Router /transactions/{id} [delete]
		transactions.DELETE("/:id", handlers.DeleteTransaction)

		// Transaction Tags routes
		// @Summary List all tags for a specific transaction
		// @Description Get a list of all tags for a specific transaction
		// @ID list-transaction-tags
		// @Accept  json
		// @Produce  json
		// @Param id path int true "Transaction ID"
		// @Success 200 {array} Tag
		// @Router /transactions/{id}/tags [get]
		transactions.GET("/:id/tags", handlers.ListTransactionTags)

		// @Summary Add a tag to a specific transaction
		// @Description Add a tag to a specific transaction
		// @ID add-tag-to-transaction
		// @Accept  json
		// @Produce  json
		// @Param id path int true "Transaction ID"
		// @Param tag body Tag true "Tag info to add"
		// @Success 200 {object} Tag
		// @Router /transactions/{id}/tags [post]
		transactions.POST("/:id/tags", handlers.AddTagToTransaction)

		// @Summary Remove a tag from a specific transaction
		// @Description Remove a tag from a specific transaction
		// @ID remove-tag-from-transaction
		// @Accept  json
		// @Produce  json
		// @Param id path int true "Transaction ID"
		// @Param tagID path int true "Tag ID"
		// @Success 200 {object} DeleteResponse
		// @Router /transactions/{id}/tags/{tagID} [delete]
		transactions.DELETE("/:id/tags/:tagID", handlers.RemoveTagFromTransaction)
	}
}
