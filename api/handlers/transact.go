package handlers

import (
	"log"
	"net/http"
	"strconv"
	"xspends/models"

	"github.com/gin-gonic/gin"
)

// CreateTransaction creates a new transaction for the authenticated user.
func CreateTransaction(c *gin.Context) {
	userID, _ := c.Get("userID")
	var newTransaction models.Transaction
	if err := c.ShouldBindJSON(&newTransaction); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := models.InsertTransaction(newTransaction); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to create transaction"})
		log.Println(err.Error())
		return
	}

	// Associate tags with the transaction
	if len(newTransaction.Tags) > 0 {
		transactionID, err := strconv.ParseInt(newTransaction.ID, 10, 64)
		if err == nil {
			models.AddTagsToTransaction(transactionID, newTransaction.Tags, userID.(int64))
		}
	}

	c.JSON(http.StatusOK, newTransaction)
}

// GetTransaction fetches a specific transaction by its ID.
func GetTransaction(c *gin.Context) {
	transactionID := c.Param("id")
	userID, _ := c.Get("userID")
	transaction, err := models.GetTransactionByID(transactionID, userID.(string))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "transaction not found"})
		return
	}

	c.JSON(http.StatusOK, transaction)
}

// UpdateTransaction modifies the details of an existing transaction.
func UpdateTransaction(c *gin.Context) {
	userID, _ := c.Get("userID")
	var updatedTransaction models.Transaction
	if err := c.ShouldBindJSON(&updatedTransaction); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := models.UpdateTransaction(updatedTransaction); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to update transaction"})
		return
	}

	// Update the associated tags with the transaction
	transactionID, err := strconv.ParseInt(updatedTransaction.ID, 10, 64)
	if err == nil {
		models.UpdateTagsForTransaction(transactionID, updatedTransaction.Tags, userID.(int64))
	}

	c.JSON(http.StatusOK, updatedTransaction)
}

// DeleteTransaction removes a transaction by its ID.
func DeleteTransaction(c *gin.Context) {
	transactionID := c.Param("id")
	userID, _ := c.Get("userID")
	if err := models.DeleteTransaction(transactionID, userID.(string)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to delete transaction"})
		return
	}

	// Remove all tags associated with the transaction
	transID, err := strconv.ParseInt(transactionID, 10, 64)
	if err == nil {
		models.RemoveTagsFromTransaction(transID)
	}

	c.JSON(http.StatusOK, gin.H{"message": "transaction deleted successfully"})
}

// ListTransactions fetches all transactions for the authenticated user, with optional filters.
func ListTransactions(c *gin.Context) {
	// Create a filter from the query parameters.
	filter := models.TransactionFilter{
		// ... (rest of the filter initialization remains unchanged)
	}

	transactions, err := models.GetTransactionsByFilter(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to fetch transactions"})
		return
	}

	if len(transactions) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "no transactions found"})
		return
	}

	// Fetching associated tags for each transaction
	for idx := range transactions {
		transactionID, err := strconv.ParseInt(transactions[idx].ID, 10, 64)
		if err != nil {
			log.Printf("Error converting transaction ID %s to integer: %v", transactions[idx].ID, err)
			continue
		}

		tags, err := models.GetTagsByTransactionID(transactionID)
		if err != nil {
			// Handle error, perhaps log it and continue to the next iteration
			log.Printf("Error fetching tags for transaction %d: %v", transactionID, err)
			continue
		}

		// Extract tag names from the tags
		tagNames := make([]string, len(tags))
		for i, tag := range tags {
			tagNames[i] = tag.Name
		}

		transactions[idx].Tags = tagNames
	}

	c.JSON(http.StatusOK, transactions)
}

// Helper function to safely retrieve float values from query parameters.
func getFloatFromQuery(c *gin.Context, key string, defaultValue float64) float64 {
	valueStr := c.DefaultQuery(key, "")
	if value, err := strconv.ParseFloat(valueStr, 64); err == nil {
		return value
	}
	return defaultValue
}

// Helper function to safely retrieve integer values from query parameters.
func getIntFromQuery(c *gin.Context, key string, defaultValue int) int {
	valueStr := c.DefaultQuery(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}
