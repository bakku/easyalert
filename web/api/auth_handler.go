package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/bakku/easyalert"
	"github.com/bakku/easyalert/random"
)

// AuthHandler accepts a JSON object containing email and password.
// It then validates this object using the database and returns the
// users token if everything is valid.
type AuthHandler struct {
	UserRepo easyalert.UserRepository
}

type authRequestBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type authResponseBody struct {
	Token string `json:"token"`
}

// ServeHTTP handles the HTTP request.
func (h AuthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not read http body")
		return
	}

	var authBody authRequestBody

	err = json.Unmarshal(bytes, &authBody)
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, "invalid json")
		return
	}

	if authBody.Email == "" || authBody.Password == "" {
		writeError(w, http.StatusBadRequest, "Empty email or password.")
		return
	}

	user, err := h.UserRepo.FindUser("WHERE email = $1", authBody.Email)
	if err != nil {
		if err == easyalert.ErrRecordDoesNotExist {
			writeError(w, http.StatusUnauthorized, "Invalid credentials.")
			return
		}

		writeError(w, http.StatusInternalServerError, "an unknown error occured")
		return
	}

	if !user.ValidPassword(authBody.Password) {
		writeError(w, http.StatusUnauthorized, "Invalid credentials.")
		return
	}

	responseBody := authResponseBody{user.Token}

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

// AuthRefreshHandler requires the token of the user in the Authorization
// header and refreshes and returns a new token for the user.
type AuthRefreshHandler struct {
	UserRepo easyalert.UserRepository
}

// ServeHTTP handles the HTTP request.
func (h AuthRefreshHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	newToken, err := random.String(easyalert.UserTokenLength)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not generate token")
		return
	}

	user.Token = newToken

	user, err = h.UserRepo.UpdateUser(user)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not update token")
		return
	}

	var responseBody authResponseBody
	responseBody.Token = user.Token

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
