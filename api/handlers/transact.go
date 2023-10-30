package handlers

import (
	"log"
	"net/http"
	"strconv"
	"xspends/models"
	"xspends/util"

	"github.com/gin-gonic/gin"
)

// CreateTransaction creates a new transaction for the authenticated user.
func CreateTransaction(c *gin.Context) {
	userID := c.MustGet("userID").(int64)

	var newTransaction models.Transaction
	if err := c.ShouldBindJSON(&newTransaction); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	newTransaction.UserID = userID
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
	userID := c.MustGet("userID").(int64)

	transaction, err := models.GetTransactionByID(transactionID, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "transaction not found"})
		return
	}

	c.JSON(http.StatusOK, transaction)
}

// UpdateTransaction modifies the details of an existing transaction.
func UpdateTransaction(c *gin.Context) {
	userID := c.MustGet("userID").(int64)
	// bodyBytes, err := c.GetRawData()
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to retrieve request body"})
	// 	return
	// }
	// bodyString := string(bodyBytes)
	// log.Println(bodyString)
	transactionID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid transaction ID"})
		return
	}
	var uTxn models.Transaction
	if err := c.ShouldBindJSON(&uTxn); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	uTxn.UserID = userID
	uTxn.ID = transactionID
	oTxn, err := models.GetTransactionByID(transactionID, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "unable to find transaction:"})
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
	if err := models.UpdateTransaction(*oTxn); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to update transaction:" + err.Error()})
		log.Print(uTxn)
		return
	}

	c.JSON(http.StatusOK, oTxn)
}

func DeleteTransaction(c *gin.Context) {
	transactionID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid transaction ID"})
		return
	}
	userID := c.MustGet("userID").(int64)
	if err := models.DeleteTransaction(transactionID, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "user not authenticated"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "transaction deleted successfully"})
}

// ListTransactions fetches all transactions for the authenticated user, with optional filters.
func ListTransactions(c *gin.Context) {
	userID := c.MustGet("userID").(int64)

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
