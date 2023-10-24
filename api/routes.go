package api

import (
	"xspends/api/handlers"
	"xspends/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	// Authentication routes
	r.POST("/register", handlers.Register)
	r.POST("/login", handlers.Login)

	// Middleware to authenticate JWT token and fetch user details.
	r.Use(middleware.AuthMiddleware())

	// Source routes
	r.GET("/sources", handlers.ListSources)
	r.POST("/sources", handlers.CreateSource)
	r.GET("/sources/:id", handlers.GetSource)
	r.PUT("/sources/:id", handlers.UpdateSource)
	r.DELETE("/sources/:id", handlers.DeleteSource)
	// Category routes
	r.GET("/categories", handlers.ListCategories)
	r.POST("/categories", handlers.CreateCategory)
	r.GET("/categories/:id", handlers.GetCategory)
	r.PUT("/categories/:id", handlers.UpdateCategory)
	r.DELETE("/categories/:id", handlers.DeleteCategory)

	// Tag routes
	r.GET("/tags", handlers.ListTags)
	r.POST("/tags", handlers.CreateTag)
	r.GET("/tags/:id", handlers.GetTag)
	r.PUT("/tags/:id", handlers.UpdateTag)
	r.DELETE("/tags/:id", handlers.DeleteTag)

	// Transaction routes
	r.GET("/transactions", handlers.ListTransactions)
	r.POST("/transactions", handlers.CreateTransaction)
	r.GET("/transactions/:id", handlers.GetTransaction)
	r.PUT("/transactions/:id", handlers.UpdateTransaction)
	r.DELETE("/transactions/:id", handlers.DeleteTransaction)

	// Transaction Tags routes
	r.GET("/transactions/:id/tags", handlers.ListTransactionTags)
	r.POST("/transactions/:id/tags", handlers.AddTagToTransaction)
	r.DELETE("/transactions/:id/tags/:tagID", handlers.RemoveTagFromTransaction)

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "UP",
		})
	})
}
