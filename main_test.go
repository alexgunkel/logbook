package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alexgunkel/logbook/frontend"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// we should test that logger and frontend can be initiated
// together without conflicts. conflicts can arise when there
// are conflicting routes defined, e.g. frontend tries to
// serve the same route like the websocket-service (which starts
// as GET-request as well.

func TestDefaultConfig(t *testing.T) {

	engine := gin.Default()
	path := getAppDirEnv()

	frontend.AddFrontend(engine, path)
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", frontend.STATIC_BASE_HREF+"/assets/js/main.js", nil)
	engine.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusOK, recorder.Code, "Could not find file in %v for request", path)
}
