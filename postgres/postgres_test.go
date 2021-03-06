package postgres_test

import (
	"database/sql"
	"os"

	_ "github.com/lib/pq"
)

func setupDB() (*sql.DB, error) {
	connStr := os.Getenv("DATABASE_URL")

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func cleanDB(db *sql.DB) {
	db.Exec("DELETE FROM users")
}
