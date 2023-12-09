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

package impl

import (
	"context"
	"net/http"
	"strings"
	"time"
	"xspends/kvstore"

	"github.com/pkg/errors"
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
		return CustomClientState{}, errors.New("Missing token")
	}

	// Assuming the token is a bearer token
	const bearerPrefix = "Bearer "
	if !strings.HasPrefix(token, bearerPrefix) {
		return CustomClientState{}, errors.New("Malformed token")
	}

	return CustomClientState{"token": token[len(bearerPrefix):]}, nil
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
	if len(data) == 0 {
		return "", nil
	}
	return string(data), nil
}

// Save associates the session data with the session token in the store.
func (s *SessionStorer) Save(ctx context.Context, sid string, state string, ttl time.Duration) error {
	if sid == "" {
		return errors.New("Invalid session ID")
	}
	if state == "" {
		return errors.New("Invalid session state")
	}
	return s.kvClient.Put(ctx, []byte(sid), []byte(state))
	//TODO: not been able to enabled ttl on tikv store.
	//return s.kvClient.PutWithTTL(ctx, []byte(sid), []byte(state), uint64(ttl.Seconds()))
}

// Delete removes the session data associated with the session token from the store.
func (s *SessionStorer) Delete(ctx context.Context, sid string) error {
	return s.kvClient.Delete(ctx, []byte(sid))
}

// Add any other methods necessary for session management here...
