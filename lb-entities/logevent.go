package lb_entities

type PostMessage struct {
	Header LogHeader
	Event  LogEvent
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
