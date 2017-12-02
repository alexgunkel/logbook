package frontend

import (
	"os"

	"github.com/gin-gonic/gin"
)

const STATIC_APP_DIR_ENV = "STATIC_APP"
const STATIC_RELATIVE_PATH = "/logbook/app/"

func SetStaticApp(engine *gin.Engine) {
	path := os.Getenv(STATIC_APP_DIR_ENV)
	engine.Static(STATIC_RELATIVE_PATH+"public/", path)
}
