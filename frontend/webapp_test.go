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

func GetDispatcher() *gin.Engine {
	engine := gin.Default()

	return engine
}

func TestServeIndexHtml(t *testing.T) {
	tmp, _ := ioutil.TempDir("", "")
	content := "<html><head>{{.PathToStatic}}</head> <body>{{.Uri}} {{.Identifier}}</body></html>"
	ioutil.WriteFile(tmp+"/Index.html", []byte(content), os.ModePerm)

	os.Setenv(STATIC_APP_DIR_ENV, tmp)
	contents := make(map[string]string)
	contents["{{.PathToStatic}}"] = STATIC_RELATIVE_PATH
	contents["{{.Uri}}"] = "ws://localhost:8080/logbook"

	for in, out := range contents {
		t.Run(in, func(t *testing.T) {
			engine := gin.Default()
			SetApplication(engine, tmp)
			recorder := httptest.NewRecorder()
			request, _ := http.NewRequest("GET", STATIC_RELATIVE_PATH, nil)
			engine.ServeHTTP(recorder, request)

			assert.Equal(t, http.StatusOK, recorder.Code)
			assert.Contains(t, recorder.Body.String(), out)
		})
	}

	t.Run("identifier", func(t *testing.T) {
		engine := gin.Default()
		SetApplication(engine, tmp)
		recorder := httptest.NewRecorder()
		request, _ := http.NewRequest("GET", STATIC_RELATIVE_PATH, nil)
		engine.ServeHTTP(recorder, request)

		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.NotEqual(t, "", recorder.Body.String())
	})
}

func TestServeIndexHtmlWithEndSlash(t *testing.T) {
	tmp, _ := ioutil.TempDir("", "")
	content := "html"
	ioutil.WriteFile(tmp+"/Index.html", []byte(content), os.ModePerm)

	os.Setenv(STATIC_APP_DIR_ENV, tmp+"/")

	engine := gin.Default()
	SetApplication(engine, tmp)
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", STATIC_RELATIVE_PATH, nil)
	engine.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, content, recorder.Body.String())
}

// A normal session starts by a HTTP GET-request at <domain>/logbook. We assume that no cookie is
// set. Therefore, we generate a client-id and set a cookie.
//
// There is no need to redirect. Instead we send a small application that shows the success.
func TestInitLogBookWithoutCookie(t *testing.T) {
	request, err := http.NewRequest("GET", STATIC_RELATIVE_PATH, nil)
	if nil != err {
		t.Fatal(err)
	}

	router := GetDispatcher()
	AddFrontend(router, "../public")
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	logBookId, _ := getRecorderCookie(recorder)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.NotEqual(t, "", logBookId)
}

// When receiving a GET-request to the root page and the client has logbook-cookie, then we
// directly show the app.
func TestInitLogBookWithCookie(t *testing.T) {
	router := GetDispatcher()
	AddFrontend(router, "../public")
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", STATIC_RELATIVE_PATH, nil)
	cookie := &http.Cookie{Name: "logbook", Value: "1234"}
	request.AddCookie(cookie)
	router.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusOK, recorder.Code)
}

// Two different sessions should receive different identifiers
// otherwise we would mix up logs
func TestGetNewIdentifierSetsDifferentIds(t *testing.T) {
	generator := &IdGenerator{}

	assert.NotEqual(t, generator.getNewIdentifier(), generator.getNewIdentifier())
}

// Test the standard logbook frontend-application
// this should be reached at the logbook root path
// and contain:
// * html-body,
// * the given id,
// * the name of the application, and
// * a link to the websocket-service
func TestInitLogBookClientApplication(t *testing.T) {
	router := gin.Default()
	app := &WebApplication{}
	app.SetTemplateDirPath("../public")
	generator := &IdGenerator{}
	router.GET(STATIC_RELATIVE_PATH, func(context *gin.Context) {
		app.InitLogBookClientApplication(context, generator)
	})
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", STATIC_RELATIVE_PATH, nil)
	cookie := &http.Cookie{Name: "logbook", Value: "1234"}
	request.AddCookie(cookie)
	router.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "<body>")
	assert.Contains(t, recorder.Body.String(), "1234")
	assert.Contains(t, recorder.Body.String(), "LogBook")
	assert.Contains(t, recorder.Body.String(), "ws://localhost:8080"+application.API_ROOT_PATH+"/1234/logs")
}

func TestWebApplication_InitLogBookClientApplication_RespectsPort(t *testing.T) {
	os.Setenv("PORT", "1234")
	os.Setenv("HOST", "127.0.0.1")
	router := gin.Default()
	app := &WebApplication{}
	app.SetTemplateDirPath("../public")
	generator := &IdGenerator{}
	router.GET(STATIC_RELATIVE_PATH, func(context *gin.Context) {
		app.InitLogBookClientApplication(context, generator)
	})
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", STATIC_RELATIVE_PATH, nil)
	cookie := &http.Cookie{Name: "logbook", Value: "1234"}
	request.AddCookie(cookie)
	router.ServeHTTP(recorder, request)

	assert.Contains(t, recorder.Body.String(), "ws://127.0.0.1:1234"+application.API_ROOT_PATH+"/1234/logs")
	assert.Contains(t, recorder.Body.String(), "<base href=\"/logbook")
}

// Helper function to get cookie values out of response recorders
func getRecorderCookie(r *httptest.ResponseRecorder) (clientId string, err error) {
	newRequest := &http.Request{Header: http.Header{"Cookie": r.HeaderMap["Set-Cookie"]}}
	clientIdCookie, err := newRequest.Cookie("logbook")

	if nil == err {
		clientId = clientIdCookie.Value
	} else {
		clientId = ""
	}

	return clientId, err
}
