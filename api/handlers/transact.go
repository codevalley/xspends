package handlers

import (
	"log"
	"net/http"
	"strconv"
	"xspends/models"
	"xspends/util"

	"github.com/gin-gonic/gin"
)

func getTransactionID(c *gin.Context) (int64, bool) {
	transactionIDStr := c.Param("id")
	if transactionIDStr == "" {
		log.Printf("[getTransactionID] Error: transaction ID is required")
		c.JSON(http.StatusBadRequest, gin.H{"error": "transaction ID is required"})
		return 0, false
	}

	transactionID, err := strconv.ParseInt(transactionIDStr, 10, 64)
	if err != nil {
		log.Printf("[getTransactionID] Error: invalid transaction ID format")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid transaction ID format"})
		return 0, false
	}

	return transactionID, true
}

// CreateTransaction creates a new transaction for the authenticated user.
func CreateTransaction(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	var newTransaction models.Transaction
	if err := c.ShouldBindJSON(&newTransaction); err != nil {
		log.Printf("[CreateTransaction] Error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	newTransaction.UserID = userID
	if err := models.InsertTransaction(c, newTransaction); err != nil {
		log.Printf("[CreateTransaction] Error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to create transaction"})
		return
	}

	c.JSON(http.StatusOK, newTransaction)
}

// GetTransaction fetches a specific transaction by its ID.
func GetTransaction(c *gin.Context) {
	transactionID, ok := getTransactionID(c)
	if !ok {
		return
	}
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	transaction, err := models.GetTransactionByID(c, transactionID, userID)
	if err != nil {
		log.Printf("[GetTransaction] Error: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "transaction not found"})
		return
	}

	c.JSON(http.StatusOK, transaction)
}

// UpdateTransaction modifies the details of an existing transaction.
func UpdateTransaction(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}
	transactionID, ok := getTransactionID(c)
	if !ok {
		return
	}

	var uTxn models.Transaction
	if err := c.ShouldBindJSON(&uTxn); err != nil {
		log.Printf("[UpdateTransaction] Error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	uTxn.UserID = userID
	uTxn.ID = transactionID
	oTxn, err := models.GetTransactionByID(c, transactionID, userID)
	if err != nil {
		log.Printf("[UpdateTransaction] Error: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "unable to find transaction"})
		return
	}
	if uTxn.Amount != 0 {
		oTxn.Amount = uTxn.Amount
	}
	if uTxn.Description != "" {
		oTxn.Description = uTxn.Description
	}
	if uTxn.Tags != nil {
		oTxn.Tags = uTxn.Tags
	}
	if uTxn.Type != "" {
		oTxn.Type = uTxn.Type
	}
	if uTxn.SourceID != 0 {
		oTxn.SourceID = uTxn.SourceID
	}
	if uTxn.CategoryID != 0 {
		oTxn.CategoryID = uTxn.CategoryID
	}
	if err := models.UpdateTransaction(c, *oTxn); err != nil {
		log.Printf("[UpdateTransaction] Error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to update transaction"})
		return
	}

	c.JSON(http.StatusOK, oTxn)
}

func DeleteTransaction(c *gin.Context) {
	transactionID, ok := getTransactionID(c)
	if !ok {
		return
	}
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	if err := models.DeleteTransaction(c, transactionID, userID); err != nil {
		log.Printf("[DeleteTransaction] Error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to delete transaction"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "transaction deleted successfully"})
}

// ListTransactions fetches all transactions for the authenticated user, with optional filters.
func ListTransactions(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	// Create a filter from the query parameters.
	filter := models.TransactionFilter{
		UserID:       userID,
		StartDate:    c.DefaultQuery("start_date", ""),
		EndDate:      c.DefaultQuery("end_date", ""),
		Category:     c.DefaultQuery("category", ""),
		Type:         c.DefaultQuery("type", ""),
		Tags:         c.QueryArray("tags"),
		MinAmount:    util.GetFloatFromQuery(c, "min_amount", 0),
		MaxAmount:    util.GetFloatFromQuery(c, "max_amount", 0),
		SortBy:       c.DefaultQuery("sort_by", "timestamp"), // defaulting to timestamp
		SortOrder:    c.DefaultQuery("sort_order", "DESC"),   // defaulting to descending
		Page:         util.GetIntFromQuery(c, "page", 1),
		ItemsPerPage: util.GetIntFromQuery(c, "items_per_page", 10), // defaulting to 10 items per page
	}

	transactions, err := models.GetTransactionsByFilter(c, filter)
	if err != nil {
		log.Printf("[ListTransactions] Error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to fetch transactions"})
		return
	}

	if len(transactions) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "no transactions found"})
		return
	}

	c.JSON(http.StatusOK, transactions)
}
