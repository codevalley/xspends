package middleware

import (
	"context"
	"net/http"
	"time"
	"xspends/kvstore"

	"github.com/volatiletech/authboss/v3"
)

// CustomClientState is a concrete implementation of authboss.ClientState.
type CustomClientState map[string]string

// Get method retrieves the value for a key from the CustomClientState map.
func (cs CustomClientState) Get(key string) (string, bool) {
	value, ok := cs[key]
	return value, ok
}

// CookieStorer implements authboss.ClientStateReadWriter for cookie management.
type CookieStorer struct {
	kvClient kvstore.RawKVClientInterface
}

// NewCookieStorer creates a new CookieStorer.
func NewCookieStorer(kvClient kvstore.RawKVClientInterface) *CookieStorer {
	return &CookieStorer{
		kvClient: kvClient,
	}
}

// ReadState reads the state from the request and returns it as a map.
func (c *CookieStorer) ReadState(r *http.Request) (authboss.ClientState, error) {
	// Extract the remember cookie.
	cookie, err := r.Cookie(authboss.CookieRemember)
	if err != nil {
		if err == http.ErrNoCookie {
			return CustomClientState{}, nil
		}
		return nil, err
	}

	return CustomClientState{authboss.CookieRemember: cookie.Value}, nil
}

// WriteState writes the state to the response.
func (c *CookieStorer) WriteState(w http.ResponseWriter, state authboss.ClientState, events []authboss.ClientStateEvent) error {
	for _, event := range events {
		switch event.Kind {
		case authboss.ClientStateEventPut:
			cookie := &http.Cookie{
				Name:     event.Key,
				Value:    event.Value,
				Path:     "/",
				HttpOnly: true,
				Secure:   true, // Should be set to true if using HTTPS
			}
			http.SetCookie(w, cookie)

		case authboss.ClientStateEventDel:
			cookie := &http.Cookie{
				Name:     event.Key,
				Value:    "",
				Path:     "/",
				HttpOnly: true,
				Secure:   true,            // Should be set to true if using HTTPS
				Expires:  time.Unix(1, 0), // Set to the past to delete the cookie
				MaxAge:   -1,
			}
			http.SetCookie(w, cookie)
		}
	}

	return nil
}

// Load loads the data associated with a cookie from the store.
func (c *CookieStorer) Load(sid string) (string, error) {
	data, err := c.kvClient.Get(context.Background(), []byte(sid))
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// Save saves the data associated with a cookie to the store.
func (c *CookieStorer) Save(sid string, state string) error {
	return c.kvClient.Put(context.Background(), []byte(sid), []byte(state))
}

// Delete deletes the data associated with a cookie from the store.
func (c *CookieStorer) Delete(sid string) error {
	return c.kvClient.Delete(context.Background(), []byte(sid))
}
