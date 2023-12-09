/*
MIT License

# Copyright (c) 2023 Narayan Babu

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

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
	models.InitModelsService()
	kvstore.SetupKV(context.Background(), false)
	kv := kvstore.GetClientFromPool()
	api.SetupRoutes(r, models.GetDB(), kv) // to be refactored (getDB)

	r.Run() // Defaults to :8080
}
