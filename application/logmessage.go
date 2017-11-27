package application

type Message struct {
	logBookId string
	Event     Event
	Origin    Origin
}

type Event struct {
	Timestamp int         `json:"timestamp"`
	Message   string      `json:"message"`
	Severity  int         `json:"severity"`
	Context   interface{} `json:"context"`
}

type Origin struct {
	Application string
	LoggerName  string
	RequestUri  string
}

func createNewLogMessage(logBookId string) (m *Message) {
	m = &Message{}
	m.logBookId = logBookId

	return
}
