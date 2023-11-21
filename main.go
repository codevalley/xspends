package main

import (
	"context"
	"xspends/api"
	"xspends/kvstore"
	"xspends/models"
	"xspends/util"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all the routes for the application.
// It takes a gin Engine, a sql DB, and a kvstore RawKVClientInterface as parameters.
// @title XSpends API
// @version 1.0
// @description This is the API for the XSpends application.
// @host localhost:8080
// @BasePath /
func main() {
	r := gin.Default()
	util.InitializeSnowflake()
	models.InitDB()
	kvstore.SetupKV(context.Background(), false)
	kv := kvstore.GetClientFromPool()
	api.SetupRoutes(r, models.GetDB(), kv)

	r.Run() // Defaults to :8080
}
