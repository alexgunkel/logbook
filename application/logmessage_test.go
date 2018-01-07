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
func TestNormalizeStartingWithText(t *testing.T) {
	var output int
	for input, out := range dataProviderForNormalization() {
		if res, ok := out.(float64); ok {
			output = int(res)
		}
		if res, ok := out.(int); ok {
			output = res
		}
		outInt, outText := analyzeLogLevel(input)
		assert.Equal(t, output, outInt, "Expected transformed value to be %v, got %v from %v", output, outInt, input)
		assert.Equal(t, output, severityValues[outText], "Expected transformed value to be %v, got %v from %v", output, outInt, input)
	}
}

// Test normalizing function
func TestNormalizeStartingWithInt(t *testing.T) {
	var input int
	for output, in := range dataProviderForNormalization() {
		if res, ok := in.(int); ok {
			input = res
		} else {
			t.Error("Error in assumption that input is integer.")
		}
		out, outText := analyzeLogLevel(input)
		assert.Equal(t, input, out, "Expected output value to be numerical %v, got %v from %v", input, out, input)
		assert.Equal(t, outText, output, "Expected output value to be string %v, got %v from %v", output, outText, input)
	}
}

// Test normalizing function
func TestNormalizeStartingWithFloat(t *testing.T) {
	var input float64
	for output, in := range dataProviderForNormalization() {
		if res, ok := in.(int); ok {
			input = float64(res)
		}
		out, outText := analyzeLogLevel(input)
		assert.Equal(t, in, out, "Expected output value to be numerical %v, got %v from %v", in, out, input)
		assert.Equal(t, outText, output, "Expected output value to be string %v, got %v from %v", output, outText, input)
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

	return
}

//
// Benchmarking
//
var result int

func BenchmarkNormalizeString(b *testing.B) {
	var res int
	inc := "debug"
	for n := 0; n < b.N; n++ {
		res, _ = analyzeLogLevel(inc)
	}

	result = res
}

func BenchmarkNormalizeFloat(b *testing.B) {
	var res int
	inc := float64(3)
	for n := 0; n < b.N; n++ {
		res, _ = analyzeLogLevel(inc)
	}

	result = res
}
