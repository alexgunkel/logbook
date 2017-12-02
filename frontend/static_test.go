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

	engine := gin.Default()
	SetStaticApp(engine)
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", STATIC_BASE_HREF+"app.js", nil)
	engine.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, content, recorder.Body.String())
}
