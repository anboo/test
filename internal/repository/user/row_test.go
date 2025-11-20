package user

import (
	"testing"
	"time"

	ent "test-question/internal/entity/user"

	"github.com/stretchr/testify/require"
)

func Test_toEntityUser(t *testing.T) {
	now := time.Now()

	row := &userRow{
		ID:        "uuid-1",
		Username:  "test",
		Password:  "pass",
		CreatedAt: now,
	}

	u := toEntityUser(row)
	require.NotNil(t, u)
	require.Equal(t, "uuid-1", u.ID)
	require.Equal(t, "test", u.Username)
	require.Equal(t, "pass", u.Password)
	require.Equal(t, now, u.CreatedAt)
}

func Test_toEntityUser_nil(t *testing.T) {
	require.Nil(t, toEntityUser(nil))
}

func Test_fromEntityUser(t *testing.T) {
	now := time.Now()

	e := &ent.User{
		ID:        "uuid-2",
		Username:  "hello",
		Password:  "123",
		CreatedAt: now,
	}

	row := fromEntityUser(e)
	require.NotNil(t, row)
	require.Equal(t, "uuid-2", row.ID)
	require.Equal(t, "hello", row.Username)
	require.Equal(t, "123", row.Password)
	require.Equal(t, now, row.CreatedAt)
}

func Test_fromEntityUser_nil(t *testing.T) {
	require.Nil(t, fromEntityUser(nil))
}
