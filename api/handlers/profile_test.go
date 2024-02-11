package handlers

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"
	"xspends/models/interfaces"
	xmock "xspends/models/mock"
	"xspends/testutils"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Initialize mock and context for the tests
func initUserProfileTest(t *testing.T) *xmock.MockUserModel {
	gin.SetMode(gin.TestMode)
	_, modelsService, _, _, tearDown := testutils.SetupModelTestEnvironment(t)
	defer tearDown()

	mockUserModel := new(xmock.MockUserModel)
	modelsService.UserModel = mockUserModel
	return mockUserModel
}

func TestGetUserProfile(t *testing.T) {
	mockUserModel := initUserProfileTest(t)
	defer mockUserModel.AssertExpectations(t)

	tests := []struct {
		name           string
		userID         string
		setupMock      func(userID int64)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:   "Successful retrieval",
			userID: "1",
			setupMock: func(userID int64) {
				mockUserModel.On("GetUserByID", mock.Anything, userID, mock.Anything).Return(&interfaces.User{
					ID:        userID,
					Name:      "John Doe",
					Scope:     1,
					CreatedAt: time.Date(0001, 01, 01, 00, 00, 00, 00, time.UTC),
					UpdatedAt: time.Date(0001, 01, 01, 00, 00, 00, 00, time.UTC),
					Email:     "",
					Username:  "",
					Currency:  "",
					// Other fields as per your User struct
				}, nil).Once()
			},
			expectedStatus: http.StatusOK,
			expectedBody: `{
				"created_at": "0001-01-01T00:00:00Z",
				"currency": "",
				"email": "",
				"user_id": 1,
				"name": "John Doe",
				"scope_id":1,
				"updated_at": "0001-01-01T00:00:00Z",
				"username": ""
			}`,
		},
		{
			name:   "Invalid user ID",
			userID: "1",
			setupMock: func(userID int64) {
				mockUserModel.On("GetUserByID", mock.Anything, userID, mock.Anything).Return(&interfaces.User{}, errors.New(`{"error":"user not found"}`)).Once()
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"error":"user not found"}`,
		},

		// ... other test cases
	}

	//test case for invalid userID format (non int64)

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			userIDInt, _ := strconv.ParseInt(tc.userID, 10, 64)
			c.Set("userID", userIDInt)
			if tc.setupMock != nil {
				tc.setupMock(userIDInt)
			}

			GetUserProfile(c)

			assert.Equal(t, tc.expectedStatus, w.Code)
			assert.JSONEq(t, tc.expectedBody, w.Body.String())
		})
	}

	t.Run("Invalid userID", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Set("userID", "invalid")
		GetUserProfile(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, `{"error":"failed to convert userID to int64"}`, w.Body.String())
	})
}

func TestUpdateUserProfile(t *testing.T) {
	mockUserModel := initUserProfileTest(t)
	defer mockUserModel.AssertExpectations(t)

	tests := []struct {
		name           string
		userID         string
		body           string
		setupMock      func(userID int64)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:   "Successful update",
			userID: "1",
			body:   `{"name":"Jane Doe"}`, // Assume this is a valid JSON for updating the user
			setupMock: func(userID int64) {
				updatedUser := interfaces.User{ID: userID, Name: "Jane Doe"} // Construct expected updated user
				mockUserModel.On("UpdateUser", mock.Anything, &updatedUser, mock.Anything).Return(nil).Once()
			},
			expectedStatus: http.StatusOK,
			expectedBody: `{
				"created_at": "0001-01-01T00:00:00Z",
				"currency": "",
				"email": "",
				"user_id": 1,
				"name": "Jane Doe",
				"scope_id":0,
				"updated_at": "0001-01-01T00:00:00Z",
				"username": ""
			}`,
		},
		{
			name:           "Fail to update",
			userID:         "1",
			body:           `"name":"Jane Doe"`, // pass invalid JSON
			setupMock:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error": "Invalid user json"}`,
		}, {
			name:   "Failed to update",
			userID: "1",
			body:   `{"name":"Jane Doe"}`, // Assume this is a valid JSON for updating the user
			setupMock: func(userID int64) {
				updatedUser := interfaces.User{ID: userID, Name: "Jane Doe"} // Construct expected updated user
				mockUserModel.On("UpdateUser", mock.Anything, &updatedUser, mock.Anything).Return(errors.New("unable to update user")).Once()
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error": "unable to update user"}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest("PUT", "/user", strings.NewReader(tc.body))
			c, _ := gin.CreateTestContext(w)
			c.Request = r

			userIDInt, _ := strconv.ParseInt(tc.userID, 10, 64)
			c.Set("userID", userIDInt)
			if tc.setupMock != nil {
				tc.setupMock(userIDInt)
			}

			UpdateUserProfile(c)

			assert.Equal(t, tc.expectedStatus, w.Code)
			assert.JSONEq(t, tc.expectedBody, w.Body.String())
		})
	}
	t.Run("Invalid userID", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Set("userID", "invalid")
		UpdateUserProfile(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, `{"error":"failed to convert userID to int64"}`, w.Body.String())
	})
}

func TestDeleteUser(t *testing.T) {
	mockUserModel := initUserProfileTest(t)
	defer mockUserModel.AssertExpectations(t)

	tests := []struct {
		name           string
		userID         string
		setupMock      func(userID int64)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:   "Successful deletion",
			userID: "1",
			setupMock: func(userID int64) {
				mockUserModel.On("DeleteUser", mock.Anything, userID, mock.Anything).Return(nil).Once()
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"message":"user deleted successfully"}`,
		}, {
			name:   "Failed deletion",
			userID: "1",
			setupMock: func(userID int64) {
				mockUserModel.On("DeleteUser", mock.Anything, userID, mock.Anything).Return(errors.New("unable to delete user")).Once()
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error": "unable to delete user"}`,
		},
		// ... other test cases
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			userIDInt, _ := strconv.ParseInt(tc.userID, 10, 64)
			c.Set("userID", userIDInt)
			if tc.setupMock != nil {
				tc.setupMock(userIDInt)
			}

			DeleteUser(c)

			assert.Equal(t, tc.expectedStatus, w.Code)
			assert.JSONEq(t, tc.expectedBody, w.Body.String())
		})
	}
	t.Run("Invalid userID", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Set("userID", "invalid")
		DeleteUser(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, `{"error":"failed to convert userID to int64"}`, w.Body.String())
	})
}
