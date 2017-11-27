package application

import (
	"testing"
	"encoding/json"
	"github.com/stretchr/testify/assert"
)

func TestUnMarshallSeverity(t *testing.T)  {
	origJson := "{\"severity\": 12}"
	origJson2 := "{\"severity\": \"error\"}"
	logEv := &Incoming{}
	logEv2 := &Incoming{}
	json.Unmarshal([]byte(origJson), logEv)
	json.Unmarshal([]byte(origJson2), logEv2)

	assert.Equal(t, float64(12), logEv.Severity)
	assert.Equal(t, "error", logEv2.Severity)
}

func TestNormalize(t *testing.T)  {
	for input, output := range dataProviderForNormalization() {
		inc := Incoming{Severity: input}
		out := inc.normalize()
		assert.Equal(t, output, out.Severity)
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

	return
}