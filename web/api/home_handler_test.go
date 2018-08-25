package api_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bakku/easyalert/web/api"

	"github.com/stretchr/testify/require"
)

func TestHomeHandler_Success(t *testing.T) {
	req, err := http.NewRequest("GET", "/api", nil)
	require.Nil(t, err)

	rr := httptest.NewRecorder()
	handler := api.HomeHandler{}

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	require.Equal(t, "application/json; charset=UTF-8", rr.Header().Get("Content-Type"))
	require.Equal(t, `{"easyalert":"Alerting made easy"}`, rr.Body.String())
}
