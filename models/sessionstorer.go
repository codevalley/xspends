package models

import (
	"context"
	"net/http"
	"strings"
	"time"
	"xspends/kvstore"

	"github.com/volatiletech/authboss/v3"
)

// SessionStorer implements authboss.ClientStateReadWriter for session management.
type SessionStorer struct {
	kvClient kvstore.RawKVClientInterface
}

// NewSessionStorer creates a new SessionStorer.
func NewSessionStorer(kvClient kvstore.RawKVClientInterface) *SessionStorer {
	return &SessionStorer{
		kvClient: kvClient,
	}
}

// ReadState reads the state from the request headers instead of cookies.
func (s *SessionStorer) ReadState(r *http.Request) (authboss.ClientState, error) {
	token := r.Header.Get("Authorization")
	if token == "" {
		return CustomClientState{}, authboss.ErrUserNotFound // Or your custom error indicating missing token
	}

	// Assuming the token is a bearer token
	splitToken := strings.Split(token, "Bearer ")
	if len(splitToken) != 2 {
		return CustomClientState{}, authboss.ErrUserNotFound // Or your custom error indicating malformed token
	}

	return CustomClientState{"token": splitToken[1]}, nil
}

// WriteState is not used for API calls, as the state is managed client-side.
func (s *SessionStorer) WriteState(w http.ResponseWriter, state authboss.ClientState, events []authboss.ClientStateEvent) error {
	// No action needed for APIs since we don't use cookies.
	return nil
}

// Load retrieves the session data (if any) associated with the session token.
func (s *SessionStorer) Load(ctx context.Context, sid string) (string, error) {
	data, err := s.kvClient.Get(ctx, []byte(sid))
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// Save associates the session data with the session token in the store.
func (s *SessionStorer) Save(ctx context.Context, sid string, state string, ttl time.Duration) error {
	// return s.kvClient.Put(context.Background(), []byte(sid), []byte(state))
	return s.kvClient.PutWithTTL(ctx, []byte(sid), []byte(state), uint64(ttl.Seconds()))
}

// Delete removes the session data associated with the session token from the store.
func (s *SessionStorer) Delete(ctx context.Context, sid string) error {
	return s.kvClient.Delete(ctx, []byte(sid))
}

// Add any other methods necessary for session management here...
