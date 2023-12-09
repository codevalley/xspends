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

// ListTransactionTags
// @Summary List all tags for a specific transaction
// @Description Get a list of all tags for a specific transaction
// @ID list-transaction-tags
// @Accept  json
// @Produce  json
// @Param id path int true "Transaction ID"
// @Success 200 {array} impl.Tag
// @Failure 500 {object} map[string]string "Unable to fetch tags for the transaction"
// @Router /transactions/{id}/tags [get]

func ListTransactionTags(c *gin.Context) {
	transactionID, ok := getTxnTagID(c)
	if !ok {
		return
	}

	tags, err := impl.GetTagsByTransactionID(c, transactionID, nil)
	if err != nil {
		log.Printf("[ListTransactionTags] Error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to fetch tags for the transaction"})
		return
	}
	c.JSON(http.StatusOK, tags)
}

// AddTagToTransaction
// @Summary Add a tag to a specific transaction
// @Description Add a tag to a specific transaction
// @ID add-tag-to-transaction
// @Accept  json
// @Produce  json
// @Param id path int true "Transaction ID"
// @Param tag body impl.Tag true "Tag info to add"
// @Success 200 {object} map[string]string "Message: Tag added successfully to the transaction"
// @Failure 400 {object} map[string]string "Invalid tag data"
// @Failure 500 {object} map[string]string "Unable to add tag to the transaction"
// @Router /transactions/{id}/tags [post]
func AddTagToTransaction(c *gin.Context) {
	transactionID, ok := getTxnTagID(c)
	if !ok {
		return
	}

	var tag impl.Tag
	if err := c.ShouldBindJSON(&tag); err != nil {
		log.Printf("[AddTagToTransaction] Error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := impl.InsertTransactionTag(c, transactionID, tag.ID, nil); err != nil {
		log.Printf("[AddTagToTransaction] Error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to add tag to the transaction"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "tag added successfully to the transaction"})
}

// RemoveTagFromTransaction
// @Summary Remove a tag from a specific transaction
// @Description Remove a tag from a specific transaction
// @ID remove-tag-from-transaction
// @Accept  json
// @Produce  json
// @Param id path int true "Transaction ID"
// @Param tagID path int true "Tag ID"
// @Success 200 {object} map[string]string "Message: Tag removed successfully from the transaction"
// @Failure 500 {object} map[string]string "Unable to remove tag from the transaction"
// @Router /transactions/{id}/tags/{tagID} [delete]
func RemoveTagFromTransaction(c *gin.Context) {
	transactionID, ok := getTxnTagID(c)
	if !ok {
		return
	}

	tagID, ok := getTagID(c)
	if !ok {
		return
	}

	if err := impl.DeleteTransactionTag(c, transactionID, tagID, nil); err != nil {
		log.Printf("[RemoveTagFromTransaction] Error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to remove tag from the transaction"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "tag removed successfully from the transaction"})
}
