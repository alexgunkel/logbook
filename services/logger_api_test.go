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
	router := gin.Default()
	incoming := make(chan entities.LogEvent, 20)
	router.POST("/logbook/:client/logs", func(context *gin.Context) {
		err := Log(context, incoming)
		if nil != err {
			t.Fatal(err)
		}
	})
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest("POST", "/logbook/12345/logs", strings.NewReader(getTestJson()))
	if nil != err {
		t.Fatal(err)
	}

	original := &entities.LogEvent{}
	json.Unmarshal([]byte(getTestJson()), original)

	// run
	router.ServeHTTP(recorder, request)
	event := <-incoming

	// validate

	assert.Equal(t, original.Message,   event.Message)
	assert.Equal(t, original.Timestamp, event.Timestamp)
	assert.Equal(t, original.Severity,  event.Severity)

	assert.NotEqual(t, original.Message,   "")
	assert.NotEqual(t, original.Timestamp, 0)
	assert.NotEqual(t, original.Severity,  0)
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
