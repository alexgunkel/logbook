package logbook

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"strings"
	"github.com/stretchr/testify/assert"
	"github.com/posener/wstest"
	"github.com/alexgunkel/logbook/lb-entities"
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

	logBookId, _ := getRecorderCookie(recorder)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusTemporaryRedirect, recorder.Code)
	assert.NotEqual(t, "", logBookId)
	path := "/logbook/" + logBookId + "/logs"
	assert.Equalf(t, path, recorder.Header().Get("Location"), "Expected path %v", path)
}

// When receiving a GET-request to the root page and the client has logbook-cookie, then we redirect her to the
// logs without setting a new cookie
func TestInitLogBookWithCookie(t *testing.T) {
	router := Application()
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/logbook", nil)
	cookie := &http.Cookie{Name: "logbook", Value: "1234"}
	request.AddCookie(cookie)
	router.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusTemporaryRedirect, recorder.Code)
	assert.Equal(t, "/logbook/1234/logs", recorder.Header().Get("Location"))
}

// GET request to a specific display path without cookie results in a redirect to the start page
func TestDisplayWithoutCookie(t *testing.T) {
	router := Application()
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/logbook/1234/logs", nil)

	router.ServeHTTP(recorder, request)
	logBookId, err := getRecorderCookie(recorder)
	newPath := "/logbook/" + logBookId + "/logs"

	assert.Nil(t, err)
	assert.Equal(t, http.StatusTemporaryRedirect, recorder.Code)
	assert.Equalf(t, newPath, recorder.Header().Get("Location"), "Expected path %v", newPath)
}

// Test the log-message-receiver
func TestEmptyLogEvent(t *testing.T) {
	request, err := http.NewRequest("POST", "/logbook/1234/logs", strings.NewReader(""))
	if nil != err {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	router := Application()
	router.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestValidLogAccepted(t *testing.T) {
	router := Application()
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest("POST", "/logbook/12345/logs", strings.NewReader("{ \"message\": \"Test\" }"))
	if nil != err {
		t.Fatal(err)
	}

	router.ServeHTTP(recorder, request)
	assert.Equal(t, http.StatusOK, recorder.Code)
}

func TestInValidJsonRefused(t *testing.T) {
	router := Application()
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest("POST", "/logbook/12345/logs", strings.NewReader("{asdasd}"))
	if nil != err {
		t.Fatal(err)
	}

	router.ServeHTTP(recorder, request)
	assert.NotEqual(t, http.StatusOK, recorder.Code)
}

// Test the websocket functionality
func TestWebsocketHandlerSwitchesProtocol(t *testing.T) {
	var err error

	h := Application()
	d := wstest.NewDialer(h, nil)  // or t.Log instead of nil

	c, resp, err := d.Dial("ws://localhost/logbook/123/ws", nil)
	if err != nil {
		t.Fatal(err)
	}

	if got, want := resp.StatusCode, http.StatusSwitchingProtocols; got != want {
		t.Errorf("resp.StatusCode = %q, want %q", got, want)
	}

	err = c.WriteJSON("test")
	if err != nil {
		t.Fatal(err)
	}
}

func TestWebsocketRecievesMessagesThatAreSentToTheReceiver(t *testing.T)  {
	var err  error

	logBook := Application()
	dialer := wstest.NewDialer(logBook, nil)
	conn, _, err := dialer.Dial("ws://localhost/logbook/123/ws", nil)
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	request, err := http.NewRequest("POST", "/logbook/12345/logs", strings.NewReader("{ \"message\": \"Test\" }"))
	if err != nil {
		t.Fatal(err)
	}

	logBook.ServeHTTP(recorder, request)

	wsMessage := &lb_entities.PostMessage{}
	conn.ReadJSON(wsMessage)

	assert.Equal(t, "Test", wsMessage.Event.Message)
}

// Helper function to get cookie values out of response recorders
func getRecorderCookie(r *httptest.ResponseRecorder) (clientId string, err error) {
	newRequest := &http.Request{Header: http.Header{"Cookie": r.HeaderMap["Set-Cookie"]}}
	clientIdCookie, err := newRequest.Cookie("logbook")

	if nil == err {
		clientId = clientIdCookie.Value
	} else {
		clientId = ""
	}

	return clientId, err
}
