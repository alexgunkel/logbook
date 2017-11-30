package frontend

import (
	"testing"
	"os"
	"io/ioutil"
	"github.com/gin-gonic/gin"
	"net/http/httptest"
	"net/http"
	"github.com/stretchr/testify/assert"
)

func TestServeStaticFiles(t *testing.T) {
	tmp, _ := ioutil.TempDir("", "")
	content := "html"
	ioutil.WriteFile(tmp + "/index.html", []byte(content), os.ModePerm)

	os.Setenv(STATIC_APP_DIR_ENV, tmp)

	engine := gin.Default()
	SetStaticApp(engine)
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/logbook", nil)
	engine.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, content, recorder.Body.String())
}

func TestServeStaticFilesWithEndSlash(t *testing.T) {
	tmp, _ := ioutil.TempDir("", "")
	content := "html"
	ioutil.WriteFile(tmp + "/index.html", []byte(content), os.ModePerm)

	os.Setenv(STATIC_APP_DIR_ENV, tmp + "/")

	engine := gin.Default()
	SetStaticApp(engine)
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/logbook", nil)
	engine.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, content, recorder.Body.String())
}

func TestServeStaticJs(t *testing.T) {
	tmp, _ := ioutil.TempDir("", "")
	content := "test"
	ioutil.WriteFile(tmp + "/app.js", []byte(content), os.ModePerm)

	os.Setenv(STATIC_APP_DIR_ENV, tmp)

	engine := gin.Default()
	SetStaticApp(engine)
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/logbook/app.js", nil)
	engine.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, content, recorder.Body.String())
}