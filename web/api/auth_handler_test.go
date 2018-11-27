package api_test

import (
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

func TestPOSTAuth_ShouldReturnErrorIfUserWasGivenIncorrectly(t *testing.T) {
	payload := "invalid"

	req, err := http.NewRequest("POST", "/api/auth", strings.NewReader(payload))
	require.Nil(t, err)

	rr := httptest.NewRecorder()
	handler := api.AuthHandler{}
	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
	require.Equal(t, "application/json; charset=UTF-8", rr.Header().Get("Content-Type"))
	require.Equal(t, "{\n  \"error\": \"invalid json\"\n}", rr.Body.String())
}

func TestPOSTAuth_ShouldReturnErrorIfEmailOrPasswordAreEmpty(t *testing.T) {

	payload := `
		{
			"email" : "",
			"password" : ""
		}
	`

	req, err := http.NewRequest("POST", "/api/auth", strings.NewReader(payload))
	require.Nil(t, err)

	rr := httptest.NewRecorder()
	handler := api.AuthHandler{}
	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)
	require.Equal(t, "application/json; charset=UTF-8", rr.Header().Get("Content-Type"))
	require.Equal(t, "{\n  \"error\": \"Empty email or password.\"\n}", rr.Body.String())
}

func TestPOSTAuth_ShouldReturnErrorIfUserCouldNotBeFound(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	userRepo := mocks.NewMockUserRepository(mockCtrl)
	userRepo.EXPECT().FindUser(gomock.Any(), gomock.Any()).Return(easyalert.User{}, easyalert.ErrRecordDoesNotExist)

	payload := `
		{
			"email" : "test@mail.com",
			"password" : "test1234"
		}
	`

	req, err := http.NewRequest("POST", "/api/auth", strings.NewReader(payload))
	require.Nil(t, err)

	rr := httptest.NewRecorder()
	handler := api.AuthHandler{
		UserRepo: userRepo,
	}
	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnauthorized, rr.Code)
	require.Equal(t, "application/json; charset=UTF-8", rr.Header().Get("Content-Type"))
	require.Equal(t, "{\n  \"error\": \"Invalid credentials.\"\n}", rr.Body.String())
}

func TestPOSTAuth_ShouldReturnErrorIfPasswordWasIncorrect(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	user := easyalert.User{
		ID:             1,
		Email:          "test@mail.com",
		PasswordDigest: "12345",
		Token:          "12345",
		Admin:          false,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	userRepo := mocks.NewMockUserRepository(mockCtrl)
	userRepo.EXPECT().FindUser(gomock.Any(), gomock.Any()).Return(user, nil)

	payload := `
		{
			"email" : "test@mail.com",
			"password" : "test1234"
		}
	`

	req, err := http.NewRequest("POST", "/api/auth", strings.NewReader(payload))
	require.Nil(t, err)

	rr := httptest.NewRecorder()
	handler := api.AuthHandler{
		UserRepo: userRepo,
	}
	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnauthorized, rr.Code)
	require.Equal(t, "application/json; charset=UTF-8", rr.Header().Get("Content-Type"))
	require.Equal(t, "{\n  \"error\": \"Invalid credentials.\"\n}", rr.Body.String())
}

func TestPOSTAuth_ShouldReturnTokenIfPasswordWasCorrect(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	user := easyalert.User{
		ID:             1,
		Email:          "test@mail.com",
		PasswordDigest: "$2a$10$zWmZyQoDOafIAOX0RJSsHuyY8DLBb3q9TKbNJQDroF1XVCwtIsamC",
		Token:          "12345",
		Admin:          false,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	userRepo := mocks.NewMockUserRepository(mockCtrl)
	userRepo.EXPECT().FindUser(gomock.Any(), gomock.Any()).Return(user, nil)

	payload := `
		{
			"email" : "test@mail.com",
			"password" : "test1234"
		}
	`

	req, err := http.NewRequest("POST", "/api/auth", strings.NewReader(payload))
	require.Nil(t, err)

	rr := httptest.NewRecorder()
	handler := api.AuthHandler{
		UserRepo: userRepo,
	}
	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	require.Equal(t, "application/json; charset=UTF-8", rr.Header().Get("Content-Type"))
	require.Equal(t, "{\n  \"token\": \"12345\"\n}", rr.Body.String())
}

func TestPUTAuthRefresh_ShouldReturnUnauthorizedIfAuthorizationHeaderIsNotPresent(t *testing.T) {
	req, err := http.NewRequest("PUT", "/api/auth/refresh", nil)
	require.Nil(t, err)

	rr := httptest.NewRecorder()
	handler := api.AuthRefreshHandler{}
	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnauthorized, rr.Code)
	require.Equal(t, "application/json; charset=UTF-8", rr.Header().Get("Content-Type"))
	require.Equal(t, "{\n  \"error\": \"Missing or invalid Authorization header.\"\n}", rr.Body.String())
}

func TestPUTAuthRefresh_ShouldReturnUnauthorizedIfAuthorizationHeaderIsInvalid(t *testing.T) {
	req, err := http.NewRequest("PUT", "/api/auth/refresh", nil)
	require.Nil(t, err)

	req.Header.Set("Authorization", "Invalid")

	rr := httptest.NewRecorder()
	handler := api.AuthRefreshHandler{}
	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnauthorized, rr.Code)
	require.Equal(t, "application/json; charset=UTF-8", rr.Header().Get("Content-Type"))
	require.Equal(t, "{\n  \"error\": \"Missing or invalid Authorization header.\"\n}", rr.Body.String())
}

func TestPUTAuthRefresh_ShouldReturnUnauthorizedIfNoUserExistsForToken(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	userRepo := mocks.NewMockUserRepository(mockCtrl)
	userRepo.EXPECT().FindUser(gomock.Any(), gomock.Any()).Return(easyalert.User{}, easyalert.ErrRecordDoesNotExist)

	req, err := http.NewRequest("PUT", "/api/auth/refresh", nil)
	require.Nil(t, err)

	req.Header.Set("Authorization", "Bearer 12345")

	rr := httptest.NewRecorder()
	handler := api.AuthRefreshHandler{
		UserRepo: userRepo,
	}
	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnauthorized, rr.Code)
	require.Equal(t, "application/json; charset=UTF-8", rr.Header().Get("Content-Type"))
	require.Equal(t, "{\n  \"error\": \"Invalid token.\"\n}", rr.Body.String())
}

func TestPUTAuthRefresh_ShouldRefreshTheToken(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	foundUser := easyalert.User{
		ID:             1,
		Email:          "test@mail.com",
		PasswordDigest: "$2a$10$zWmZyQoDOafIAOX0RJSsHuyY8DLBb3q9TKbNJQDroF1XVCwtIsamC",
		Token:          "12345",
		Admin:          false,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	updatedUser := easyalert.User{
		ID:             1,
		Email:          "test@mail.com",
		PasswordDigest: "$2a$10$zWmZyQoDOafIAOX0RJSsHuyY8DLBb3q9TKbNJQDroF1XVCwtIsamC",
		Token:          "67890",
		Admin:          false,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	userRepo := mocks.NewMockUserRepository(mockCtrl)
	userRepo.EXPECT().FindUser(gomock.Any(), gomock.Any()).Return(foundUser, nil)
	userRepo.EXPECT().UpdateUser(gomock.Any()).Return(updatedUser, nil)

	req, err := http.NewRequest("PUT", "/api/auth/refresh", nil)
	require.Nil(t, err)

	req.Header.Set("Authorization", "Bearer 12345")

	rr := httptest.NewRecorder()
	handler := api.AuthRefreshHandler{
		UserRepo: userRepo,
	}
	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	require.Equal(t, "application/json; charset=UTF-8", rr.Header().Get("Content-Type"))
	require.Equal(t, "{\n  \"token\": \"67890\"\n}", rr.Body.String())
}
