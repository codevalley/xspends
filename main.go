package main

import (
	"log"
	"xspends/api"
	"xspends/models"
	"xspends/util"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	util.InitializeSnowflake()
	models.InitDB()
	api.SetupRoutes(r)
	log.Println("MATERIAL!@!!!!")
	r.Run() // Defaults to :8080
}
