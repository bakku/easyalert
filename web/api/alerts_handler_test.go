package api_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/bakku/easyalert"
	"github.com/bakku/easyalert/mocks"
	"github.com/bakku/easyalert/web/api"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestPOSTAlerts_ShouldReturnUnauthorizedIfAuthorizationHeaderIsNotPresent(t *testing.T) {
	req, err := http.NewRequest("POST", "/api/alerts", nil)
	require.Nil(t, err)

	rr := httptest.NewRecorder()
	handler := api.CreateAlertsHandler{}
	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnauthorized, rr.Code)
	require.Equal(t, "application/json; charset=UTF-8", rr.Header().Get("Content-Type"))
	require.Equal(t, "{\n  \"error\": \"Missing or invalid Authorization header.\"\n}", rr.Body.String())
}

func TestPOSTAlerts_ShouldReturnUnauthorizedIfNoUserExistsForToken(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	userRepo := mocks.NewMockUserRepository(mockCtrl)
	userRepo.EXPECT().FindUser(gomock.Any(), gomock.Any()).Return(easyalert.User{}, easyalert.ErrRecordDoesNotExist)

	req, err := http.NewRequest("POST", "/api/alerts", nil)
	require.Nil(t, err)

	req.Header.Set("Authorization", "Bearer 12345")

	rr := httptest.NewRecorder()
	handler := api.CreateAlertsHandler{
		UserRepo: userRepo,
	}
	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnauthorized, rr.Code)
	require.Equal(t, "application/json; charset=UTF-8", rr.Header().Get("Content-Type"))
	require.Equal(t, "{\n  \"error\": \"Invalid token.\"\n}", rr.Body.String())
}

func TestPOSTAlerts_ShouldReturnErrorIfAlertWasGivenIncorrectly(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	userRepo := mocks.NewMockUserRepository(mockCtrl)
	userRepo.EXPECT().FindUser(gomock.Any(), gomock.Any()).Return(easyalert.User{}, nil)

	payload := "invalid"

	req, err := http.NewRequest("POST", "/api/alerts", strings.NewReader(payload))
	require.Nil(t, err)

	req.Header.Set("Authorization", "Bearer 12345")

	rr := httptest.NewRecorder()
	handler := api.CreateAlertsHandler{
		UserRepo: userRepo,
	}
	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
	require.Equal(t, "application/json; charset=UTF-8", rr.Header().Get("Content-Type"))
	require.Equal(t, "{\n  \"error\": \"invalid json\"\n}", rr.Body.String())
}

func TestPOSTAlerts_ShouldReturnErrorIfSubjectOrMessageAreNotGiven(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	userRepo := mocks.NewMockUserRepository(mockCtrl)
	userRepo.EXPECT().FindUser(gomock.Any(), gomock.Any()).Return(easyalert.User{}, nil)

	payload := `{
		"subject": "",
		"message": ""
	}`

	req, err := http.NewRequest("POST", "/api/alerts", strings.NewReader(payload))
	require.Nil(t, err)

	req.Header.Set("Authorization", "Bearer 12345")

	rr := httptest.NewRecorder()
	handler := api.CreateAlertsHandler{
		UserRepo: userRepo,
	}
	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
	require.Equal(t, "application/json; charset=UTF-8", rr.Header().Get("Content-Type"))
	require.Equal(t, "{\n  \"error\": \"Subject or message not given.\"\n}", rr.Body.String())
}

func TestPOSTAlerts_ShouldReturnErrorIfCreationWasUnsuccessful(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	userRepo := mocks.NewMockUserRepository(mockCtrl)
	userRepo.EXPECT().FindUser(gomock.Any(), gomock.Any()).Return(easyalert.User{}, nil)

	alertRepo := mocks.NewMockAlertRepository(mockCtrl)
	alertRepo.EXPECT().CreateAlert(gomock.Any()).Return(easyalert.Alert{}, errors.New("Error!!"))

	payload := `{
		"subject": "Hi",
		"message": "Hi"
	}`

	req, err := http.NewRequest("POST", "/api/alerts", strings.NewReader(payload))
	require.Nil(t, err)

	req.Header.Set("Authorization", "Bearer 12345")

	rr := httptest.NewRecorder()
	handler := api.CreateAlertsHandler{
		UserRepo:  userRepo,
		AlertRepo: alertRepo,
	}
	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
	require.Equal(t, "application/json; charset=UTF-8", rr.Header().Get("Content-Type"))
	require.Equal(t, "{\n  \"error\": \"could not create alert\"\n}", rr.Body.String())
}

func TestPOSTAlerts_ShouldCreateAlert(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	userRepo := mocks.NewMockUserRepository(mockCtrl)
	userRepo.EXPECT().FindUser(gomock.Any(), gomock.Any()).Return(easyalert.User{}, nil)

	alertRepo := mocks.NewMockAlertRepository(mockCtrl)
	alertRepo.EXPECT().CreateAlert(gomock.Any()).Return(easyalert.Alert{}, nil)

	payload := `{
		"subject": "Hi",
		"message": "Hi"
	}`

	req, err := http.NewRequest("POST", "/api/alerts", strings.NewReader(payload))
	require.Nil(t, err)

	req.Header.Set("Authorization", "Bearer 12345")

	rr := httptest.NewRecorder()
	handler := api.CreateAlertsHandler{
		UserRepo:  userRepo,
		AlertRepo: alertRepo,
	}
	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusCreated, rr.Code)
}
