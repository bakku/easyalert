package api_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/bakku/easyalert"
	"github.com/bakku/easyalert/mocks"
	"github.com/bakku/easyalert/web/api"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestPOSTUsers_ShouldReturnErrorIfUserWasGivenIncorrectly(t *testing.T) {
	payload := "invalid"

	req, err := http.NewRequest("POST", "/api/v1/users", strings.NewReader(payload))
	require.Nil(t, err)

	rr := httptest.NewRecorder()
	handler := api.CreateUsersHandler{}
	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
	require.Equal(t, "application/json; charset=UTF-8", rr.Header().Get("Content-Type"))
	require.Equal(t, "{\n  \"error\": \"invalid json\"\n}", rr.Body.String())
}

func TestPOSTUsers_ShouldReturnErrorFromDatabase(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	userRepo := mocks.NewMockUserRepository(mockCtrl)
	userRepo.EXPECT().CreateUser(gomock.Any()).Return(easyalert.User{}, errors.New("Email is already taken."))

	payload := `
		{
			"email" : "test@mail.com",
			"password" : "test1234"
		}
	`

	req, err := http.NewRequest("POST", "/api/v1/users", strings.NewReader(payload))
	require.Nil(t, err)

	rr := httptest.NewRecorder()
	handler := api.CreateUsersHandler{
		UserRepo: userRepo,
	}
	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)
	require.Equal(t, "application/json; charset=UTF-8", rr.Header().Get("Content-Type"))
	require.Equal(t, "{\n  \"error\": \"Email is already taken.\"\n}", rr.Body.String())
}

func TestPOSTUsers_ShouldNotAcceptEmptyEmailOrPassword(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	userRepo := mocks.NewMockUserRepository(mockCtrl)

	payload := `
		{
			"email" : "",
			"password" : ""
		}
	`

	req, err := http.NewRequest("POST", "/api/v1/users", strings.NewReader(payload))
	require.Nil(t, err)

	rr := httptest.NewRecorder()
	handler := api.CreateUsersHandler{
		UserRepo: userRepo,
	}
	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)
	require.Equal(t, "application/json; charset=UTF-8", rr.Header().Get("Content-Type"))
	require.Equal(t, "{\n  \"error\": \"Empty email or password.\"\n}", rr.Body.String())
}

func TestPOSTUsers_ShouldCreateUser(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	user := easyalert.User{
		ID:             1,
		Email:          "test@mail.com",
		PasswordDigest: "12345",
		Token:          "12345",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	userRepo := mocks.NewMockUserRepository(mockCtrl)
	userRepo.EXPECT().CreateUser(gomock.Any()).Return(user, nil)

	payload := `
		{
			"email" : "test@mail.com",
			"password" : "test1234"
		}
	`

	req, err := http.NewRequest("POST", "/api/v1/users", strings.NewReader(payload))
	require.Nil(t, err)

	rr := httptest.NewRecorder()
	handler := api.CreateUsersHandler{
		UserRepo: userRepo,
	}
	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusCreated, rr.Code)
	require.Equal(t, "application/json; charset=UTF-8", rr.Header().Get("Content-Type"))
	require.Equal(t, "{\n  \"token\": \"12345\"\n}", rr.Body.String())
}
