package api

import (
	"net/http"
)

// HomeHandler just returns a simple welcome message
type HomeHandler struct{}

// ServeHTTP handles the HTTP request
func (h HomeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"easyalert":"Alerting made easy"}`))
}
