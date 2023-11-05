package models

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/volatiletech/authboss/v3"
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
	u, ok := user.(*User)
	if !ok {
		log.Printf("[UserStorer Save] Error: user is not of type * User")
		return errors.New("user is not of type * User")
	}
	err := UpdateUser(ctx, u)
	if err != nil {
		log.Printf("[UserStorer Save] Error: %v", err)
	}
	return err
}

func (s *UserStorer) Create(ctx context.Context, user authboss.User) error {
	u, ok := user.(*User)
	if !ok {
		log.Printf("[UserStorer Create] Error: user is not of type * User")
		return errors.New("user is not of type * User")
	}
	err := InsertUser(ctx, u)
	if err != nil {
		log.Printf("[UserStorer Create] Error: %v", err)
	}
	return err
}

func (s *UserStorer) LoadByConfirmSelector(ctx context.Context, selector string) (authboss.ConfirmableUser, error) {
	// Implement this method based on your application's requirements.
	// This method is used for email confirmation features.
	log.Printf("[LoadByConfirmSelector] Error: method not implemented")
	return nil, nil
}

func (s *UserStorer) LoadByRecoverSelector(ctx context.Context, selector string) (authboss.RecoverableUser, error) {
	// Implement this method based on your application's requirements.
	// This method is used for account recovery features.
	log.Printf("[LoadByRecoverSelector] Error: method not implemented")
	return nil, nil
}

// ... and so on for other methods required by different modules of authboss you decide to use.
