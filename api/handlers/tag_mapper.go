package handlers

import (
	"log"
	"net/http"
	"strconv"
	"xspends/models"

	"github.com/gin-gonic/gin"
)

func getTxnTagID(c *gin.Context) (int64, bool) {
	transactionIDStr := c.Param("transaction_id")
	if transactionIDStr == "" {
		log.Printf("[getTagTransactionID] Error: transaction ID is required")
		c.JSON(http.StatusBadRequest, gin.H{"error": "transaction ID is required"})
		return 0, false
	}

	transactionID, err := strconv.ParseInt(transactionIDStr, 10, 64)
	if err != nil {
		log.Printf("[getTagTransactionID] Error: invalid transaction ID format")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid transaction ID format"})
		return 0, false
	}

	return transactionID, true
}

func ListTransactionTags(c *gin.Context) {
	transactionID, ok := getTxnTagID(c)
	if !ok {
		return
	}

	tags, err := models.GetTagsByTransactionID(c, transactionID)
	if err != nil {
		log.Printf("[ListTransactionTags] Error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to fetch tags for the transaction"})
		return
	}
	c.JSON(http.StatusOK, tags)
}

func AddTagToTransaction(c *gin.Context) {
	transactionID, ok := getTxnTagID(c)
	if !ok {
		return
	}

	var tag models.Tag
	if err := c.ShouldBindJSON(&tag); err != nil {
		log.Printf("[AddTagToTransaction] Error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := models.InsertTransactionTag(c, transactionID, tag.ID); err != nil {
		log.Printf("[AddTagToTransaction] Error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to add tag to the transaction"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "tag added successfully to the transaction"})
}

func RemoveTagFromTransaction(c *gin.Context) {
	transactionID, ok := getTxnTagID(c)
	if !ok {
		return
	}

	tagID, ok := getTagID(c)
	if !ok {
		return
	}

	if err := models.DeleteTransactionTag(c, transactionID, tagID); err != nil {
		log.Printf("[RemoveTagFromTransaction] Error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to remove tag from the transaction"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "tag removed successfully from the transaction"})
}
