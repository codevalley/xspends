package handlers

import (
	"net/http"
	"strconv"
	"xspends/models"

	"github.com/gin-gonic/gin"
)

// ListTransactionTags retrieves all tags for a specific transaction.
func ListTransactionTags(c *gin.Context) {
	transactionIDStr := c.Param("transaction_id")
	transactionID, err := strconv.Atoi(transactionIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid transaction ID"})
		return
	}

	tags, err := models.GetTagsByTransactionID(transactionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to fetch tags for the transaction"})
		return
	}
	c.JSON(http.StatusOK, tags)
}

// AddTagToTransaction adds a new tag to a specific transaction.
func AddTagToTransaction(c *gin.Context) {
	transactionIDStr := c.Param("transaction_id")
	transactionID, err := strconv.Atoi(transactionIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid transaction ID"})
		return
	}

	var tag models.Tag
	if err := c.ShouldBindJSON(&tag); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := models.InsertTransactionTag(transactionID, tag.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to add tag to the transaction"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "tag added successfully to the transaction"})
}

// RemoveTagFromTransaction removes a specific tag from a specific transaction.
func RemoveTagFromTransaction(c *gin.Context) {
	transactionIDStr := c.Param("transaction_id")
	transactionID, err := strconv.Atoi(transactionIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid transaction ID"})
		return
	}

	tagIDStr := c.Param("tag_id")
	tagID, err := strconv.Atoi(tagIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid tag ID"})
		return
	}

	if err := models.DeleteTransactionTag(transactionID, tagID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to remove tag from the transaction"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "tag removed successfully from the transaction"})
}
