package main

import (
	"xspends/api"
	"xspends/models"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	models.InitDB()

	api.SetupRoutes(r)

	r.Run() // Defaults to :8080
}
