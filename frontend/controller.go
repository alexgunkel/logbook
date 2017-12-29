package frontend

import (
	"github.com/gin-gonic/gin"
	"github.com/satori/go.uuid"
)

const STATIC_APP_DIR_ENV = "STATIC_APP"
const STATIC_RELATIVE_PATH = "/"
const STATIC_BASE_HREF = "/public/"

func AddFrontend(engine *gin.Engine, defaultTemplateDir string) {
	SetApplication(engine, defaultTemplateDir)
	SetStaticApp(engine)
}

func getNewIdentifier() string {
	return uuid.NewV4().String()
}
