package application

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"strings"
	"github.com/stretchr/testify/assert"
	"github.com/posener/wstest"
	"github.com/gin-gonic/gin"
	"time"
)

func GetDispatcher() *gin.Engine {
	engine := gin.Default()
	app := LogBookApplication{}
	app.AddApplicationToEngine(engine)

	return engine
}

// Test the log-message-inbox
func TestEmptyLogEvent(t *testing.T) {
	request, err := http.NewRequest("POST", "/logbook/1234/logs", strings.NewReader(""))
	if nil != err {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	router := GetDispatcher()
	router.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestValidLogAccepted(t *testing.T) {
	var err error

	h := GetDispatcher()
	d := wstest.NewDialer(h, nil)  // or t.submit instead of nil

	d.Dial("ws://localhost/logbook/12345/logs", nil)

	router := GetDispatcher()
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest("POST", "/logbook/12345/logs", strings.NewReader("{ \"message\": \"Test\" }"))
	if nil != err {
		t.Fatal(err)
	}

	router.ServeHTTP(recorder, request)
	assert.Equal(t, http.StatusOK, recorder.Code)
}

func TestInValidJsonRefused(t *testing.T) {
	router := GetDispatcher()
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

	h := GetDispatcher()
	d := wstest.NewDialer(h, nil)  // or t.submit instead of nil

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

func TestWebsocketHandlerCreatesChannel(t *testing.T) {
	var err  error
	engine := gin.Default()
	app := LogBookApplication{}
	app.AddApplicationToEngine(engine)

	if _, ok := app.d.channels["123"]; ok {
		t.Error("Map entry for 123 should not exist yet.")
	}

	dialer := wstest.NewDialer(engine, nil)
	_, _, err = dialer.Dial("ws://localhost/logbook/123/logs", nil)
	if err != nil {
		t.Fatal(err)
	}

	if _, ok := app.d.channels["123"]; !ok {
		t.Error("Map entry for 123 should exist now.")
	}
}

func TestWebsocketRecievesMessagesThatAreSentToTheReceiver(t *testing.T)  {
	var err  error

	logBook := GetDispatcher()
	dialer := wstest.NewDialer(logBook, nil)
	conn, _, err := dialer.Dial("ws://localhost/logbook/123/logs", nil)
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	body := "{ \"message\": \"Test\", \"severity\": 4, \"timestamp\": 123123123 }"
	request, err := http.NewRequest("POST", "/logbook/123/logs", strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	request.Header.Add(LogHeaderLoggerName, "MyLogger")
	request.Header.Add(LogHeaderAppIdentifier, "MyMicroService")
	request.Header.Add(LogHeaderRequestUri, "http://my.web.app")

	logBook.ServeHTTP(recorder, request)

	wsMessage := &Message{}
	go conn.ReadJSON(wsMessage)
	time.Sleep(time.Millisecond * 10)

	assert.Equal(t, "Test", wsMessage.Event.Message)
	assert.Equal(t, 123123123, wsMessage.Event.Timestamp)
	assert.Equal(t, 4, wsMessage.Event.Severity)
	assert.Equal(t, "MyMicroService", wsMessage.Origin.Application)
	assert.Equal(t, "MyLogger", wsMessage.Origin.LoggerName)
	assert.Equal(t, "http://my.web.app", wsMessage.Origin.RequestUri)
}

func TestWebsocketReceivesOnlyMessagesWithSameLogBookId(t *testing.T) {
	var err  error

	logBook := GetDispatcher()
	dialer := wstest.NewDialer(logBook, nil)
	conn, _, err := dialer.Dial("ws://localhost/logbook/123/logs", nil)
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	body := "{ \"message\": \"Test\", \"severity\": 4, \"timestamp\": 123123123 }"
	request, err := http.NewRequest("POST", "/logbook/321/logs", strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}

	logBook.ServeHTTP(recorder, request)

	wsMessage := &Message{}
	go conn.ReadJSON(wsMessage)

	time.Sleep(time.Millisecond * 10)

	assert.Equal(t, "", wsMessage.Event.Message)
}
