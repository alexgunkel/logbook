package application

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestDispatch(t *testing.T) {
	d := &messageDispatcher{}
	d.incoming = make(chan NewMessage, 20)
	d.channels = make(map[string]chan Message)
	d.channels["1"] = make(chan Message, 20)

	go d.dispatch()

	m := NewMessage{logBookId: "1", Event: Incoming{Message: "test", Timestamp: 123, Severity: "debug"}}
	d.incoming<- m

	res :=<- d.channels["1"]

	assert.Equal(t, "test", res.Event.Message)
	assert.Equal(t, 123, res.Event.Timestamp)
}

func TestDispatchConvertsEventType(t *testing.T) {
	d := &messageDispatcher{}
	d.incoming = make(chan NewMessage, 20)
	d.channels = make(map[string]chan Message)
	d.channels["1"] = make(chan Message, 20)

	go d.dispatch()

	m := NewMessage{logBookId: "1", Event: Incoming{Severity: "debug"}}
	d.incoming<- m

	res :=<- d.channels["1"]

	assert.Equal(t, 7, res.Event.Severity)
}