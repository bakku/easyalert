package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/bakku/easyalert"
	"github.com/bakku/easyalert/random"
)

// CreateUsersHandler should accept a JSON object and create a user from it.
type CreateUsersHandler struct {
	UserRepo easyalert.UserRepository
}

type createUserRequestBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type createUserResponseBody struct {
	Token string `json:"token"`
}

// ServeHTTP handles the HTTP request.
func (h CreateUsersHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not read http body")
		return
	}

	var userBody createUserRequestBody

	err = json.Unmarshal(bytes, &userBody)
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, "invalid json")
		return
	}

	if userBody.Email == "" || userBody.Password == "" {
		writeError(w, http.StatusBadRequest, "Empty email or password.")
		return
	}

	token, err := random.String(easyalert.UserTokenLength)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not generate token")
		return
	}

	user := easyalert.User{
		Email: userBody.Email,
		Token: token,
	}

	err = user.HashPassword(userBody.Password)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not hash password")
		return
	}

	user, err = h.UserRepo.CreateUser(user)
	if err != nil {
		writeError(w, http.StatusBadRequest, fmt.Sprintf("%v", err))
		return
	}

	responseBody := createUserResponseBody{user.Token}

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
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(body))
}

// UpdateUserHandler should accept a JSON object and update the user given by the auth header from it.
type UpdateUserHandler struct {
	UserRepo easyalert.UserRepository
}

type updateUserRequestBody struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}

type updateUserResponseBody struct {
	Email string `json:"email"`
	Token string `json:"token"`
}

// ServeHTTP handles the HTTP request.
func (h UpdateUserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	var userBody updateUserRequestBody

	err = json.Unmarshal(bytes, &userBody)
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, "invalid json")
		return
	}

	if userBody.Email != "" {
		user.Email = userBody.Email
	}

	if userBody.Password != "" {
		user.HashPassword(userBody.Password)
	}

	user, err = h.UserRepo.UpdateUser(user)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not update user")
		return
	}

	responseBody := updateUserResponseBody{user.Email, user.Token}

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
