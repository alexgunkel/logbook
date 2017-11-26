package logbook

import (
	"github.com/alexgunkel/logbook/lb-receiver"
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/alexgunkel/logbook/lb-entities"
	"github.com/alexgunkel/logbook/lb-websocket"
	"github.com/alexgunkel/logbook/lb-frontend"
)


type Webapp struct {
	engine *gin.Engine
}

func (app *Webapp) ServeHTTP(w http.ResponseWriter, req *http.Request)  {
	app.engine.ServeHTTP(w, req)
}

func Default() *gin.Engine {
	app := new(Webapp)
	app.engine = gin.Default()
	gen := &lb_frontend.IdGenerator{}
	incoming := make(chan lb_entities.PostMessage, 20)
	outbound := make(chan lb_entities.PostMessage, 10)

	app.engine.GET("/logbook", func(context *gin.Context) {
		lb_frontend.InitLogBookClientApplication(context, gen)
	})

	app.engine.GET("/logbook/:client/logs", func(context *gin.Context) {
		lb_frontend.DisplayLogs(context, gen)
	})

	app.engine.POST("/logbook/:client/logs", func(context *gin.Context) {
		lb_receiver.Log(context, incoming)
	})

	app.engine.GET("logbook/:client/ws", func(context *gin.Context) {
		lb_websocket.WebsocketHandler(context.Writer, context.Request, outbound)
	})

	return app.engine
}

// @API
func Application() *gin.Engine {
	return Default()
}

