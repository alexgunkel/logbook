package application

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/posener/wstest"
	"github.com/stretchr/testify/assert"
)

func GetDispatcher() *gin.Engine {
	engine := gin.Default()
	app := LogBookApplication{}
	app.AddApplicationToEngine(engine)

	return engine
}

// Test the log-message-inbox
func TestEmptyLogEvent(t *testing.T) {
	request, err := http.NewRequest("POST", API_ROOT_PATH+"/1234/logs", strings.NewReader(""))
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
	d := wstest.NewDialer(h, nil) // or t.submit instead of nil

	d.Dial("ws://localhost"+API_ROOT_PATH+"/logbook/12345/logs", nil)

	router := GetDispatcher()
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest("POST", API_ROOT_PATH+"/12345/logs", strings.NewReader("{ \"message\": \"Test\" }"))
	if nil != err {
		t.Fatal(err)
	}

	router.ServeHTTP(recorder, request)
	assert.Equal(t, http.StatusAccepted, recorder.Code)
}

func TestSetAlternativeRootPath(t *testing.T) {
	os.Setenv("API_ROOT_PATH", "/alternative/root/path")
	defer os.Setenv("API_ROOT_PATH", "")

	h := GetDispatcher()
	d := wstest.NewDialer(h, nil)

	d.Dial("ws://localhost/alternative/root/path/logbook/12345/logs", nil)

	router := GetDispatcher()
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "/alternative/root/path/12345/logs", strings.NewReader("{ \"message\": \"Test\" }"))

	router.ServeHTTP(recorder, request)
	assert.Equal(t, http.StatusAccepted, recorder.Code)
}

func TestInValidJsonRefused(t *testing.T) {
	router := GetDispatcher()
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest("POST", API_ROOT_PATH+"/12345/logs", strings.NewReader("{asdasd}"))
	if nil != err {
		t.Fatal(err)
	}

	router.ServeHTTP(recorder, request)
	assert.NotEqual(t, http.StatusAccepted, recorder.Code)
}

// Test the websocket functionality
func TestWebsocketHandlerSwitchesProtocol(t *testing.T) {
	var err error

	h := GetDispatcher()
	d := wstest.NewDialer(h, nil) // or t.submit instead of nil

	c, resp, err := d.Dial("ws://localhost"+API_ROOT_PATH+"/123/logs", nil)
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
	var err error
	engine := gin.Default()
	app := LogBookApplication{}
	app.AddApplicationToEngine(engine)

	if _, ok := app.dispatcher.channels["123"]; ok {
		t.Error("Map entry for 123 should not exist yet.")
	}

	dialer := wstest.NewDialer(engine, nil)
	_, _, err = dialer.Dial("ws://localhost"+API_ROOT_PATH+"/123/logs", nil)
	if err != nil {
		t.Fatal(err)
	}

	if _, ok := app.dispatcher.channels["123"]; !ok {
		t.Error("Map entry for 123 should exist now.")
	}
}

func TestWebsocketRecievesMessagesThatAreSentToTheReceiver(t *testing.T) {
	var err error

	logBook := GetDispatcher()
	dialer := wstest.NewDialer(logBook, nil)
	conn, _, err := dialer.Dial("ws://localhost"+API_ROOT_PATH+"/123/logs", nil)
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	body := "{ \"message\": \"Test\", \"severity\": 4, \"time\": 123123123 }"
	request, err := http.NewRequest("POST", API_ROOT_PATH+"/123/logs", strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	request.Header.Add(LogHeaderLoggerName, "MyLogger")
	request.Header.Add(LogHeaderAppIdentifier, "MyMicroService")
	request.Header.Add(LogHeaderRequestUri, "http://my.web.app")

	logBook.ServeHTTP(recorder, request)

	wsMessage := &LogBookEntry{}
	go conn.ReadJSON(wsMessage)
	time.Sleep(time.Millisecond * 10)

	assert.Equal(t, "Test", wsMessage.Message)
	assert.Equal(t, 123123123, wsMessage.Timestamp)
	assert.Equal(t, 4, wsMessage.Severity)
	assert.Equal(t, "MyMicroService", wsMessage.Application)
	assert.Equal(t, "MyLogger", wsMessage.LoggerName)
	assert.Equal(t, "http://my.web.app", wsMessage.RequestUri)
}

func TestWebsocketReceivesOnlyMessagesWithSameLogBookId(t *testing.T) {
	var err error

	logBook := GetDispatcher()
	dialer := wstest.NewDialer(logBook, nil)
	conn, _, err := dialer.Dial("ws://localhost"+API_ROOT_PATH+"/123/logs", nil)
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	body := "{ \"message\": \"Test\", \"severity\": 4, \"timestamp\": 123123123 }"
	request, err := http.NewRequest("POST", API_ROOT_PATH+"/321/logs", strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}

	logBook.ServeHTTP(recorder, request)

	wsMessage := &IncomingMessage{}
	go conn.ReadJSON(wsMessage)

	time.Sleep(time.Millisecond * 10)

	assert.Equal(t, "", wsMessage.Body.Message)
}
