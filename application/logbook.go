package application

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/alexgunkel/logbook/frontend"
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
	fe := &frontend.WebApplication{}
	fe.SetTemplateDirPath(templateDir)
	gen := &frontend.IdGenerator{}
	incoming := make(chan PostMessage, 20)
	//outbound := make(chan lb_entities.PostMessage, 10)

	app.engine.GET("/logbook", func(context *gin.Context) {
		fe.InitLogBookClientApplication(context, gen)
	})

	app.engine.POST("/logbook/:client/logs", func(context *gin.Context) {
		Log(context, incoming)
	})

	app.engine.GET("logbook/:client/logs", func(context *gin.Context) {
		ForceCookie(context)
		WebsocketHandler(context.Writer, context.Request, incoming)
	})

	return app.engine
}

// @API
func Application(templateDir string) *gin.Engine {
	return Default(templateDir)
}

