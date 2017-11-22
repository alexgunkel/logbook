package logbook

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"encoding/json"
	"strings"
	"github.com/stretchr/testify/assert"
)

// When start page is called without a cookie, then the cookie is set and the client is redirected
// to the log-list page
func TestInitLogBookWithoutCookie(t *testing.T) {
	request, err := http.NewRequest("GET", "/logbook", nil)
	if nil != err {
		t.Fatal(err)
	}

	router := Application()
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)


	// Copy the Cookie over to a new Request
	newRequest := &http.Request{Header: http.Header{"Cookie": recorder.HeaderMap["Set-Cookie"]}}

	// Extract the dropped cookie from the request.
	clientIdCookie, _ := newRequest.Cookie("logbook")
	clientId := clientIdCookie.Value

	assert.Equal(t, http.StatusTemporaryRedirect, recorder.Code)
	assert.NotEqual(t, "", clientId)
	path := "/" + clientId + "/logs"
	assert.Equalf(t, path, recorder.Header().Get("Location"), "Expected path %v", path)
}

// GET request to a specific display path without cookie results in a redirect to the start page
func TestDisplayWithoutCookie(t *testing.T)  {
	router := Application()
	recorder := httptest.NewRecorder()
	request,_ := http.NewRequest("GET", "/logbook/1234/logs", nil)

	router.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusTemporaryRedirect, recorder.Code)
	assert.Equal(t, "/logbook", recorder.Header().Get("Location"))
}

func TestEmptyLogEventLogEvent(t *testing.T) {
	request, err := http.NewRequest("POST", "/logbook/1234/logs", nil)
	if nil != err {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	router := Application()
	router.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestValidLogAccepted(t *testing.T)  {
	router := Application()
	recorder := httptest.NewRecorder()
	requestBody,_ := json.Marshal(logEvent{"Test"})
	request, err := http.NewRequest("POST", "/logbook/12345/logs", strings.NewReader(string(requestBody)))
	if nil != err {
		t.Fatal(err)
	}

	router.ServeHTTP(recorder, request)
	assert.Equal(t, http.StatusOK, recorder.Code)
}
