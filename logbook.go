package logbook

import (
	"github.com/alexgunkel/logbook/services"
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/alexgunkel/logbook/entities"
	"github.com/gorilla/websocket"
)


type Webapp struct {
	engine *gin.Engine
}

func (app *Webapp) ServeHTTP(w http.ResponseWriter, req *http.Request)  {
	app.engine.ServeHTTP(w, req)
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func handler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		panic(err)
		return
	}

	for {
		t, msg, err := conn.ReadMessage()
		if err != nil {
			break
		}
		conn.WriteMessage(t, msg)
	}
}

func Default() *gin.Engine {
	app := new(Webapp)
	app.engine = gin.Default()
	gen := &services.IdGenerator{}
	incoming := make(chan entities.PostMessage, 20)

	app.engine.GET("/logbook", func(context *gin.Context) {
		services.InitLogBookClientApplication(context, gen)
	})

	app.engine.GET("/logbook/:client/logs", func(context *gin.Context) {
		services.DisplayLogs(context, gen)
	})

	app.engine.POST("/logbook/:client/logs", func(context *gin.Context) {
		services.Log(context, incoming)
	})

	app.engine.GET("logbook/:client/ws", func(context *gin.Context) {
		handler(context.Writer, context.Request)
	})

	return app.engine
}

// @API
func Application() *gin.Engine {
	return Default()
}

