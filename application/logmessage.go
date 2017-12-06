package application

// This the data structure for incoming messages.
// It consists of header data and the information
// in the http body.
//
// Here our ontology is oriented at the technical
// prospective of http-requests
type IncomingMessage struct {
	logBookId string
	Body      LogMessageBody
	Origin    HeaderData
}

// The body of the message is sent as a json
// object
type LogMessageBody struct {
	Timestamp int         `json:"timestamp"`
	Message   string      `json:"message"`
	Severity  interface{} `json:"severity"`
	Context   interface{} `json:"context"`
}

// Header data contain information about
// the origin of the message, like app and
// logger or the request uri
type HeaderData struct {
	Application string
	LoggerName  string
	RequestUri  string
}

// The LogBook-entry is the transformed data object
// which contains data for the LogBook frontend.
//
// Here we use the genuine LogBook-ontology
type LogBookEntry struct {
	logBookId   string
	Application string
	LoggerName  string
	RequestUri  string
	Timestamp   int         `json:"time"`
	Message     string      `json:"message"`
	Severity    int         `json:"severity"`
	Context     interface{} `json:"context"`
}

type Event struct {
}

func createNewLogMessage(logBookId string) (m *IncomingMessage) {
	m = &IncomingMessage{}
	m.logBookId = logBookId

	return
}

// Mapping of digital and textual versions of
// our loglevels. This map follows the recommendations
// in RFC 5424
//
// see https://tools.ietf.org/html/rfc5424
var severityValues = map[string]int{"debug": 7,
	"informational": 6,
	"notice":        5,
	"warning":       4,
	"error":         3,
	"critical":      2,
	"alert":         1,
	"emergency":     0}

func normalize(input interface{}) (digit int, textual string) {
	if level, ok := input.(string); ok {
		digit = severityValues[level]
		return
	}

	if level, ok := input.(float64); ok {
		digit = int(level)
		return
	}

	if level, ok := input.(int); ok {
		digit = level
		return
	}

	return
}

// This function is responsible for the transition
// from technical terminology to LogBook-ontology
func processMessage(inbound IncomingMessage) (outbound LogBookEntry) {
	outbound.Timestamp = inbound.Body.Timestamp
	outbound.Severity, _ = normalize(inbound.Body.Severity)
	outbound.Message = inbound.Body.Message
	outbound.Context = inbound.Body.Context
	outbound.Application = inbound.Origin.Application
	outbound.LoggerName = inbound.Origin.LoggerName
	outbound.RequestUri = inbound.Origin.RequestUri
	outbound.logBookId = inbound.logBookId

	return
}
