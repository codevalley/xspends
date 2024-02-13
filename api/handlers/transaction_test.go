package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"xspends/models/interfaces"
	xmock "xspends/models/mock"
	"xspends/testutils"
	"xspends/util"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func initTransactionTest(t *testing.T) *xmock.MockTransactionModel {
	gin.SetMode(gin.TestMode)
	_, modelsService, _, _, tearDown := testutils.SetupModelTestEnvironment(t)
	defer tearDown()

	mockTransactionModel := new(xmock.MockTransactionModel)
	modelsService.TransactionModel = mockTransactionModel
	return mockTransactionModel
}

func TestCreateTransaction(t *testing.T) {
	mockTransactionModel := initTransactionTest(t)
	defer mockTransactionModel.AssertExpectations(t)

	tests := []struct {
		name           string
		userID         string
		scopeID        string
		requestBody    string
		setupMock      func(userID int64, transaction interfaces.Transaction)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:        "Successful creation",
			userID:      "1",
			scopeID:     "1",
			requestBody: `{"description":"Test transaction","amount":100,"scope_id":1}`,
			setupMock: func(userID int64, transaction interfaces.Transaction) {
				// Adjusting the mock setup to include the third argument
				mockTransactionModel.On("InsertTransaction", mock.AnythingOfType("*gin.Context"), transaction, mock.AnythingOfType("[]*sql.Tx")).Return(nil).Once()
			},
			expectedStatus: http.StatusCreated,
			expectedBody: `{
				"amount": 100,
				"category_id": 0,
				"description": "Test transaction",
				"transaction_id": 0,
				"scope_id":1,
				"source_id": 0,
				"tags": null,
				"timestamp": "0001-01-01T00:00:00Z",
				"type": "",
				"user_id": 1
			}`,
		},
		// ... other test cases ...
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest("POST", "/transactions", strings.NewReader(tc.requestBody))
			c, _ := gin.CreateTestContext(w)
			c.Request = r

			if tc.setupMock != nil {
				// Convert userID to int64 and create a transaction object for the mock setup
				userIDInt, _ := strconv.ParseInt(tc.userID, 10, 64)
				scopeIDInt, _ := strconv.ParseInt(tc.scopeID, 10, 64)
				c.Set("scopeID", scopeIDInt)
				c.Set("userID", userIDInt)
				var transaction interfaces.Transaction
				_ = json.Unmarshal([]byte(tc.requestBody), &transaction) // assuming requestBody is a valid JSON for Transaction
				transaction.UserID = userIDInt

				tc.setupMock(userIDInt, transaction)
			}

			CreateTransaction(c)

			assert.Equal(t, tc.expectedStatus, w.Code)
			assert.JSONEq(t, tc.expectedBody, w.Body.String())
		})
	}
}

func TestGetTransaction(t *testing.T) {
	mockTransactionModel := initTransactionTest(t)
	defer mockTransactionModel.AssertExpectations(t)

	tests := []struct {
		name           string
		userID         string
		scopeID        string
		transactionID  string
		setupMock      func(userID int64, transactionID int64)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:          "Successful retrieval",
			userID:        "1",
			scopeID:       "1",
			transactionID: "123",
			setupMock: func(userID int64, transactionID int64) {
				mockTransactionModel.On("GetTransactionByID", mock.AnythingOfType("*gin.Context"), transactionID, []int64{1}, mock.AnythingOfType("[]*sql.Tx")).Return(&interfaces.Transaction{
					ID:          transactionID,
					UserID:      userID,
					ScopeID:     1,
					Description: "Sample Transaction",
					Amount:      100,
				}, nil).Once()
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"amount":100, "category_id":0, "description":"Sample Transaction", "scope_id":1, "source_id":0, "tags":null, "timestamp":"0001-01-01T00:00:00Z", "type":"", "transaction_id":123, "user_id":1}`,
		},
		// ... other test cases ...
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", fmt.Sprintf("/transactions/%s", tc.transactionID), nil)
			c, _ := gin.CreateTestContext(w)
			c.Request = r
			c.Params = gin.Params{{Key: "id", Value: tc.transactionID}}

			userIDInt, _ := strconv.ParseInt(tc.userID, 10, 64)
			scopeIDInt, _ := strconv.ParseInt(tc.scopeID, 10, 64)
			c.Set("userID", userIDInt)
			c.Set("scopeID", scopeIDInt)
			transactionIDInt, _ := strconv.ParseInt(tc.transactionID, 10, 64)

			if tc.setupMock != nil {
				tc.setupMock(userIDInt, transactionIDInt)
			}

			GetTransaction(c)

			assert.Equal(t, tc.expectedStatus, w.Code)
			assert.JSONEq(t, tc.expectedBody, w.Body.String())
		})
	}
}

func TestUpdateTransaction(t *testing.T) {
	mockTransactionModel := initTransactionTest(t)
	defer mockTransactionModel.AssertExpectations(t)

	tests := []struct {
		name           string
		userID         string
		scopeID        string
		transactionID  string
		requestBody    string
		setupMock      func(userID int64, transactionID int64, updatedTransaction interfaces.Transaction)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:          "Successful update",
			userID:        "1",
			scopeID:       "1",
			transactionID: "123",
			requestBody:   `{"description":"Updated description","amount":150,"scope_id":1}`,
			setupMock: func(userID int64, transactionID int64, updatedTransaction interfaces.Transaction) {
				mockTransactionModel.On("GetTransactionByID", mock.AnythingOfType("*gin.Context"), transactionID, []int64{1}, mock.AnythingOfType("[]*sql.Tx")).Return(&interfaces.Transaction{ID: transactionID, UserID: 1, ScopeID: 1}, nil).Once()
				mockTransactionModel.On("UpdateTransaction", mock.AnythingOfType("*gin.Context"), updatedTransaction, mock.AnythingOfType("[]*sql.Tx")).Return(nil).Once()
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"amount":150, "category_id":0, "description":"Updated description", "transaction_id":123, "scope_id":1, "source_id":0, "tags":null, "timestamp":"0001-01-01T00:00:00Z", "type":"", "user_id":1}`,
		},
		// ... other test cases ...
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest("PUT", fmt.Sprintf("/transactions/%s", tc.transactionID), strings.NewReader(tc.requestBody))
			c, _ := gin.CreateTestContext(w)
			c.Request = r
			c.Params = gin.Params{{Key: "id", Value: tc.transactionID}}

			userIDInt, _ := strconv.ParseInt(tc.userID, 10, 64)
			scopeIDInt, _ := strconv.ParseInt(tc.scopeID, 10, 64)
			c.Set("userID", userIDInt)
			c.Set("scopeID", scopeIDInt)

			transactionIDInt, _ := strconv.ParseInt(tc.transactionID, 10, 64)

			if tc.setupMock != nil {
				var updatedTransaction interfaces.Transaction
				_ = json.Unmarshal([]byte(tc.requestBody), &updatedTransaction)
				updatedTransaction.ID = transactionIDInt
				updatedTransaction.UserID = userIDInt
				tc.setupMock(userIDInt, transactionIDInt, updatedTransaction)
			}

			UpdateTransaction(c)

			assert.Equal(t, tc.expectedStatus, w.Code)
			assert.JSONEq(t, tc.expectedBody, w.Body.String())
		})
	}
}

func TestDeleteTransaction(t *testing.T) {
	mockTransactionModel := initTransactionTest(t)
	defer mockTransactionModel.AssertExpectations(t)

	tests := []struct {
		name           string
		userID         string
		scopeID        string
		transactionID  string
		setupMock      func(userID int64, transactionID int64)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:          "Successful deletion",
			userID:        "1",
			scopeID:       "1",
			transactionID: "123",
			setupMock: func(userID int64, transactionID int64) {
				mockTransactionModel.On("DeleteTransaction", mock.AnythingOfType("*gin.Context"), transactionID, []int64{1}, mock.AnythingOfType("[]*sql.Tx")).Return(nil).Once()
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"message":"transaction deleted successfully"}`,
		},
		// ... other test cases ...
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest("DELETE", fmt.Sprintf("/transactions/%s", tc.transactionID), nil)
			c, _ := gin.CreateTestContext(w)
			c.Request = r
			c.Params = gin.Params{{Key: "id", Value: tc.transactionID}}

			userIDInt, _ := strconv.ParseInt(tc.userID, 10, 64)
			scopeIDInt, _ := strconv.ParseInt(tc.scopeID, 10, 64)
			c.Set("userID", userIDInt)
			c.Set("scopeID", scopeIDInt)
			transactionIDInt, _ := strconv.ParseInt(tc.transactionID, 10, 64)

			if tc.setupMock != nil {
				tc.setupMock(userIDInt, transactionIDInt)
			}

			DeleteTransaction(c)

			assert.Equal(t, tc.expectedStatus, w.Code)
			assert.JSONEq(t, tc.expectedBody, w.Body.String())
		})
	}
}

func TestListTransactions(t *testing.T) {
	mockTransactionModel := initTransactionTest(t)
	defer mockTransactionModel.AssertExpectations(t)

	tests := []struct {
		name           string
		userID         string
		scopeID        string
		queryParams    map[string]string // assuming filters are passed as query parameters
		setupMock      func(userID int64, filter interfaces.TransactionFilter)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:    "Successful retrieval",
			userID:  "1",
			scopeID: "1",
			queryParams: map[string]string{
				"page":           "1",
				"items_per_page": "10",
				// Add other filter parameters as needed
			},
			setupMock: func(userID int64, filter interfaces.TransactionFilter) {
				// Ensure the mock setup matches the actual method call including all parameters
				mockTransactionModel.On("GetTransactionsByFilter", mock.AnythingOfType("*gin.Context"), filter, mock.AnythingOfType("[]*sql.Tx")).Return([]interfaces.Transaction{
					{ID: 1, UserID: userID, Amount: 100, ScopeID: 1, Description: "Transaction 1"},
					{ID: 2, UserID: userID, Amount: 200, ScopeID: 1, Description: "Transaction 2"},
				}, nil).Once()
			},

			expectedStatus: http.StatusOK,
			expectedBody: `[
				{
					"transaction_id": 1,
					"user_id": 1,
					"scope_id": 1,
					"amount": 100,
					"description": "Transaction 1",
					"category_id": 0,
					"source_id": 0,
					"tags": null,
					"timestamp": "0001-01-01T00:00:00Z",
					"type": ""
				},
				{
					"transaction_id": 2,
					"user_id": 1,
					"scope_id": 1,
					"amount": 200,
					"description": "Transaction 2",
					"category_id": 0,
					"source_id": 0,
					"tags": null,
					"timestamp": "0001-01-01T00:00:00Z",
					"type": ""
				}
			]`,
		},
		// ... other test cases ...
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			reqURL := "/transactions?" + generateQueryParams(tc.queryParams)
			r := httptest.NewRequest("GET", reqURL, nil)
			c, _ := gin.CreateTestContext(w)
			c.Request = r

			userIDInt, _ := strconv.ParseInt(tc.userID, 10, 64)
			scopeIDInt, _ := strconv.ParseInt(tc.scopeID, 10, 64)
			c.Set("userID", userIDInt)
			c.Set("scopeID", scopeIDInt)

			if tc.setupMock != nil {
				// Here, convert query parameters into a TransactionFilter and pass to setupMock
				filter := interfaces.TransactionFilter{
					UserID:       userIDInt,
					ScopeID:      scopeIDInt,
					Page:         util.GetIntFromQuery(c, "page", 1),
					ItemsPerPage: util.GetIntFromQuery(c, "items_per_page", 10),
					SortBy:       c.DefaultQuery("sort_by", "timestamp"),
					SortOrder:    c.DefaultQuery("sort_order", "DESC"),
					// Include additional fields as necessary based on query parameters
				}
				tc.setupMock(userIDInt, filter)
			}

			ListTransactions(c)

			assert.Equal(t, tc.expectedStatus, w.Code)
			assert.JSONEq(t, tc.expectedBody, w.Body.String())
		})
	}
}

// Helper function to generate URL query parameters from a map
func generateQueryParams(params map[string]string) string {
	if len(params) == 0 {
		return ""
	}

	var queryParams []string
	for k, v := range params {
		queryParams = append(queryParams, fmt.Sprintf("%s=%s", k, v))
	}
	return strings.Join(queryParams, "&")
}
