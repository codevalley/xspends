package models

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
	user, err := GetUserByUsername(ctx, key)
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
	return UpdateUser(ctx, u)
}

func (s *UserStorer) Create(ctx context.Context, user authboss.User) error {
	u, ok := assertUserType(user)
	if !ok {
		return fmt.Errorf("%w: %s", errors.New(userTypeAssertionFailed), "Create")
	}
	return InsertUser(ctx, u)
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
