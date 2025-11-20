package user

import (
	"time"

	"github.com/pkg/errors"
)

var (
	ErrUsernameOrPasswordIncorrect = errors.New("username or password incorrect")
)

type User struct {
	ID        string
	Username  string
	Password  string
	CreatedAt time.Time
}
