package router

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPricePlanRouteSuccess(t *testing.T) {
	rr := httptest.NewRecorder()

	handler := RegisterRoutes()

	resp, err := http.NewRequest(http.MethodPost, "/client", strings.NewReader(""))

	handler.ServeHTTP(rr, resp)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rr.Code)

	expectedContentType := "application/json"
	actualContentType := rr.Header().Get("Content-Type")
	assert.Equal(t, expectedContentType, actualContentType)
	
	expectedResponse := "Hello"
	actualResponse := rr.Body.String()
	assert.Equal(t, expectedResponse, actualResponse)
}
