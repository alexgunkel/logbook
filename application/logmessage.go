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
	Timestamp int         `json:"time"`
	Message   string      `json:"message"`
	Severity  interface{} `json:"severity"`

	// The context is not a natural part of
	// a log message in general, but only for
	// psr-loggers in PHP. This should therefore
	// be eliminated in the future
	Context interface{} `json:"context"`
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
	logBookId    string `json:"log_book_id"`
	Application  string `json:"application"`
	LoggerName   string `json:"logger"`
	RequestUri   string `json:"request_uri"`
	Timestamp    int    `json:"time"`
	Message      string `json:"message"`
	Severity     int    `json:"severity"`
	SeverityText string `json:"severity_text"`

	// The Context will be eliminated in a future version
	// because it only affects the working of psr-compatible
	// loggers in PHP
	Context interface{} `json:"context"`
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
var severityValuesIntString = map[int]string{7: "debug",
	6: "informational",
	5: "notice",
	4: "warning",
	3: "error",
	2: "critical",
	1: "alert",
	0: "emergency"}

// Function to normalize severity levels
// Takes an input in int, float or string
// returns loglevel as int and string
func analyzeLogLevel(input interface{}) (int, string) {
	if level, ok := input.(string); ok {
		return severityValues[level], level
	}

	if level, ok := input.(float64); ok {
		digit := int(level)
		return digit, severityValuesIntString[digit]
	}

	if level, ok := input.(int); ok {
		return level, severityValuesIntString[level]
	}

	return -1, ""
}
