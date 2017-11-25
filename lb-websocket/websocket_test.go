package lb_websocket

import (
	"github.com/posener/wstest"
	"testing"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/alexgunkel/logbook/entities"
)

func TestWebsocketHandlerEstablishesConnection(t *testing.T) {
	var err error

	channel := make(chan entities.PostMessage, 10)
	c, resp, err := createServer(channel)
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

func TestWebsocketHandlerSendsMessagesWhenReceiving(t *testing.T) {
	var err error
	channel := make(chan entities.PostMessage, 10)
	c, _, err := createServer(channel)
	if err != nil {
		t.Fatal(err)
	}


	input := entities.PostMessage{}
	channel<- input
	message := &entities.PostMessage{}
	err = c.ReadJSON(message)

	if err != nil {
		t.Fatal(err)
	}
}

func createServer(c chan entities.PostMessage) (*websocket.Conn, *http.Response, error) {
	h := gin.Default()
	h.GET("/logbook", func(context *gin.Context) {
		WebsocketHandler(context.Writer, context.Request, c)
	})
	d := wstest.NewDialer(h, nil)
	// or t.Log instead of nil
	conn, resp, err := d.Dial("ws://localhost/logbook", nil)
	return conn, resp, err
}
