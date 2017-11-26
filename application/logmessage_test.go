package application

import (
	"testing"
	"encoding/json"
	"github.com/stretchr/testify/assert"
)

func TestUnMarshallSeverity(t *testing.T)  {
	t.SkipNow()
	origJson := "{\"severity\": 12}"
	origJson2 := "{\"severity\": \"error\"}"
	logEv := &LogEvent{}
	logEv2 := &LogEvent{}
	json.Unmarshal([]byte(origJson), logEv)
	json.Unmarshal([]byte(origJson2), logEv2)

	assert.Equal(t, 12, logEv.Severity)
	assert.Equal(t, "error", logEv2.Severity)
}