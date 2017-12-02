package frontend

import (
	"net/http/httptest"
	"github.com/stretchr/testify/assert"
	"github.com/alexgunkel/logbook/application"
	"testing"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"io/ioutil"
)

// if STATIC_APP is not set, a request to /logbook should
// be handled by the default frontend-application
func TestLogBookRequestWithourEnvSet(t *testing.T) {
	router := gin.Default()
	AddFrontend(router, "../public")

	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/logbook", nil)
	router.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "<body>")
	assert.Contains(t, recorder.Body.String(), "LogBook")
	assert.Contains(t, recorder.Body.String(), "ws://localhost:8080" + application.API_ROOT_PATH)
}


// if STATIC APP is set, a request to /logbook should be handled
// by the respective static files
func TestServeTemplateWithStaticAppConfigured(t *testing.T) {
	tmp, _ := ioutil.TempDir("", "")
	content := "html"
	ioutil.WriteFile(tmp + "/index.html", []byte(content), os.ModePerm)

	os.Setenv(STATIC_APP_DIR_ENV, tmp)

	engine := gin.Default()
	AddFrontend(engine, tmp)
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", STATIC_RELATIVE_PATH, nil)
	engine.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, content, recorder.Body.String())
}

// in either case, static files should be returned from the
// given static-directory
func TestServeStaticFilesWithStaticAppConfigured(t *testing.T) {
	tmp, _ := ioutil.TempDir("", "")
	js := "test javascript content"
	ioutil.WriteFile(tmp + "/test.js", []byte(js), os.ModePerm)

	os.Setenv(STATIC_APP_DIR_ENV, tmp)

	engine := gin.Default()
	AddFrontend(engine, tmp)
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", STATIC_RELATIVE_PATH + "/test.js", nil)
	engine.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, js, recorder.Body.String())
}

// Test handling of env variables:
// HOST
// PORT
func TestEnvVariables(t *testing.T) {
	os.Setenv("PORT", "1234")
	defer os.Setenv("PORT", "")

	os.Setenv("HOST", "123.456.789.132")
	defer os.Setenv("HOST", "")

	router := gin.Default()
	AddFrontend(router, "../public")

	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/logbook", nil)
	router.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "<body>")
	assert.Contains(t, recorder.Body.String(), "LogBook")
	assert.Contains(t, recorder.Body.String(), "ws://123.456.789.132:1234" + application.API_ROOT_PATH)
}