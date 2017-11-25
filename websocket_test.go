package logbook

import (
	"github.com/posener/wstest"
	"testing"
	"net/http"
)

func TestHandler(t *testing.T) {
	var err error

	h := Application()
	d := wstest.NewDialer(h, nil)  // or t.Log instead of nil

	c, resp, err := d.Dial("ws://localhost/logbook/123/ws", nil)
	if err != nil {
		t.Fatal(err)
	}

	if got, want := resp.StatusCode, http.StatusSwitchingProtocols; got != want {
		t.Errorf("resp.StatusCode = %q, want %q", got, want)
	}

	err = c.WriteJSON("test")
	if err != nil {
		t.Fatal(err)
	}
}
