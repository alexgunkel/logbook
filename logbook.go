package logbook

import (
	"github.com/alexgunkel/logbook/services"
)

const CONTENT_TYPE_JSON  = "application/json"
const CONTENT_TYPE_HTML  = "text/html"

// @API
func Application() *services.Webapp {
	return services.Default()
}

