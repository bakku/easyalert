package main

import (
	"fmt"
	"os"

	"github.com/bakku/easyalert/web"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		fmt.Println("no PORT env given")
		return
	}

	server := web.NewServer(port)
	server.Start()
}
