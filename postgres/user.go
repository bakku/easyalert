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
func (repo UserRepository) FindUser(ID uint) (easyalert.User, error) {
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

// FindUsers fetches all users and returns them.
func (repo UserRepository) FindUsers() ([]easyalert.User, error) {
	var users []easyalert.User

	rows, err := repo.DB.Query(`
				SELECT id, email, password_digest, token
				FROM users
			`)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var u easyalert.User

		if err := rows.Scan(&u.ID, &u.Email, &u.PasswordDigest, &u.Token); err != nil {
			return nil, err
		}

		users = append(users, u)
	}

	return users, nil
}

// CreateUser creates a user in the Postgres database and returns it with ID and created_at/updated_at filled.
func (repo UserRepository) CreateUser(user easyalert.User) (easyalert.User, error) {
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

// UpdateUser updates all fields of the user in the Postgres database and returns it with updated_at refreshed.
func (repo UserRepository) UpdateUser(user easyalert.User) (easyalert.User, error) {
	row := repo.DB.QueryRow(`
			UPDATE users
			SET email = $1, password_digest = $2,
				token = $3, admin = $4, updated_at = NOW()
			WHERE users.id = $5
			RETURNING updated_at
		`, user.Email, user.PasswordDigest, user.Token, user.Admin, user.ID)

	err := row.Scan(&user.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return easyalert.User{}, easyalert.ErrRecordDoesNotExist
		}

		return easyalert.User{}, err
	}

	return user, nil
}

// DeleteUser deletes a user and returns an error if one occurs.
func (repo UserRepository) DeleteUser(user easyalert.User) error {
	_, err := repo.DB.Exec(`
			DELETE FROM users
			WHERE id = $1
		`, user.ID)

	return err
}
