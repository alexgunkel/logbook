package frontend

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/alexgunkel/logbook/application"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// if STATIC_APP is not set, a request to /logbook should
// be handled by the default frontend-application
func TestLogBookRequestWithoutEnvSet(t *testing.T) {
	router := gin.Default()
	AddFrontend(router, "../public")

	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", STATIC_RELATIVE_PATH, nil)
	router.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "<body>")
	assert.Contains(t, recorder.Body.String(), "LogBook")
	assert.Contains(t, recorder.Body.String(), "ws://localhost:8080"+application.API_ROOT_PATH)
}

// if STATIC APP is set, a request to /logbook should be handled
// by the respective static files
func TestServeTemplateWithStaticAppConfigured(t *testing.T) {
	tmp, _ := ioutil.TempDir("", "")
	os.Setenv(STATIC_APP_DIR_ENV, tmp)

	type contentObj struct {
		content string
		name    string
		path    string
		result  string
	}
	var contents []contentObj
	contents = append(contents, contentObj{"html", "Index.html", "", "html"})
	contents = append(contents, contentObj{"{{.BaseHref}}", "Index.html", "", STATIC_BASE_HREF})
	contents = append(contents, contentObj{"test-js", "test.js", "public/test.js", "test-js"})

	for _, content := range contents {
		t.Run(content.name, func(t *testing.T) {
			ioutil.WriteFile(tmp+"/"+content.name, []byte(content.content), os.ModePerm)

			engine := gin.Default()
			AddFrontend(engine, tmp)
			recorder := httptest.NewRecorder()
			request, _ := http.NewRequest("GET", STATIC_RELATIVE_PATH+content.path, nil)
			engine.ServeHTTP(recorder, request)

			assert.Equal(t, http.StatusOK, recorder.Code)
			assert.Equal(t, content.result, recorder.Body.String())

		})
	}
}

// in either case, static files should be returned from the
// given static-directory
func TestServeStaticFilesWithStaticAppConfigured(t *testing.T) {
	tmp, _ := ioutil.TempDir("", "")
	js := "test javascript content"
	ioutil.WriteFile(tmp+"/test.js", []byte(js), os.ModePerm)

	os.Setenv(STATIC_APP_DIR_ENV, tmp)

	engine := gin.Default()
	AddFrontend(engine, tmp)
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", STATIC_BASE_HREF+"test.js", nil)
	engine.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, js, recorder.Body.String())
}

// Test handling of env variables:
// HOST
// PORT
func TestEnvVariables(t *testing.T) {
	eVars := provideEnvVariables()
	for _, eVar := range eVars {
		t.Run(eVar.host+":"+eVar.port, func(t *testing.T) {
			os.Setenv("PORT", eVar.port)
			defer os.Setenv("PORT", "")

			os.Setenv("HOST", eVar.host)
			defer os.Setenv("HOST", "")

			router := gin.Default()
			AddFrontend(router, "../public")

			recorder := httptest.NewRecorder()
			request, _ := http.NewRequest("GET", STATIC_RELATIVE_PATH, nil)
			router.ServeHTTP(recorder, request)

			assert.Equal(t, http.StatusOK, recorder.Code)
			assert.Contains(t, recorder.Body.String(), "<body>")
			assert.Contains(t, recorder.Body.String(), "LogBook")
			assert.Contains(t, recorder.Body.String(), eVar.result+application.API_ROOT_PATH)

		})
	}
}

type envVariables struct {
	host   string
	port   string
	result string
}

func provideEnvVariables() (data []envVariables) {
	data = append(data, envVariables{"localhost", "80", "ws://localhost:80"})
	data = append(data, envVariables{"", "", "ws://localhost:8080"})
	data = append(data, envVariables{"127.0.0.1", "80", "ws://127.0.0.1:80"})
	data = append(data, envVariables{"www.homepage.io", "123", "ws://www.homepage.io:123"})
	return
}
