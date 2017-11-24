package services

import (
	"testing"
	"github.com/gin-gonic/gin"
	"github.com/alexgunkel/logbook/entities"
	"net/http/httptest"
	"strings"
	"github.com/stretchr/testify/assert"
	"encoding/json"
	"net/http"
)

func TestValidLogSentToDispatcher(t *testing.T) {
	// setup
	request, err := http.NewRequest("POST", "/logbook/12345/logs", strings.NewReader(getTestJson()))
	if nil != err {
		t.Fatal(err)
	}

	router := gin.Default()
	incoming := make(chan entities.PostMessage, 20)
	router.POST("/logbook/:client/logs", func(context *gin.Context) {
		err := Log(context, incoming)
		if nil != err {
			t.Fatal(err)
		}
	})
	recorder := httptest.NewRecorder()

	original := &entities.LogEvent{}
	json.Unmarshal([]byte(getTestJson()), original)

	// run
	router.ServeHTTP(recorder, request)
	postMessage := <-incoming

	// validate

	assert.Equal(t, original.Message,   postMessage.Event.Message)
	assert.Equal(t, original.Timestamp, postMessage.Event.Timestamp)
	assert.Equal(t, original.Severity,  postMessage.Event.Severity)

	assert.NotEqual(t, postMessage.Event.Message,   "")
	assert.NotEqual(t, postMessage.Event.Timestamp, 0)
	assert.NotEqual(t, postMessage.Event.Severity,  0)
}

func TestLogStoresHeaderDataInLogInfo(t *testing.T)  {
	// setup
	request, err := http.NewRequest("POST", "/logbook/12345/logs", strings.NewReader(getTestJson()))
	if nil != err {
		t.Fatal(err)
	}
	request.Header.Set(entities.LogHeaderAppIdentifier,  "MyApp")
	request.Header.Set(entities.LogHeaderLoggerName,  "MyLogger")
	request.Header.Set(entities.LogHeaderRequestUri,  "https://www.logbook.io")

	router := gin.Default()
	incoming := make(chan entities.PostMessage, 20)
	router.POST("/logbook/:client/logs", func(context *gin.Context) {
		err := Log(context, incoming)
		if nil != err {
			t.Fatal(err)
		}
	})
	recorder := httptest.NewRecorder()

	// run
	router.ServeHTTP(recorder, request)
	postMessage := <-incoming

	// validate
	assert.NotNil(t, postMessage.Header)
}

// Helperfunctions to make testing easier.
// They are mainly thought of as dataproviders.

// Helper to build a simple JSON string for testing
func getTestJson() string {
	res, _ := json.Marshal(struct {
		Message   string
		Timestamp int
		Severity  int
	}{
		"Test",
		123123123,
		3,
	})

	return string(res)
}
