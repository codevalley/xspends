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

package handlers

import (
	"log"
	"net/http"
	"strconv"
	"xspends/models/impl"
	"xspends/models/interfaces"
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

// CreateTransaction
// @Summary Create a new transaction
// @Description Create a new transaction with the provided information
// @ID create-transaction
// @Accept  json
// @Produce  json
// @Param transaction body impl.Transaction true "Transaction info for creation"
// @Success 201 {object} impl.Transaction
// @Failure 400 {object} map[string]string "Invalid transaction data"
// @Failure 500 {object} map[string]string "Unable to create transaction"
// @Router /transactions [post]
func CreateTransaction(c *gin.Context) {
	userID, scopes, ok := getUserAndScopes(c, impl.RoleWrite)
	if !ok {
		log.Printf("[CreateTransaction] Error: Missing user or scope information")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing user or scope information"})
		return
	}

	var newTransaction interfaces.Transaction
	if err := c.ShouldBindJSON(&newTransaction); err != nil {
		log.Printf("[CreateTransaction] Error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	newTransaction.UserID = userID
	newTransaction.ScopeID = scopes[0]
	if err := impl.GetModelsService().TransactionModel.InsertTransaction(c, newTransaction); err != nil {
		log.Printf("[CreateTransaction] Error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to create transaction"})
		return
	}

	c.JSON(http.StatusCreated, newTransaction)
}

// GetTransaction
// @Summary Get a specific transaction
// @Description Get a specific transaction by its ID
// @ID get-transaction
// @Accept  json
// @Produce  json
// @Param id path int true "Transaction ID"
// @Success 200 {object} impl.Transaction
// @Failure 404 {object} map[string]string "Transaction not found"
// @Router /transactions/{id} [get]
func GetTransaction(c *gin.Context) {
	_, scopes, ok := getUserAndScopes(c, impl.RoleView)
	if !ok {
		log.Printf("[GetTransaction] Error: %v", "Missing user or scope information")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing user or scope information"})
		return
	}
	transactionID, ok := getTransactionID(c)
	if !ok {
		log.Printf("[GetTransaction] Error: %v", "Invalid transaction ID")
		return
	}
	transaction, err := impl.GetModelsService().TransactionModel.GetTransactionByID(c, transactionID, scopes)
	if err != nil {
		log.Printf("[GetTransaction] Error: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "transaction not found"})
		return
	}

	c.JSON(http.StatusOK, transaction)
}

// UpdateTransaction
// @Summary Update a specific transaction
// @Description Update a specific transaction by its ID
// @ID update-transaction
// @Accept  json
// @Produce  json
// @Param id path int true "Transaction ID"
// @Param transaction body impl.Transaction true "Transaction info for update"
// @Success 200 {object} impl.Transaction
// @Failure 400 {object} map[string]string "Invalid transaction data"
// @Failure 404 {object} map[string]string "Transaction not found"
// @Failure 500 {object} map[string]string "Unable to update transaction"
// @Router /transactions/{id} [put]
func UpdateTransaction(c *gin.Context) {
	userID, scopes, ok := getUserAndScopes(c, impl.RoleWrite)
	if !ok {
		log.Printf("[UpdateTransaction] Error: %v", "Missing user or scope information")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing user or scope information"})
		return
	}
	transactionID, ok := getTransactionID(c)
	if !ok {
		log.Printf("[UpdateTransaction] Error: %v", "Invalid transaction ID")
		return
	}

	var uTxn interfaces.Transaction
	if err := c.ShouldBindJSON(&uTxn); err != nil {
		log.Printf("[UpdateTransaction] Error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	uTxn.UserID = userID
	uTxn.ScopeID = scopes[0]
	uTxn.ID = transactionID
	oTxn, err := impl.GetModelsService().TransactionModel.GetTransactionByID(c, transactionID, scopes)
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
	if err := impl.GetModelsService().TransactionModel.UpdateTransaction(c, *oTxn); err != nil {
		log.Printf("[UpdateTransaction] Error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to update transaction"})
		return
	}

	c.JSON(http.StatusOK, oTxn)
}

// DeleteTransaction
// @Summary Delete a specific transaction
// @Description Delete a specific transaction by its ID
// @ID delete-transaction
// @Accept  json
// @Produce  json
// @Param id path int true "Transaction ID"
// @Success 200 {object} map[string]string "Message: Transaction deleted successfully"
// @Failure 500 {object} map[string]string "Unable to delete transaction"
// @Router /transactions/{id} [delete]
func DeleteTransaction(c *gin.Context) {

	_, scopes, ok := getUserAndScopes(c, impl.RoleWrite)
	if !ok {
		log.Printf("[DeleteTransaction] Error: %v", "Missing user or scope information")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing user or scope information"})
		return
	}
	transactionID, ok := getTransactionID(c)
	if !ok {
		log.Printf("[DeleteTransaction] Error: %v", "Invalid transaction ID")
		return
	}
	if err := impl.GetModelsService().TransactionModel.DeleteTransaction(c, transactionID, scopes); err != nil {
		log.Printf("[DeleteTransaction] Error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to delete transaction"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "transaction deleted successfully"})
}

// ListTransactions
// @Summary List all transactions
// @Description Get a list of all transactions with optional filters
// @ID list-transactions
// @Accept  json
// @Produce  json
// @Param start_date query string false "Start Date"
// @Param end_date query string false "End Date"
// @Param category query string false "Category"
// @Param type query string false "Transaction Type"
// @Param tags query []string false "Tags"
// @Param min_amount query number false "Minimum Amount"
// @Param max_amount query number false "Maximum Amount"
// @Param sort_by query string false "Sort By"
// @Param sort_order query string false "Sort Order"
// @Param page query int false "Page Number"
// @Param items_per_page query int false "Items Per Page"
// @Success 200 {array} impl.Transaction
// @Failure 500 {object} map[string]string "Unable to fetch transactions"
// @Router /transactions [get]
func ListTransactions(c *gin.Context) {
	userID, scopes, ok := getUserAndScopes(c, impl.RoleView)
	if !ok {
		log.Printf("[ListTransactions] Error: %v", "Missing user or scope information")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing user or scope information"})
		return
	}

	// Create a filter from the query parameters.
	filter := interfaces.TransactionFilter{
		UserID:       userID,
		Scopes:       scopes,
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

	transactions, err := impl.GetModelsService().TransactionModel.GetTransactionsByFilter(c, filter)
	if err != nil {
		log.Printf("[ListTransactions] Error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to fetch transactions"})
		return
	}

	if len(transactions) == 0 {
		log.Printf("[ListTransactions] Info: %v", "no transactions found")
		c.JSON(http.StatusOK, gin.H{"message": "no transactions found"})
		return
	}

	c.JSON(http.StatusOK, transactions)
}
