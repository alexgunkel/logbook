package application

type PostMessage struct {
	Header    LogHeader
	Event     LogEvent
	logBookId string
}

type LogEvent struct {
	Timestamp int    `json:"timestamp"`
	Message   string `json:"message"`
	Severity  int    `json:"severity"`
}

type LogHeader struct {
	Application string
	LoggerName  string
	RequestUri  string
}
