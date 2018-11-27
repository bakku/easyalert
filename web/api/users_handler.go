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
		Admin: false,
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
