package application

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

// a test for our expectations concerning the
// standard behavior of json unmarshaling
func TestUnMarshallSeverity(t *testing.T) {
	origJson := "{\"severity\": 12}"
	origJson2 := "{\"severity\": \"error\"}"
	logEv := &LogMessageBody{}
	logEv2 := &LogMessageBody{}
	json.Unmarshal([]byte(origJson), logEv)
	json.Unmarshal([]byte(origJson2), logEv2)

	assert.Equal(t, float64(12), logEv.Severity)
	assert.Equal(t, "error", logEv2.Severity)
}

// Test normalizing function
func TestNormalize(t *testing.T) {
	for input, output := range dataProviderForNormalization() {
		out, outText := normalize(input)
		assert.Equal(t, output, out, "Expected transformed value to be %v, got %v from %v", output, out, input)
		assert.Equal(t, output, severityValues[outText], "Expected transformed value to be %v, got %v from %v", output, out, input)
	}
}

func TestMessageToLogbookEntryCorrectsLogLevel(t *testing.T) {
	for input, _ := range dataProviderForNormalization() {
		i := IncomingMessage{}
		i.Body.Severity = input

		res := processMessage(i)

		assert.NotNil(t, res.Severity)
		assert.NotNil(t, res.SeverityText)
	}
}

// According to RFC 5424 we have the following associations between numbers
// and log devels:
//
// 0       Emergency: system is unusable
// 1       Alert: action must be taken immediately
// 2       Critical: critical conditions
// 3       Error: error conditions
// 4       Warning: warning conditions
// 5       Notice: normal but significant condition
// 6       Informational: informational messages
// 7       Debug: debug-level messages
//
// see https://tools.ietf.org/html/rfc5424
func dataProviderForNormalization() (m map[interface{}]interface{}) {
	m = make(map[interface{}]interface{})
	m["debug"] = 7
	m["informational"] = 6
	m["notice"] = 5
	m["warning"] = 4
	m["error"] = 3
	m["critical"] = 2
	m["alert"] = 1
	m["emergency"] = 0
	m[float64(1)] = 1
	m[float64(0)] = 0
	m[float64(5)] = 5
	m[int(1)] = 1

	return
}

var result int

func BenchmarkNormalizeString(b *testing.B) {
	var res int
	inc := "debug"
	for n := 0; n < b.N; n++ {
		res, _ = normalize(inc)
	}

	result = res
}

func BenchmarkNormalizeFloat(b *testing.B) {
	var res int
	inc := float64(3)
	for n := 0; n < b.N; n++ {
		res, _ = normalize(inc)
	}

	result = res
}

var resultMsg LogBookEntry

func BenchmarkProcessMessage(b *testing.B) {
	var out LogBookEntry
	in := IncomingMessage{logBookId: "123",
		Body:   LogMessageBody{Timestamp: 123123123, Message: "Message", Severity: "debug"},
		Origin: HeaderData{Application: "myApp", LoggerName: "Logger", RequestUri: "http://www.google.de"}}

	for i := 0; i < b.N; i++ {
		out = processMessage(in)
	}

	resultMsg = out
}
