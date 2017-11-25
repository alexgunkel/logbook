package lb_websocket

import (
	"github.com/gorilla/websocket"
	"net/http"
	"github.com/alexgunkel/logbook/entities"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func WebsocketHandler(w http.ResponseWriter, r *http.Request, c chan entities.PostMessage) {
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
