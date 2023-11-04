package handlers

import (
	"log"
	"net/http"
	"xspends/kvstore"

	"github.com/gin-gonic/gin"
)

func TempTestKV(c *gin.Context) {
	client := kvstore.GetClientFromPool()
	err := client.Put(c, []byte("key"), []byte("blob"))
	if err != nil {
		log.Printf("[TempTestKV] Error: could not store the blob: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"KVStore": "inserted Key value pair"})
}

func GetKV(c *gin.Context) {
	client := kvstore.GetClientFromPool()
	value, err := client.Get(c, []byte("key"))
	if err != nil {
		log.Printf("[GetKV] Error: could not read the blob: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	c.JSON(http.StatusOK, string(value))
}
