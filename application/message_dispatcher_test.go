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

	d.dispatch(NUMBER_OF_DISPATCHERS)

	m := IncomingMessage{logBookId: "1", Body: LogMessageBody{Message: "test", Timestamp: 123, Severity: "debug"}}
	d.incoming <- m

	res := <-d.channels["1"]

	assert.Equal(t, "test", res.Message)
	assert.Equal(t, 123, res.Timestamp)
}

var resultingMessage LogBookEntry

func BenchmarkDispatch(b *testing.B) {
	var res LogBookEntry
	d := &messageDispatcher{}
	d.incoming = make(chan IncomingMessage, 20)
	d.channels = make(map[string]chan LogBookEntry)
	d.channels["1"] = make(chan LogBookEntry, 20)

	d.dispatch(NUMBER_OF_DISPATCHERS)

	m := IncomingMessage{logBookId: "1", Body: LogMessageBody{Severity: "debug"}}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		d.incoming <- m

		res = <-d.channels["1"]
	}

	resultingMessage = res
}

// We must escape the message because it may contain xml or even html that
// should be shown as plain text in our frontend
func TestMessageInOutgoingMessageWillBeEscaped(t *testing.T) {
	dispatcher := &messageDispatcher{}
	incoming := IncomingMessage{}
	incoming.Body.Message = "<?xml version=\"1.0\" ?><content>blablabla</content>"

	outgoing := dispatcher.processMessage(incoming)

	assert.Equal(t, "&lt;?xml version=&#34;1.0&#34; ?&gt;&lt;content&gt;blablabla&lt;/content&gt;", outgoing.Message)
}

func TestMessageToLogbookEntryCorrectsLogLevel(t *testing.T) {
	dispatcher := &messageDispatcher{}
	for input, inputInt := range dataProviderForNormalization() {
		i := IncomingMessage{}
		i.Body.Severity = input

		res := dispatcher.processMessage(i)

		assert.Equal(t, input, res.SeverityText)
		assert.Equal(t, inputInt, res.Severity)
	}
}

var resultMsg LogBookEntry

func BenchmarkProcessMessage(b *testing.B) {
	dispatcher := &messageDispatcher{}
	var out LogBookEntry
	in := IncomingMessage{logBookId: "123",
		Body:   LogMessageBody{Timestamp: 123123123, Message: "Message", Severity: "debug"},
		Origin: HeaderData{Application: "myApp", LoggerName: "Logger", RequestUri: "http://www.google.de"}}

	for i := 0; i < b.N; i++ {
		out = dispatcher.processMessage(in)
	}

	resultMsg = out
}
