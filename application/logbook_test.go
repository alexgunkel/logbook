package application

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/posener/wstest"
	"github.com/stretchr/testify/assert"
)

func TestWebsocketHandlerEstablishesConnection(t *testing.T) {
	var err error

	channel := make(chan LogBookEntry, 10)
	_, resp, err := createServer(channel)
	if err != nil {
		t.Fatal(err)
	}

	if got, want := resp.StatusCode, http.StatusSwitchingProtocols; got != want {
		t.Errorf("resp.StatusCode = %q, want %q", got, want)
	}
}

func TestWebsocketHandlerSendsMessagesWhenReceiving(t *testing.T) {
	var err error
	channel := make(chan LogBookEntry, 10)
	c, _, err := createServer(channel)
	if err != nil {
		t.Fatal(err)
	}

	input := LogBookEntry{}
	input.Timestamp = 123123123
	input.Message = "Test"
	input.Severity = 3
	channel <- input
	message := &LogBookEntry{}
	err = c.ReadJSON(message)

	assert.Equal(t, "Test", message.Message)
	assert.Equal(t, 123123123, message.Timestamp)
	assert.Equal(t, 3, message.Severity)
}

func createServer(c chan LogBookEntry) (*websocket.Conn, *http.Response, error) {
	h := gin.Default()
	h.GET(API_ROOT_PATH+"/123/logs", func(context *gin.Context) {
		lb, _ := createLogBook(context.Writer, context.Request, c)
		lb.listen()
	})
	d := wstest.NewDialer(h, nil)
	// or t.submit instead of nil
	conn, resp, err := d.Dial("ws://localhost"+API_ROOT_PATH+"/123/logs", nil)
	return conn, resp, err
}
