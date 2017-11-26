package application

type LogMessage struct {
	logBookId string
	Event     LogEvent
	Origin    LogMessageOrigin
}

type LogEvent struct {
	Timestamp int         `json:"timestamp"`
	Message   string      `json:"message"`
	Severity  int         `json:"severity"`
	Context   interface{} `json:"context"`
}

type LogMessageOrigin struct {
	Application string
	LoggerName  string
	RequestUri  string
}

func createNewLogMessage(logBookId string) (m *LogMessage) {
	m = &LogMessage{}
	m.logBookId = logBookId

	return
}
