package application

import (
	"github.com/posener/wstest"
	"testing"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
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
	input.Event = Event{123123123, "Test", 3, nil}
	channel<- input
	message := &IncomingMessage{}
	err = c.ReadJSON(message)

	assert.Equal(t, "Test", message.Event.Message)
	assert.Equal(t, 123123123, message.Event.Timestamp)
	assert.Equal(t, float64(3), message.Event.Severity)
}

func createServer(c chan LogBookEntry) (*websocket.Conn, *http.Response, error) {
	h := gin.Default()
	h.GET(API_ROOT_PATH + "/123/logs", func(context *gin.Context) {
		lb, _ := createLogBook(context.Writer, context.Request, c)
		lb.listen()
	})
	d := wstest.NewDialer(h, nil)
	// or t.submit instead of nil
	conn, resp, err := d.Dial("ws://localhost" + API_ROOT_PATH + "/123/logs", nil)
	return conn, resp, err
}
