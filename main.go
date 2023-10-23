package main

import (
	"xspends/api/handlers"
	"xspends/models"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	models.InitializeStore()
	r.POST("/register", handlers.Register)
	r.POST("/login", handlers.Login)

	r.Run() // Defaults to :8080
}
