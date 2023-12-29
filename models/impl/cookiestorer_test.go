package impl

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"xspends/kvstore/mock"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/volatiletech/authboss/v3"
)

func TestCookieStorer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockKVClient := mock.NewMockRawKVClientInterface(ctrl)
	kvClient := NewCookieStorer(mockKVClient)

	// Test ReadState
	t.Run("ReadState", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.AddCookie(&http.Cookie{Name: authboss.CookieRemember, Value: "test-cookie"})

		state, err := kvClient.ReadState(req)
		assert.NoError(t, err)
		cookieValue, exists := state.Get(authboss.CookieRemember) // capture both value and exists
		assert.True(t, exists)                                    // Assert the cookie was found
		assert.Equal(t, "test-cookie", cookieValue)               // Assert the value is correct
	})

	// Test WriteState
	t.Run("WriteState", func(t *testing.T) {
		w := httptest.NewRecorder()
		state := CustomClientState{authboss.CookieRemember: "test-cookie"}

		events := []authboss.ClientStateEvent{
			{Kind: authboss.ClientStateEventPut, Key: authboss.CookieRemember, Value: "test-cookie"},
		}

		err := kvClient.WriteState(w, state, events)
		assert.NoError(t, err)

		cookie := w.Result().Cookies()[0]
		assert.Equal(t, authboss.CookieRemember, cookie.Name)
		assert.Equal(t, "test-cookie", cookie.Value)
	})

	// Test Load
	t.Run("Load", func(t *testing.T) {
		sid := "session-id"
		mockKVClient.EXPECT().Get(context.Background(), []byte(sid)).Return([]byte("session-data"), nil)

		data, err := kvClient.Load(context.Background(), sid)
		assert.NoError(t, err)
		assert.Equal(t, "session-data", data)
	})

	// Test Save
	t.Run("Save", func(t *testing.T) {
		sid := "session-id"
		state := "session-data"
		mockKVClient.EXPECT().Put(context.Background(), []byte(sid), []byte(state)).Return(nil)

		err := kvClient.Save(context.Background(), sid, state)
		assert.NoError(t, err)
	})

	// Test Delete
	t.Run("Delete", func(t *testing.T) {
		sid := "session-id"
		mockKVClient.EXPECT().Delete(context.Background(), []byte(sid)).Return(nil)

		err := kvClient.Delete(context.Background(), sid)
		assert.NoError(t, err)
	})
}
