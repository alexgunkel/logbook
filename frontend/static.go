package frontend

import (
	"github.com/gin-gonic/gin"
	"os"
)

const STATIC_APP_DIR_ENV  = "STATIC_APP"
const STATIC_RELATIVE_PATH = "/logbook/app"

func SetStaticApp(engine *gin.Engine)  {
	path := os.Getenv(STATIC_APP_DIR_ENV)
	engine.StaticFile(STATIC_RELATIVE_PATH, path + "/index.html")
	engine.Static(STATIC_RELATIVE_PATH, path)
}
