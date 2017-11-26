package lb_logbook

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"strings"
	"github.com/stretchr/testify/assert"
	"github.com/posener/wstest"
	"github.com/alexgunkel/logbook/lb-entities"
	"github.com/alexgunkel/logbook/lb-receiver"
)



// A normal session starts by a HTTP GET-request at <domain>/logbook. We assume that no cookie is
// set. Therefore, we generate a client-id and set a cookie.
//
// There is no need to redirect. Instead we send a small application that shows the success.
func TestInitLogBookWithoutCookie(t *testing.T) {
	request, err := http.NewRequest("GET", "/logbook", nil)
	if nil != err {
		t.Fatal(err)
	}

	router := Application("../resources/private/template")
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	logBookId, _ := getRecorderCookie(recorder)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.NotEqual(t, "", logBookId)
}

// When receiving a GET-request to the root page and the client has logbook-cookie, then we
// directly show the app.
func TestInitLogBookWithCookie(t *testing.T) {
	router := Application("../resources/private/template")
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/logbook", nil)
	cookie := &http.Cookie{Name: "logbook", Value: "1234"}
	request.AddCookie(cookie)
	router.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusOK, recorder.Code)
}

// Test the log-message-receiver
func TestEmptyLogEvent(t *testing.T) {
	request, err := http.NewRequest("POST", "/logbook/1234/logs", strings.NewReader(""))
	if nil != err {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	router := Application("")
	router.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestValidLogAccepted(t *testing.T) {
	router := Application("")
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest("POST", "/logbook/12345/logs", strings.NewReader("{ \"message\": \"Test\" }"))
	if nil != err {
		t.Fatal(err)
	}

	router.ServeHTTP(recorder, request)
	assert.Equal(t, http.StatusOK, recorder.Code)
}

func TestInValidJsonRefused(t *testing.T) {
	router := Application("")
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

	h := Application("")
	d := wstest.NewDialer(h, nil)  // or t.Log instead of nil

	c, resp, err := d.Dial("ws://localhost/logbook/123/logs", nil)
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

	logBook := Application("")
	dialer := wstest.NewDialer(logBook, nil)
	conn, _, err := dialer.Dial("ws://localhost/logbook/123/logs", nil)
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	body := "{ \"message\": \"Test\", \"severity\": 4, \"timestamp\": 123123123 }"
	request, err := http.NewRequest("POST", "/logbook/12345/logs", strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	request.Header.Add(lb_receiver.LogHeaderLoggerName, "MyLogger")
	request.Header.Add(lb_receiver.LogHeaderAppIdentifier, "MyMicroService")
	request.Header.Add(lb_receiver.LogHeaderRequestUri, "http://my.web.app")

	logBook.ServeHTTP(recorder, request)

	wsMessage := &lb_entities.PostMessage{}
	conn.ReadJSON(wsMessage)

	assert.Equal(t, "Test", wsMessage.Event.Message)
	assert.Equal(t, 123123123, wsMessage.Event.Timestamp)
	assert.Equal(t, 4, wsMessage.Event.Severity)
	assert.Equal(t, "MyMicroService", wsMessage.Header.Application)
	assert.Equal(t, "MyLogger", wsMessage.Header.LoggerName)
	assert.Equal(t, "http://my.web.app", wsMessage.Header.RequestUri)
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
