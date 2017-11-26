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

func WebsocketHandler(w http.ResponseWriter, r *http.Request, c <-chan PostMessage) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		panic(err)
		return
	}

	for {
		msg :=<- c

		conn.WriteJSON(msg)
	}
}

func ForceCookie(c *gin.Context)  {
	_, err := c.Cookie("logbook")
	if nil != err {
		identifier := c.Param("client")
		c.SetCookie("logbook", identifier, 0, "", "", false, false)
	}
}
