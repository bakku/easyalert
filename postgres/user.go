package postgres

import (
	"database/sql"

	"github.com/bakku/easyalert"
)

// UserRepository is a postgres implementation of the UserRepository interface
type UserRepository struct {
	DB *sql.DB
}

// FindUser fetches a user by ID and returns it. If the user does not exist it will return easyalert.ErrRecordDoesNotExist.
func (repo *UserRepository) FindUser(ID uint) (easyalert.User, error) {
	var user easyalert.User

	row := repo.DB.QueryRow(`
				SELECT id, email, password_digest, token, admin, created_at, updated_at
				FROM users
				WHERE id = $1
			`, ID)

	err := row.Scan(&user.ID, &user.Email, &user.PasswordDigest, &user.Token, &user.Admin, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return easyalert.User{}, easyalert.ErrRecordDoesNotExist
		}

		return easyalert.User{}, err
	}

	return user, nil
}

func (repo *UserRepository) FindUsers() ([]easyalert.User, error) {
	return nil, nil
}

// CreateUser creates a user in the Postgres database and returns it with ID and created_at/updated_at filled
func (repo *UserRepository) CreateUser(user easyalert.User) (easyalert.User, error) {
	row := repo.DB.QueryRow(`
			INSERT INTO users(email, password_digest, token, admin, created_at, updated_at)
			VALUES ($1, $2, $3, $4, NOW(), NOW())
			RETURNING id, created_at, updated_at
		`, user.Email, user.PasswordDigest, user.Token, user.Admin)

	err := row.Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return easyalert.User{}, err
	}

	return user, nil
}

func (repo *UserRepository) UpdateUser(user easyalert.User) (easyalert.User, error) {
	return easyalert.User{}, nil
}

func (repo *UserRepository) DeleteUser(user easyalert.User) error {
	return nil
}
