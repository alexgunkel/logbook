package logbook

import (
	"github.com/alexgunkel/logbook/services"
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/alexgunkel/logbook/entities"
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
	gen := &services.IdGenerator{}
	incoming := make(chan entities.LogEvent, 20)

	app.engine.GET("/logbook", func(context *gin.Context) {
		services.InitLogBookClientApplication(context, gen)
	})

	app.engine.GET("/logbook/:client/logs", func(context *gin.Context) {
		services.DisplayLogs(context, gen)
	})

	app.engine.POST("/logbook/:client/logs", func(context *gin.Context) {
		services.Log(context, incoming)
	})

	return app.engine
}

// @API
func Application() *gin.Engine {
	return Default()
}

