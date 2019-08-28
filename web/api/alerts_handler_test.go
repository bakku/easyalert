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

func TestGETAlerts_ShouldReturnUnauthorizedIfAuthorizationHeaderIsNotPresent(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/alerts", nil)
	require.Nil(t, err)

	rr := httptest.NewRecorder()
	handler := api.GetAlertsHandler{}
	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnauthorized, rr.Code)
	require.Equal(t, "application/json; charset=UTF-8", rr.Header().Get("Content-Type"))
	require.Equal(t, "{\n  \"error\": \"Missing or invalid Authorization header.\"\n}", rr.Body.String())
}

func TestGETAlerts_ShouldReturnUnauthorizedIfNoUserExistsForToken(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	userRepo := mocks.NewMockUserRepository(mockCtrl)
	userRepo.EXPECT().FindUser(gomock.Any(), gomock.Any()).Return(easyalert.User{}, easyalert.ErrRecordDoesNotExist)

	req, err := http.NewRequest("GET", "/api/alerts", nil)
	require.Nil(t, err)

	req.Header.Set("Authorization", "Bearer 12345")

	rr := httptest.NewRecorder()
	handler := api.GetAlertsHandler{
		UserRepo: userRepo,
	}
	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnauthorized, rr.Code)
	require.Equal(t, "application/json; charset=UTF-8", rr.Header().Get("Content-Type"))
	require.Equal(t, "{\n  \"error\": \"Invalid token.\"\n}", rr.Body.String())
}

func TestGETAlerts_ShouldReturnErrorIfGettingAlertsReturnsError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	userRepo := mocks.NewMockUserRepository(mockCtrl)
	userRepo.EXPECT().FindUser(gomock.Any(), gomock.Any()).Return(easyalert.User{ID: 1}, nil)

	alertRepo := mocks.NewMockAlertRepository(mockCtrl)
	alertRepo.EXPECT().FindAlerts(gomock.Any(), gomock.Any()).Return(nil, errors.New("Error!!"))

	req, err := http.NewRequest("GET", "/api/alerts", nil)
	require.Nil(t, err)

	req.Header.Set("Authorization", "Bearer 12345")

	rr := httptest.NewRecorder()
	handler := api.GetAlertsHandler{
		UserRepo:  userRepo,
		AlertRepo: alertRepo,
	}
	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
	require.Equal(t, "application/json; charset=UTF-8", rr.Header().Get("Content-Type"))
	require.Equal(t, "{\n  \"error\": \"could not fetch alerts\"\n}", rr.Body.String())
}

func TestGETAlerts_ShouldReturnAllAlerts(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	userRepo := mocks.NewMockUserRepository(mockCtrl)
	userRepo.EXPECT().FindUser(gomock.Any(), gomock.Any()).Return(easyalert.User{ID: 1}, nil)

	createdAt := time.Date(2018, 5, 10, 8, 50, 0, 0, time.UTC)
	sentAt := time.Date(2018, 5, 10, 8, 53, 0, 0, time.UTC)

	expected := []easyalert.Alert{
		{
			ID:        1,
			Subject:   "Test #1",
			Status:    0,
			SentAt:    nil,
			CreatedAt: createdAt,
		},
		{
			ID:        2,
			Subject:   "Test #2",
			Status:    1,
			SentAt:    &sentAt,
			CreatedAt: createdAt,
		},
	}

	alertRepo := mocks.NewMockAlertRepository(mockCtrl)
	alertRepo.EXPECT().FindAlerts(gomock.Any(), gomock.Any()).Return(expected, nil)

	req, err := http.NewRequest("GET", "/api/alerts", nil)
	require.Nil(t, err)

	req.Header.Set("Authorization", "Bearer 12345")

	rr := httptest.NewRecorder()
	handler := api.GetAlertsHandler{
		UserRepo:  userRepo,
		AlertRepo: alertRepo,
	}
	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	require.Equal(t, "application/json; charset=UTF-8", rr.Header().Get("Content-Type"))

	expectedJsonResp := "[\n" +
		"  {\n" +
		"    \"subject\": \"Test #1\",\n" +
		"    \"status\": \"pending\",\n" +
		"    \"created_at\": \"2018-05-10T08:50:00Z\"\n" +
		"  },\n" +
		"  {\n" +
		"    \"subject\": \"Test #2\",\n" +
		"    \"status\": \"sent\",\n" +
		"    \"sent_at\": \"2018-05-10T08:53:00Z\",\n" +
		"    \"created_at\": \"2018-05-10T08:50:00Z\"\n" +
		"  }\n" +
		"]"

	require.Equal(t, expectedJsonResp, rr.Body.String())
}
