package logbook

import (
	"github.com/alexgunkel/logbook/services"
	"github.com/gin-gonic/gin"
)

// @API
func Application() *gin.Engine {
	return services.Default()
}

