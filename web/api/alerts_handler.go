package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

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

// GetAlertsHandler should return all alerts of the user.
type GetAlertsHandler struct {
	UserRepo  easyalert.UserRepository
	AlertRepo easyalert.AlertRepository
}

type getAlertsResponseBody struct {
	Subject   string `json:"subject"`
	Status    string `json:"status"`
	SentAt    string `json:"sent_at,omitempty"`
	CreatedAt string `json:"created_at"`
}

// ServeHTTP handles the HTTP request.
func (h GetAlertsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	alerts, err := h.AlertRepo.FindAlerts("WHERE user_id = $1", user.ID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not fetch alerts")
		return
	}

	responseBody := convertAlertsToResponseBody(alerts)

	responseBodyBytes, err := json.Marshal(responseBody)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not marshal response body")
		return
	}

	body, err := prettifyJSON(string(responseBodyBytes))
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not prettify json response")
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(body))
}

func convertAlertsToResponseBody(alerts []easyalert.Alert) []getAlertsResponseBody {
	responseBodyArray := make([]getAlertsResponseBody, len(alerts))

	for i, alert := range alerts {
		responseAlert := getAlertsResponseBody{
			Subject:   alert.Subject,
			Status:    alert.HumanStatus(),
			CreatedAt: alert.CreatedAt.Format(time.RFC3339),
		}

		if alert.SentAt != nil {
			responseAlert.SentAt = alert.SentAt.Format(time.RFC3339)
		}

		responseBodyArray[i] = responseAlert
	}

	return responseBodyArray
}
