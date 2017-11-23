package entities

type LogEvent struct {
	Timestamp int    `json:"timestamp"`
	Message   string `json:"message"`
	Severity  int    `json:"severity"`
}
