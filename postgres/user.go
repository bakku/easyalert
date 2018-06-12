package postgres

import (
	"database/sql"

	"github.com/bakku/easyalert"
)

// UserRepository is a postgres implementation of the UserRepository interface
type UserRepository struct {
	DB *sql.DB
}

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

func (repo *UserRepository) CreateUser(user easyalert.User) (easyalert.User, error) {
	return easyalert.User{}, nil
}

func (repo *UserRepository) UpdateUser(user easyalert.User) (easyalert.User, error) {
	return easyalert.User{}, nil
}

func (repo *UserRepository) DeleteUser(user easyalert.User) error {
	return nil
}
