package postgres_test

import (
	"fmt"
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

func TestFindUsers_NoUsers(t *testing.T) {
	db, err := setupDB()
	require.Nil(t, err)

	defer cleanDB(db)

	repo := postgres.UserRepository{DB: db}

	users, err := repo.FindUsers()
	require.Nil(t, err)

	require.Len(t, users, 0)
}
func TestFindUsers_UsersExist(t *testing.T) {
	db, err := setupDB()
	require.Nil(t, err)

	defer cleanDB(db)

	_, err = db.Exec(`
		INSERT INTO users(id, email, password_digest,
			token, admin, created_at, updated_at)
		VALUES (1, 'test@mail.com', '1234',
				'1234', TRUE, NOW(), NOW()),
				(2, 'test@mail2.com', '1234',
				'1235', FALSE, NOW(), NOW()),
				(3, 'test@mail3.com', '1234',
				'1236', FALSE, NOW(), NOW())
	`)
	require.Nil(t, err)

	repo := postgres.UserRepository{DB: db}

	users, err := repo.FindUsers()
	require.Nil(t, err)

	require.Len(t, users, 3)

	first := users[0]
	require.Equal(t, uint(1), first.ID)
	require.Equal(t, "test@mail.com", first.Email)
	require.Equal(t, "1234", first.PasswordDigest)
	require.Equal(t, "1234", first.Token)

	second := users[1]
	require.Equal(t, uint(2), second.ID)
	require.Equal(t, "test@mail2.com", second.Email)
	require.Equal(t, "1234", second.PasswordDigest)
	require.Equal(t, "1235", second.Token)
}

func TestCreateUser_Success(t *testing.T) {
	db, err := setupDB()
	require.Nil(t, err)

	defer cleanDB(db)

	repo := postgres.UserRepository{DB: db}

	user := easyalert.User{Email: "test@user.com", PasswordDigest: "1234", Token: "1234", Admin: false}

	user, err = repo.CreateUser(user)
	require.Nil(t, err)

	var defaultTime time.Time

	require.NotEqual(t, uint(0), user.ID)
	require.NotEqual(t, defaultTime, user.CreatedAt)
	require.NotEqual(t, defaultTime, user.UpdatedAt)

	var exists bool

	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)", user.ID).Scan(&exists)
	require.Nil(t, err)

	require.True(t, exists)
}

func TestCreateUser_ShouldReturnRecordAlreadyExists(t *testing.T) {
	db, err := setupDB()
	require.Nil(t, err)

	defer cleanDB(db)

	_, err = db.Exec(`
			INSERT INTO users(email, password_digest, token, admin, created_at, updated_at)
			VALUES ('test@user.com', '1234', '1234', false, NOW(), NOW())
	`)

	require.Nil(t, err)

	repo := postgres.UserRepository{DB: db}

	user := easyalert.User{Email: "test@user.com", PasswordDigest: "1234", Token: "1234", Admin: false}

	user, err = repo.CreateUser(user)
	require.NotNil(t, err)

	require.Equal(t, "Email is already taken.", fmt.Sprintf("%v", err))
}

func TestUpdateUser_Success(t *testing.T) {
	db, err := setupDB()
	require.Nil(t, err)

	defer cleanDB(db)

	row := db.QueryRow(`
		INSERT INTO users(id, email, password_digest,
			token, admin, created_at, updated_at)
		VALUES (1, 'test@mail.com', '1234',
			'1234', TRUE, NOW(), NOW())
		RETURNING updated_at
	`)

	var oldUpdatedAt time.Time
	err = row.Scan(&oldUpdatedAt)
	require.Nil(t, err)

	user := easyalert.User{}
	user.ID = 1
	user.Email = "updated@mail.com"
	user.PasswordDigest = "5678"
	user.Token = "5678"
	user.Admin = false

	repo := postgres.UserRepository{DB: db}

	user, err = repo.UpdateUser(user)
	require.Nil(t, err)

	require.Equal(t, uint(1), user.ID)
	require.Equal(t, "updated@mail.com", user.Email)
	require.Equal(t, "5678", user.PasswordDigest)
	require.Equal(t, "5678", user.Token)
	require.False(t, user.Admin)
	require.NotEqual(t, oldUpdatedAt, user.UpdatedAt)

	var newUser easyalert.User

	row = db.QueryRow(`
		SELECT id, email, password_digest,
			token, admin, created_at, updated_at
		FROM users
		WHERE id = 1
	`)

	err = row.Scan(&newUser.ID, &newUser.Email, &newUser.PasswordDigest,
		&newUser.Token, &newUser.Admin, &newUser.CreatedAt, &newUser.UpdatedAt)
	require.Nil(t, err)

	require.Equal(t, uint(1), newUser.ID)
	require.Equal(t, "updated@mail.com", newUser.Email)
	require.Equal(t, "5678", newUser.PasswordDigest)
	require.Equal(t, "5678", newUser.Token)
	require.False(t, newUser.Admin)
	require.Equal(t, user.UpdatedAt, newUser.UpdatedAt)
}

func TestUpdateUser_NotExists(t *testing.T) {
	db, err := setupDB()
	require.Nil(t, err)

	defer cleanDB(db)

	user := easyalert.User{}
	user.ID = 1
	user.Email = "updated@mail.com"
	user.PasswordDigest = "5678"
	user.Token = "5678"
	user.Admin = false

	repo := postgres.UserRepository{DB: db}

	_, err = repo.UpdateUser(user)
	require.Equal(t, easyalert.ErrRecordDoesNotExist, err)
}

func TestDeleteUser_Success(t *testing.T) {
	db, err := setupDB()
	require.Nil(t, err)

	defer cleanDB(db)

	row := db.QueryRow(`
		INSERT INTO users(id, email, password_digest,
			token, admin, created_at, updated_at)
		VALUES (1, 'test@mail.com', '1234',
			'1234', TRUE, NOW(), NOW())
		RETURNING updated_at
	`)

	user := easyalert.User{}
	user.ID = 1

	repo := postgres.UserRepository{DB: db}

	err = repo.DeleteUser(user)
	require.Nil(t, err)

	row = db.QueryRow(`
		SELECT COUNT(*)
		FROM users
	`)

	var count int
	err = row.Scan(&count)
	require.Nil(t, err)

	require.Equal(t, 0, count)
}

func TestDeleteUser_NotExists(t *testing.T) {
	db, err := setupDB()
	require.Nil(t, err)

	defer cleanDB(db)

	user := easyalert.User{}
	user.ID = 1

	repo := postgres.UserRepository{DB: db}

	err = repo.DeleteUser(user)
	require.Nil(t, err)
}
