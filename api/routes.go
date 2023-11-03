package api

import (
	"xspends/api/handlers"
	"xspends/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "UP",
		})
	})

	// Authentication routes
	auth := r.Group("/auth")
	{
		auth.POST("/register", handlers.Register)
		auth.POST("/login", handlers.Login)
	}

	// Middleware to authenticate JWT token and fetch user details.
	r.Use(middleware.AuthMiddleware())
	r.Use(middleware.EnsureUserID())

	// Source routes
	sources := r.Group("/sources")
	{
		sources.GET("", handlers.ListSources)
		sources.POST("", handlers.CreateSource)
		sources.GET("/:id", handlers.GetSource)
		sources.PUT("/:id", handlers.UpdateSource)
		sources.DELETE("/:id", handlers.DeleteSource)
	}

	// Category routes
	categories := r.Group("/categories")
	{
		categories.GET("", handlers.ListCategories)
		categories.POST("", handlers.CreateCategory)
		categories.GET("/:id", handlers.GetCategory)
		categories.PUT("/:id", handlers.UpdateCategory)
		categories.DELETE("/:id", handlers.DeleteCategory)
	}

	// Tag routes
	tags := r.Group("/tags")
	{
		tags.GET("", handlers.ListTags)
		tags.POST("", handlers.CreateTag)
		tags.GET("/:id", handlers.GetTag)
		tags.PUT("/:id", handlers.UpdateTag)
		tags.DELETE("/:id", handlers.DeleteTag)
	}

	// Transaction routes
	transactions := r.Group("/transactions")
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
