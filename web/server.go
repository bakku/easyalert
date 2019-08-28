package web

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bakku/easyalert"
	"github.com/bakku/easyalert/web/api"
	"github.com/gorilla/mux"
)

// Server holds everything needed to use easyalert with a HTTP server
type Server struct {
	server http.Server
}

// NewServer returns a new Server with all routes set up
func NewServer(port string, userRepo easyalert.UserRepository, alertRepo easyalert.AlertRepository) *Server {
	s := &Server{
		server: http.Server{
			Addr: ":" + port,
		},
	}

	router := mux.NewRouter()

	// api handler
	home := api.HomeHandler{}
	createUsers := api.CreateUsersHandler{userRepo}
	updateUser := api.UpdateUserHandler{userRepo}

	getAlerts := api.GetAlertsHandler{userRepo, alertRepo}
	createAlerts := api.CreateAlertsHandler{userRepo, alertRepo}

	auth := api.AuthHandler{userRepo}
	authRefresh := api.AuthRefreshHandler{userRepo}

	router.Methods("GET").Path("/api").Handler(home)

	router.Methods("POST").Path("/api/users").Handler(createUsers)
	router.Methods("PUT").Path("/api/users/me").Handler(updateUser)

	router.Methods("GET").Path("/api/alerts").Handler(getAlerts)
	router.Methods("POST").Path("/api/alerts").Handler(createAlerts)

	router.Methods("POST").Path("/api/auth").Handler(auth)
	router.Methods("PUT").Path("/api/auth/refresh").Handler(authRefresh)

	s.server.Handler = router

	return s
}

// Start starts the HTTP server with graceful shutdown implemented
func (s *Server) Start() {
	shutDownSig := make(chan os.Signal, 1)
	shutDownFinished := make(chan bool)

	go func() {
		log.Fatal(s.server.ListenAndServe())
		shutDownFinished <- true
	}()

	signal.Notify(shutDownSig, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-shutDownSig
		log.Printf("Shutting down gracefully...")

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err := s.server.Shutdown(ctx)
		if err != nil {
			log.Println("Shutdown error:", err)
		}
	}()

	// Wait for shutdown to finish
	<-shutDownFinished
	log.Println("Done")
}
