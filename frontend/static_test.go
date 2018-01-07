package frontend

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestServeStaticJs(t *testing.T) {
	tmp, _ := ioutil.TempDir("", "")
	content := "test"
	ioutil.WriteFile(tmp+"/app.js", []byte(content), os.ModePerm)

	os.Setenv(STATIC_APP_DIR_ENV, tmp)
	defer os.Setenv(STATIC_APP_DIR_ENV, "")

	engine := gin.Default()
	SetStaticApp(engine)
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", STATIC_BASE_HREF+"/app.js", nil)
	engine.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusOK, recorder.Code, "Asserted that %v finds %v", STATIC_BASE_HREF+"/app.js", tmp+"/app.js")
	assert.Equal(t, content, recorder.Body.String())
}

func TestServeStaticJsWithPath(t *testing.T) {
	tmp, _ := ioutil.TempDir("", "")
	content := "test"
	os.Mkdir(tmp+"/path", os.ModePerm)
	err := ioutil.WriteFile(tmp+"/path/app.js", []byte(content), os.ModePerm)
	if nil != err {
		t.Fatal(err)
	}
	os.Setenv(STATIC_APP_DIR_ENV, tmp)
	defer os.Setenv(STATIC_APP_DIR_ENV, "")

	engine := gin.Default()
	SetStaticApp(engine)
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", STATIC_BASE_HREF+"/path/app.js", nil)
	engine.ServeHTTP(recorder, request)

	assert.FileExists(t, tmp+"/path/app.js")
	assert.Equal(t, http.StatusOK, recorder.Code, "Asserted that %v finds %v", STATIC_BASE_HREF+"/path/app.js", tmp+"/path/app.js")
	assert.Equal(t, content, recorder.Body.String())
}
