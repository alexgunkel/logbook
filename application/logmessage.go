package application

type NewMessage struct {
	logBookId string
	Event     Incoming
	Origin    Origin
}

type Message struct {
	logBookId string
	Event     Event
	Origin    Origin
}

type Incoming struct {
	Timestamp int         `json:"timestamp"`
	Message   string      `json:"message"`
	Severity  interface{} `json:"severity"`
	Context   interface{} `json:"context"`
}

type Event struct {
	Timestamp int         `json:"timestamp"`
	Message   string      `json:"message"`
	Severity  int         `json:"severity"`
	Context   interface{} `json:"context"`
}

type Origin struct {
	Application string `json:"application"`
	LoggerName  string `json:"logger_name"`
	RequestUri  string `json:"request_uri"`
}

func createNewLogMessage(logBookId string) (m *NewMessage) {
	m = &NewMessage{}
	m.logBookId = logBookId

	return
}

var severityValues = map[string]int{"debug": 7,
	"informational": 6,
	"notice": 5,
	"warning": 4,
	"error": 3,
	"critical": 2,
	"alert": 1,
	"emergency": 0}

func (i Incoming) normalize() (e Event) {
	e = copyEvent(&i)
	if level, ok := i.Severity.(string); ok {
		e.Severity = severityValues[level]
		return
	}

	if level, ok := i.Severity.(float64); ok {
		e.Severity = int(level)
		return
	}

	if level, ok := i.Severity.(int); ok {
		e.Severity = level
		return
	}

	return
}

func processMessage(inbound NewMessage) (outbound Message) {
	outbound.logBookId = inbound.logBookId
	outbound.Origin = inbound.Origin
	outbound.Event = inbound.Event.normalize()

	return
}

func copyEvent(i *Incoming) (e Event) {
	e.Timestamp = i.Timestamp
	e.Message = i.Message
	e.Context = i.Context

	return
}
