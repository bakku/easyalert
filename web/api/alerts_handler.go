package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/bakku/easyalert"
)

// CreateAlertsHandler should accept a JSON object and create an alert from it.
type CreateAlertsHandler struct {
	UserRepo  easyalert.UserRepository
	AlertRepo easyalert.AlertRepository
}

type createAlertRequestBody struct {
	Subject string `json:"subject"`
	Message string `json:"message"`
}

// ServeHTTP handles the HTTP request.
func (h CreateAlertsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	token, ok := getUserToken(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "Missing or invalid Authorization header.")
		return
	}

	user, err := h.UserRepo.FindUser("WHERE token = $1", token)
	if err != nil {
		if err == easyalert.ErrRecordDoesNotExist {
			writeError(w, http.StatusUnauthorized, "Invalid token.")
			return
		}

		writeError(w, http.StatusInternalServerError, "an unknown error occured")
		return
	}

	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not read http body")
		return
	}

	var alertBody createAlertRequestBody

	err = json.Unmarshal(bytes, &alertBody)
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, "invalid json")
		return
	}

	if alertBody.Subject == "" || alertBody.Message == "" {
		writeError(w, http.StatusUnprocessableEntity, "Subject or message not given.")
		return
	}

	alert := easyalert.Alert{
		Subject: alertBody.Subject,
		Status:  0,
		UserID:  user.ID,
	}

	_, err = h.AlertRepo.CreateAlert(alert)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not create alert")
		return
	}

	w.WriteHeader(http.StatusCreated)
}
