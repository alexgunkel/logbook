package logbook

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"encoding/json"
	"strings"
	"github.com/stretchr/testify/assert"
)

// A normal session starts by a HTTP GET-request at <domain>/logbook. We assume that no cookie is
// set. Therefore, we generate a client-id, set a cookie, and redirect to the log-message-list-page.
//
func TestInitLogBookWithoutCookie(t *testing.T) {
	request, err := http.NewRequest("GET", "/logbook", nil)
	if nil != err {
		t.Fatal(err)
	}

	router := Application()
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)


	clientId := getRecorderCookie(recorder).Value

	assert.Equal(t, http.StatusTemporaryRedirect, recorder.Code)
	assert.NotEqual(t, "", clientId)
	path := "/logbook/" + clientId + "/logs"
	assert.Equalf(t, path, recorder.Header().Get("Location"), "Expected path %v", path)
}

// When receiving a GET-request to the root page and the client has logbook-cookie, then we redirect her to the
// logs without setting a new cookie
func TestInitLogBookWithCookie(t *testing.T) {
	router := Application()
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/logbook", nil)
	cookie := &http.Cookie{Name: "logbook", Value:"1234"}
	request.AddCookie(cookie)
	router.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusTemporaryRedirect, recorder.Code)
	assert.Equal(t, "/logbook/1234/logs", recorder.Header().Get("Location"))
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

func getRecorderCookie(r *httptest.ResponseRecorder) *http.Cookie {
	newRequest := &http.Request{Header: http.Header{"Cookie": r.HeaderMap["Set-Cookie"]}}
	clientIdCookie, err := newRequest.Cookie("logbook")

	if nil != err {
		panic(err)
	}

	return clientIdCookie
}
