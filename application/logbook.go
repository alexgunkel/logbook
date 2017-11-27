package application

import (
	"github.com/gorilla/websocket"
	"net/http"
	"github.com/gin-gonic/gin"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type logBook struct {
	conn    *websocket.Conn
	mailbox <-chan Message
}

func createLogBook(w http.ResponseWriter, r *http.Request, c <-chan Message) (lb logBook, err error) {
	lb.conn, err = upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	lb.mailbox = c

	return
}

func (lb *logBook) listen() {
	for {
		msg := <-lb.mailbox
		lb.conn.WriteJSON(msg)
	}
}

func ForceCookie(c *gin.Context, id string) {
	_, err := c.Cookie("logbook")
	if nil != err {
		c.SetCookie("logbook", id, 0, "", "", false, false)
	}
}
