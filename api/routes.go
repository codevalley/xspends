package api

import (
	"database/sql"
	"xspends/api/handlers"
	"xspends/kvstore"
	"xspends/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, db *sql.DB, kvClient kvstore.RawKVClientInterface) {
	ab := middleware.SetupAuthBoss(r, db, kvClient)

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "UP"})
	})

	// Authentication routes using custom JWT handlers
	auth := r.Group("/auth")
	{
		// Use your custom JWTRegisterHandler and JWTLoginHandler here
		auth.POST("/register", handlers.JWTRegisterHandler(ab))
		auth.POST("/login", handlers.JWTLoginHandler(ab))
		// If you have logout and refresh token routes, you can add them here as well
		// auth.POST("/logout", handlers.JWTLogoutHandler(ab))
		// auth.POST("/refresh", handlers.JWTRefreshHandler(ab))
	}

	// All other routes should be protected by the AuthMiddleware
	apiRoutes := r.Group("/")
	apiRoutes.Use(middleware.AuthMiddleware(ab), middleware.EnsureUserID())

	// Source routes
	sources := apiRoutes.Group("/sources")
	{
		sources.GET("", handlers.ListSources)
		sources.POST("", handlers.CreateSource)
		sources.GET("/:id", handlers.GetSource)
		sources.PUT("/:id", handlers.UpdateSource)
		sources.DELETE("/:id", handlers.DeleteSource)
	}

	// Category routes
	categories := apiRoutes.Group("/categories")
	{
		categories.GET("", handlers.ListCategories)
		categories.POST("", handlers.CreateCategory)
		categories.GET("/:id", handlers.GetCategory)
		categories.PUT("/:id", handlers.UpdateCategory)
		categories.DELETE("/:id", handlers.DeleteCategory)
	}

	// Tag routes
	tags := apiRoutes.Group("/tags")
	{
		tags.GET("", handlers.ListTags)
		tags.POST("", handlers.CreateTag)
		tags.GET("/:id", handlers.GetTag)
		tags.PUT("/:id", handlers.UpdateTag)
		tags.DELETE("/:id", handlers.DeleteTag)
	}

	// Transaction routes
	transactions := apiRoutes.Group("/transactions")
	{
		transactions.GET("", handlers.ListTransactions)
		transactions.POST("", handlers.CreateTransaction)
		transactions.GET("/:id", handlers.GetTransaction)
		transactions.PUT("/:id", handlers.UpdateTransaction)
		transactions.DELETE("/:id", handlers.DeleteTransaction)

		// Transaction Tags routes
		transactions.GET("/:id/tags", handlers.ListTransactionTags)
		transactions.POST("/:id/tags", handlers.AddTagToTransaction)
		transactions.DELETE("/:id/tags/:tagID", handlers.RemoveTagFromTransaction)
	}
}
