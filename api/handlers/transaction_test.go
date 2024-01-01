package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"xspends/models/interfaces"
	xmock "xspends/models/mock"
	"xspends/testutils"

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
		requestBody    string
		setupMock      func(userID int64, transaction interfaces.Transaction)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:        "Successful creation",
			userID:      "1",
			requestBody: `{"description":"Test transaction","amount":100}`,
			setupMock: func(userID int64, transaction interfaces.Transaction) {
				// Adjusting the mock setup to include the third argument
				mockTransactionModel.On("InsertTransaction", mock.AnythingOfType("*gin.Context"), transaction, mock.AnythingOfType("[]*sql.Tx")).Return(nil).Once()
			},
			expectedStatus: http.StatusCreated,
			expectedBody: `{
				"amount": 100,
				"category_id": 0,
				"description": "Test transaction",
				"id": 0,
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

	// ... Add your GetTransaction test cases here
}

func TestUpdateTransaction(t *testing.T) {
	mockTransactionModel := initTransactionTest(t)
	defer mockTransactionModel.AssertExpectations(t)

	// ... Add your UpdateTransaction test cases here
}

func TestDeleteTransaction(t *testing.T) {
	mockTransactionModel := initTransactionTest(t)
	defer mockTransactionModel.AssertExpectations(t)

	// ... Add your DeleteTransaction test cases here
}

func TestListTransactions(t *testing.T) {
	mockTransactionModel := initTransactionTest(t)
	defer mockTransactionModel.AssertExpectations(t)

	// ... Add your ListTransactions test cases here
}
