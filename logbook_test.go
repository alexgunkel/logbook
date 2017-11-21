package logbook

import (
	"net/http"
	"testing"
	"net/http/httptest"
)

func TestDisplayLogBook(t *testing.T)  {
	// Create a first request to log-display
	request, err := http.NewRequest("GET", "/display", nil)
	if nil != err {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(DisplayLogBook)
	handler.ServeHTTP(recorder, request)

	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("Didn't receive status code: Got %v wanted %v", status, http.StatusOK)
	}

	// expect setting cookie
	if contentType := recorder.Header().Get("Content-type"); "application/json" != contentType {
		t.Error("No or wrong content-type set in HTTP response; expected application/json, got %q.", contentType)
	}

	// expect setting cookie
	if "" == recorder.Header().Get("Set-Cookie") {
		t.Error("No cookie set in HTTP response.")
	}
}
