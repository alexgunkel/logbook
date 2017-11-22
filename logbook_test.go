package logbook

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"encoding/json"
	"strings"
	"github.com/stretchr/testify/assert"
)

func TestDisplayLogBook(t *testing.T) {
	// Create a first request to log-display
	request, err := http.NewRequest("GET", "/display", nil)
	if nil != err {
		t.Fatal(err)
	}

	router := SetUpRouter()
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t,	"application/json", recorder.Header().Get("Content-type"))
	assert.NotEqual(t, "", recorder.Header().Get("Set-Cookie"))
}

func TestEmptyLogEventLogEvent(t *testing.T) {
	request, err := http.NewRequest("POST", "/logbook/1234/logs", nil)
	if nil != err {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	router := SetUpRouter()
	router.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestValidLogAccepted(t *testing.T)  {
	router := SetUpRouter()
	recorder := httptest.NewRecorder()
	requestBody,_ := json.Marshal(logEvent{"Test"})
	request, err := http.NewRequest("POST", "/logbook/12345/logs", strings.NewReader(string(requestBody)))
	if nil != err {
		t.Fatal(err)
	}

	router.ServeHTTP(recorder, request)
	assert.Equal(t, http.StatusOK, recorder.Code)
}
