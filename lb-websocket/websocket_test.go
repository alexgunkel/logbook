package lb_websocket

import (
	"github.com/posener/wstest"
	"testing"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func TestHandler(t *testing.T) {
	var err error

	c, resp, err := createServer()
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

func createServer() (*websocket.Conn, *http.Response, error) {
	h := gin.Default()
	h.GET("/logbook", func(context *gin.Context) {
		WebsocketHandler(context.Writer, context.Request)
	})
	d := wstest.NewDialer(h, nil)
	// or t.Log instead of nil
	c, resp, err := d.Dial("ws://localhost/logbook", nil)
	return c, resp, err
}
