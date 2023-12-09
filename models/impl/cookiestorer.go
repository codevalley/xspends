/*
MIT License

Copyright (c) 2023 Narayan Babu

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
			return nil, nil
		}
		return nil, err
	}

	return CustomClientState{authboss.CookieRemember: cookie.Value}, nil
}

// WriteState writes the state to the response.
func (c *CookieStorer) WriteState(w http.ResponseWriter, state authboss.ClientState, events []authboss.ClientStateEvent) error {
	const (
		cookiePath     = "/"
		cookieHttpOnly = true
		cookieSecure   = true
	)

	cookie := &http.Cookie{
		Name:     "",
		Value:    "",
		Path:     cookiePath,
		HttpOnly: cookieHttpOnly,
		Secure:   cookieSecure, // Should be set to true if using HTTPS
	}

	for _, event := range events {
		switch event.Kind {
		case authboss.ClientStateEventPut:
			cookie.Name = event.Key
			cookie.Value = event.Value
			http.SetCookie(w, cookie)

		case authboss.ClientStateEventDel:
			cookie.Name = event.Key
			cookie.Value = ""
			cookie.Expires = time.Unix(1, 0)
			cookie.MaxAge = -1
			http.SetCookie(w, cookie)
		}
	}

	return nil
}

// Load loads the data associated with a cookie from the store.
func (c *CookieStorer) Load(ctx context.Context, sid string) (string, error) {
	data, err := c.kvClient.Get(ctx, []byte(sid))
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// Save saves the data associated with a cookie to the store.
func (c *CookieStorer) Save(ctx context.Context, sid string, state string) error {
	return c.kvClient.Put(ctx, []byte(sid), []byte(state))
}

// Delete deletes the data associated with a cookie from the store.
func (c *CookieStorer) Delete(ctx context.Context, sid string) error {
	return c.kvClient.Delete(ctx, []byte(sid))
}
