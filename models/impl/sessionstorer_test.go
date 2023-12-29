package impl

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"xspends/kvstore/mock"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestSessionStorer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockKVClient := mock.NewMockRawKVClientInterface(ctrl)
	sessionStorer := NewSessionStorer(mockKVClient)

	t.Run("ReadState", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "Bearer test-token")

		state, err := sessionStorer.ReadState(req)
		assert.NoError(t, err)
		token, exists := state.Get("token")
		assert.True(t, exists)
		assert.Equal(t, "test-token", token)
	})

	t.Run("WriteState", func(t *testing.T) {
		// WriteState should not do anything and not return an error.
		w := httptest.NewRecorder()
		err := sessionStorer.WriteState(w, nil, nil)
		assert.NoError(t, err)
	})

	t.Run("Load", func(t *testing.T) {
		ctx := context.Background()
		sessionID := "session-id"
		mockKVClient.EXPECT().Get(ctx, []byte(sessionID)).Return([]byte("session-data"), nil)

		data, err := sessionStorer.Load(ctx, sessionID)
		assert.NoError(t, err)
		assert.Equal(t, "session-data", data)
	})

	t.Run("Save", func(t *testing.T) {
		ctx := context.Background()
		sessionID := "session-id"
		state := "session-state"
		mockKVClient.EXPECT().Put(ctx, []byte(sessionID), []byte(state)).Return(nil)

		err := sessionStorer.Save(ctx, sessionID, state, 24*time.Hour)
		assert.NoError(t, err)
	})

	t.Run("Delete", func(t *testing.T) {
		ctx := context.Background()
		sessionID := "session-id"
		mockKVClient.EXPECT().Delete(ctx, []byte(sessionID)).Return(nil)

		err := sessionStorer.Delete(ctx, sessionID)
		assert.NoError(t, err)
	})

	// Add additional test cases as needed...
}
