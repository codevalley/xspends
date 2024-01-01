package interfaces

import (
	"context"
	"database/sql"
	"strconv"
	"time"
)

type User struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Currency  string    `json:"currency"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Authboss methods,
// have to be implemented to maintain compatibility with the structure of the authboss.User interface
func (u *User) PutPID(pid string) {
	u.ID, _ = strconv.ParseInt(pid, 10, 64)
}

func (u User) GetPID() string {
	return strconv.FormatInt(u.ID, 10)
}

func (u *User) PutPassword(password string) {
	u.Password = password
}

func (u User) GetPassword() string {
	return u.Password
}

type UserService interface {
	InsertUser(ctx context.Context, user *User, otx ...*sql.Tx) error
	UpdateUser(ctx context.Context, user *User, otx ...*sql.Tx) error
	DeleteUser(ctx context.Context, id int64, otx ...*sql.Tx) error
	GetUserByID(ctx context.Context, id int64, otx ...*sql.Tx) (*User, error)
	GetUserByUsername(ctx context.Context, username string, otx ...*sql.Tx) (*User, error)
	UserExists(ctx context.Context, username, email string, otx ...*sql.Tx) (bool, error)
	UserIDExists(ctx context.Context, id int64, otx ...*sql.Tx) (bool, error)
}
