package frontend

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func AddFrontend(engine *gin.Engine, defaultTemplateDir string) {
	SetApplication(engine, defaultTemplateDir)
	SetStaticApp(engine)
}

// Generate new IDs to serve websocket requests
type IdGenerator struct {
	lastIdentifier int64
}

func (app *IdGenerator) getNewIdentifier() string {
	app.lastIdentifier++
	return strconv.FormatInt(app.lastIdentifier, 10)
}
