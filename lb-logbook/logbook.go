package lb_logbook

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

func Default(templateDir string) *gin.Engine {
	app := new(Webapp)
	app.engine = gin.Default()
	frontend := &lb_frontend.WebApplication{}
	frontend.SetTemplateDirPath(templateDir)
	gen := &lb_frontend.IdGenerator{}
	incoming := make(chan lb_entities.PostMessage, 20)
	//outbound := make(chan lb_entities.PostMessage, 10)

	app.engine.GET("/logbook", func(context *gin.Context) {
		frontend.InitLogBookClientApplication(context, gen)
	})

	app.engine.POST("/logbook/:client/logs", func(context *gin.Context) {
		lb_receiver.Log(context, incoming)
	})

	app.engine.GET("logbook/:client/logs", func(context *gin.Context) {
		lb_websocket.ForceCookie(context)
		lb_websocket.WebsocketHandler(context.Writer, context.Request, incoming)
	})

	return app.engine
}

// @API
func Application(templateDir string) *gin.Engine {
	return Default(templateDir)
}

