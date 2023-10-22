package main

import (
	"github.com/gin-gonic/gin"
	//"github.com/go-sql-driver/mysql"
)

func main() {
	r := gin.Default()
	init_db()
	r.POST("/register", register)
	r.POST("/login", login)

	r.Run() // Defaults to :8080
}
