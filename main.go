package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	init_db()
	r.POST("/register", register)
	r.POST("/login", login)

	r.Run() // Defaults to :8080
}
