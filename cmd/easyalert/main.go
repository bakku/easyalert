package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/bakku/easyalert/postgres"
	"github.com/bakku/easyalert/web"
	_ "github.com/lib/pq"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		fmt.Println("no PORT env given")
		return
	}

	dbConnStr := os.Getenv("DATABASE_URL")
	if dbConnStr == "" {
		fmt.Println("no DATABASE_URL env given")
		return
	}

	db, err := sql.Open("postgres", dbConnStr)
	if err != nil {
		fmt.Println("error while connecting to database:", err)
		return
	}

	err = db.Ping()
	if err != nil {
		fmt.Println("error while pinging database:", err)
		return
	}

	userRepo := postgres.UserRepository{db}
	alertRepo := postgres.AlertRepository{db}

	server := web.NewServer(port, userRepo, alertRepo)
	server.Start()
}
