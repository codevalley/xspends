package handlers

import (
	"fmt"
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

func initTransactionTagTest(t *testing.T) *xmock.MockTransactionTagModel {
	gin.SetMode(gin.TestMode)
	_, modelsService, _, _, tearDown := testutils.SetupModelTestEnvironment(t)
	defer tearDown()

	mockTransactionTagModel := new(xmock.MockTransactionTagModel)
	modelsService.TransactionTagModel = mockTransactionTagModel
	return mockTransactionTagModel
}

func TestListTransactionTags(t *testing.T) {
	mockTransactionTagModel := initTransactionTagTest(t)
	defer mockTransactionTagModel.AssertExpectations(t)

	tests := []struct {
		name           string
		transactionID  string
		setupMock      func(transactionID int64)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:          "Successful retrieval",
			transactionID: "123",
			setupMock: func(transactionID int64) {
				mockTransactionTagModel.On("GetTagsByTransactionID", mock.AnythingOfType("*gin.Context"), transactionID, mock.AnythingOfType("[]*sql.Tx")).Return([]interfaces.Tag{
					{ID: 1, Name: "Tag1"},
					{ID: 2, Name: "Tag2"},
				}, nil).Once()
			},
			expectedStatus: http.StatusOK,
			expectedBody: `[
				{"tag_id":1, "name":"Tag1", "scope_id":0, "created_at":"0001-01-01T00:00:00Z", "updated_at":"0001-01-01T00:00:00Z", "user_id":0},
				{"tag_id":2, "name":"Tag2", "scope_id":0, "created_at":"0001-01-01T00:00:00Z", "updated_at":"0001-01-01T00:00:00Z", "user_id":0}
			]`,
		},
		// Add more test cases as necessary
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/transactions/"+tc.transactionID+"/tags", nil)
			c, _ := gin.CreateTestContext(w)
			c.Request = r
			c.Params = []gin.Param{{Key: "transaction_id", Value: tc.transactionID}}

			// Convert string transactionID to int64
			transactionIDInt, _ := strconv.ParseInt(tc.transactionID, 10, 64)

			// Setup the mock expectations
			if tc.setupMock != nil {
				tc.setupMock(transactionIDInt)
			}

			// Call the handler
			ListTransactionTags(c)

			// Assert the expectations
			assert.Equal(t, tc.expectedStatus, w.Code)
			assert.JSONEq(t, tc.expectedBody, w.Body.String())
		})
	}
}

func TestAddTagToTransaction(t *testing.T) {
	mockTransactionTagModel := initTransactionTagTest(t)
	defer mockTransactionTagModel.AssertExpectations(t)

	tests := []struct {
		name           string
		transactionID  string
		tagID          string // assume tagID is needed as part of the body or logic
		setupMock      func(transactionID int64, tagID int64)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:          "Successful addition",
			transactionID: "123",
			tagID:         "1",
			setupMock: func(transactionID int64, tagID int64) {
				mockTransactionTagModel.On("InsertTransactionTag", mock.AnythingOfType("*gin.Context"), transactionID, tagID, mock.AnythingOfType("[]*sql.Tx")).Return(nil).Once()
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"message": "tag added successfully to the transaction"}`,
		},
		// Add more test cases as necessary
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest("POST", "/transactions/"+tc.transactionID+"/tags", strings.NewReader(fmt.Sprintf(`{"tag_id":%s}`, tc.tagID)))
			c, _ := gin.CreateTestContext(w)
			c.Request = r
			c.Params = []gin.Param{
				{Key: "transaction_id", Value: tc.transactionID},
				{Key: "tag_id", Value: tc.tagID},
				// Assuming tag ID is part of the URL or body
			}

			transactionIDInt, _ := strconv.ParseInt(tc.transactionID, 10, 64)
			tagIDInt, _ := strconv.ParseInt(tc.tagID, 10, 64)

			if tc.setupMock != nil {
				tc.setupMock(transactionIDInt, tagIDInt)
			}

			AddTagToTransaction(c)

			assert.Equal(t, tc.expectedStatus, w.Code)
			assert.JSONEq(t, tc.expectedBody, w.Body.String())
		})
	}
}

func TestRemoveTagFromTransaction(t *testing.T) {
	mockTransactionTagModel := initTransactionTagTest(t)
	defer mockTransactionTagModel.AssertExpectations(t)

	tests := []struct {
		name           string
		transactionID  string
		tagID          string
		setupMock      func(transactionID int64, tagID int64)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:          "Successful removal",
			transactionID: "123",
			tagID:         "1",
			setupMock: func(transactionID int64, tagID int64) {
				mockTransactionTagModel.On("DeleteTransactionTag", mock.AnythingOfType("*gin.Context"), transactionID, tagID, mock.AnythingOfType("[]*sql.Tx")).Return(nil).Once()
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"message": "tag removed successfully from the transaction"}`,
		},
		// Add more test cases as necessary
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest("DELETE", "/transactions/"+tc.transactionID+"/tags/"+tc.tagID, nil)
			c, _ := gin.CreateTestContext(w)
			c.Request = r
			c.Params = []gin.Param{
				{Key: "transaction_id", Value: tc.transactionID},
				{Key: "id", Value: tc.tagID},
			}

			transactionIDInt, _ := strconv.ParseInt(tc.transactionID, 10, 64)
			tagIDInt, _ := strconv.ParseInt(tc.tagID, 10, 64)

			if tc.setupMock != nil {
				tc.setupMock(transactionIDInt, tagIDInt)
			}

			RemoveTagFromTransaction(c)

			assert.Equal(t, tc.expectedStatus, w.Code)
			assert.JSONEq(t, tc.expectedBody, w.Body.String())
		})
	}
}
