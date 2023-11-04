package main

import (
	"xspends/api"
	"xspends/kvstore"
	"xspends/models"
	"xspends/util"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	util.InitializeSnowflake()
	models.InitDB()
	kvstore.SetupKV(false)
	kv := kvstore.GetClientFromPool()
	api.SetupRoutes(r, models.GetDB(), kv)

	r.Run() // Defaults to :8080
}
