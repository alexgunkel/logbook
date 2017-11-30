package frontend

import (
	"github.com/gin-gonic/gin"
	"os"
)

const STATIC_APP_DIR_ENV  = "STATIC_APP"

func SetStaticApp(engine *gin.Engine)  {
	path := os.Getenv(STATIC_APP_DIR_ENV) + "/index.html"
	engine.StaticFile("/logbook", path)
}
