package frontend

import (
	"os"

	"github.com/gin-gonic/gin"
)

func SetStaticApp(engine *gin.Engine) {
	path := os.Getenv(STATIC_APP_DIR_ENV)
	engine.Static(STATIC_BASE_HREF, path)
}
