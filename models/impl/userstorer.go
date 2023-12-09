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
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/volatiletech/authboss/v3"
)

const (
	userTypeAssertionFailed = "user type assertion failed: user is not of type *User"
)

type UserStorer struct {
	db *sql.DB
}

func NewUserStorer(db *sql.DB) *UserStorer {
	return &UserStorer{db: db}
}

func (s *UserStorer) Load(ctx context.Context, key string) (authboss.User, error) {
	user, err := GetUserByUsername(ctx, key, nil) //TODO add DBService / transaction support
	if err != nil {
		log.Printf("[UserStorer Load] Error: %v", err)
		if errors.Is(err, ErrUserNotFound) {
			return nil, authboss.ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}

func (s *UserStorer) Save(ctx context.Context, user authboss.User) error {
	u, ok := assertUserType(user)
	if !ok {
		return fmt.Errorf("%w: %s", errors.New(userTypeAssertionFailed), "Save")
	}
	return UpdateUser(ctx, u, nil) //TODO add DBService / transaction support
}

func (s *UserStorer) Create(ctx context.Context, user authboss.User) error {
	u, ok := assertUserType(user)
	if !ok {
		return fmt.Errorf("%w: %s", errors.New(userTypeAssertionFailed), "Create")
	}

	return InsertUser(ctx, u, nil) //TODO add DBService / transaction support
}

func (s *UserStorer) LoadByConfirmSelector(ctx context.Context, selector string) (authboss.ConfirmableUser, error) {
	log.Printf("[LoadByConfirmSelector] Error: method not implemented")
	return nil, errors.New("LoadByConfirmSelector method not implemented")
}

func (s *UserStorer) LoadByRecoverSelector(ctx context.Context, selector string) (authboss.RecoverableUser, error) {
	log.Printf("[LoadByRecoverSelector] Error: method not implemented")
	return nil, errors.New("LoadByRecoverSelector method not implemented")
}

func assertUserType(user authboss.User) (*User, bool) {
	u, ok := user.(*User)
	if !ok {
		log.Printf("[UserStorer] Error: user is not of type *User")
		return nil, false
	}
	return u, true
}
