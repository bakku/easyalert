package postgres_test

import (
	"testing"
	"time"

	"github.com/bakku/easyalert"

	"github.com/bakku/easyalert/postgres"
	"github.com/stretchr/testify/require"
)

func TestFindUser_Success(t *testing.T) {
	db, err := setupDB()
	require.Nil(t, err)

	defer cleanDB(db)

	_, err = db.Exec(`
		INSERT INTO users(id, email, password_digest, 
			token, admin, created_at, updated_at)
		VALUES (1, 'test@mail.com', '1234', 
			'1234', TRUE, NOW(), NOW())
	`)
	require.Nil(t, err)

	repo := postgres.UserRepository{DB: db}

	user, err := repo.FindUser(1)
	require.Nil(t, err)

	var defaultTime time.Time

	require.Equal(t, uint(1), user.ID)
	require.Equal(t, "test@mail.com", user.Email)
	require.Equal(t, "1234", user.PasswordDigest)
	require.Equal(t, "1234", user.Token)
	require.True(t, user.Admin)
	require.NotEqual(t, defaultTime, user.CreatedAt)
	require.NotEqual(t, defaultTime, user.UpdatedAt)
}

func TestFindUser_NotExists(t *testing.T) {
	db, err := setupDB()
	require.Nil(t, err)

	defer cleanDB(db)

	repo := postgres.UserRepository{DB: db}

	_, err = repo.FindUser(1)

	require.Equal(t, easyalert.ErrRecordDoesNotExist, err)
}