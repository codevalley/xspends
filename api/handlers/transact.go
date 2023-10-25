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

	userID, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "user not authenticated"})
		return
	}

	var newTransaction models.Transaction
	if err := c.ShouldBindJSON(&newTransaction); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	newTransaction.UserID = userID.(int64)
	if err := models.InsertTransaction(newTransaction); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to create transaction"})
		log.Println(err.Error())
		return
	}

	c.JSON(http.StatusOK, newTransaction)
}

// GetTransaction fetches a specific transaction by its ID.
func GetTransaction(c *gin.Context) {
	transactionID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid transaction ID"})
		return
	}

	userID, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "user not authenticated"})
		return
	}

	transaction, err := models.GetTransactionByID(transactionID, userID.(int64))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "transaction not found"})
		return
	}

	c.JSON(http.StatusOK, transaction)
}

// UpdateTransaction modifies the details of an existing transaction.
func UpdateTransaction(c *gin.Context) {
	userID, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "user not authenticated"})
		return
	}
	// bodyBytes, err := c.GetRawData()
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to retrieve request body"})
	// 	return
	// }
	// bodyString := string(bodyBytes)
	// log.Println(bodyString)
	var updatedTransaction models.Transaction
	if err := c.ShouldBindJSON(&updatedTransaction); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	updatedTransaction.UserID = userID.(int64)

	if err := models.UpdateTransaction(updatedTransaction); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to update transaction"})
		return
	}

	c.JSON(http.StatusOK, updatedTransaction)
}

func DeleteTransaction(c *gin.Context) {
	transactionID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid transaction ID"})
		return
	}
	userID, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "user ID not found in context"})
		return
	}

	if err := models.DeleteTransaction(transactionID, userID.(int64)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "user not authenticated"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "transaction deleted successfully"})
}

// ListTransactions fetches all transactions for the authenticated user, with optional filters.
func ListTransactions(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}
	// Create a filter from the query parameters.
	filter := models.TransactionFilter{
		UserID:       userID.(int64),
		StartDate:    c.DefaultQuery("start_date", ""),
		EndDate:      c.DefaultQuery("end_date", ""),
		Category:     c.DefaultQuery("category", ""),
		Type:         c.DefaultQuery("type", ""),
		Tags:         c.QueryArray("tags"),
		MinAmount:    getFloatFromQuery(c, "min_amount", 0),
		MaxAmount:    getFloatFromQuery(c, "max_amount", 0),
		SortBy:       c.DefaultQuery("sort_by", "timestamp"), // defaulting to timestamp
		SortOrder:    c.DefaultQuery("sort_order", "DESC"),   // defaulting to descending
		Page:         getIntFromQuery(c, "page", 1),
		ItemsPerPage: getIntFromQuery(c, "items_per_page", 10), // defaulting to 10 items per page
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
