package application

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestDispatch(t *testing.T) {
	d := &messageDispatcher{}
	d.incoming = make(chan IncomingMessage, 20)
	d.channels = make(map[string]chan LogBookEntry)
	d.channels["1"] = make(chan LogBookEntry, 20)

	go d.dispatch()

	m := IncomingMessage{logBookId: "1", Body: LogMessageBody{Message: "test", Timestamp: 123, Severity: "debug"}}
	d.incoming<- m

	res :=<- d.channels["1"]

	assert.Equal(t, "test", res.Event.Message)
	assert.Equal(t, 123, res.Event.Timestamp)
}

var resultingMessage LogBookEntry

func BenchmarkDispatch(b *testing.B) {
	var res LogBookEntry
	d := &messageDispatcher{}
	d.incoming = make(chan IncomingMessage, 20)
	d.channels = make(map[string]chan LogBookEntry)
	d.channels["1"] = make(chan LogBookEntry, 20)

	go d.dispatch()

	m := IncomingMessage{logBookId: "1", Body: LogMessageBody{Severity: "debug"}}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		d.incoming<- m

		res =<- d.channels["1"]
	}

	resultingMessage = res
}