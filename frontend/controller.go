package frontend

import (
	"strconv"
	"github.com/gin-gonic/gin"
)

func AddFrontend(engine *gin.Engine, templateDir string) {
	SetApplication(engine, templateDir)
	SetStaticApp(engine)
}

func SetApplication(engine *gin.Engine, templateDir string)  {
	fe := &WebApplication{}
	fe.SetTemplateDirPath(templateDir)
	gen := &IdGenerator{}

	engine.GET("/logbook", func(context *gin.Context) {
		fe.InitLogBookClientApplication(context, gen)
	})
}

// Generate new IDs to serve websocket requests
type IdGenerator struct {
	lastIdentifier int64
}

func (app *IdGenerator) getNewIdentifier() string {
	app.lastIdentifier++
	return strconv.FormatInt(app.lastIdentifier, 10)
}
