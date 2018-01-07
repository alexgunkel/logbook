package application

import "html"

type messageDispatcher struct {
	// for every registered logbook there is
	// one channel to the LogBook
	channels map[string]chan LogBookEntry

	// One channel for incoming messages
	// This has to be watched for log messages
	incoming chan IncomingMessage
}

// Central function to start a worker
func (d *messageDispatcher) dispatch(numberOfWorkers int) {
	for index := 0; index < numberOfWorkers; index++ {
		go d.createWorker()
	}
}

func (d *messageDispatcher) createWorker() {
	for {
		inMsg := <-d.incoming
		postMsg := d.processMessage(inMsg)

		if c, ok := d.channels[postMsg.logBookId]; ok {
			c <- postMsg
		}
	}
}

// This function is responsible for the transition
// from technical terminology to LogBook-ontology
func (dispatcher *messageDispatcher) processMessage(inbound IncomingMessage) (outbound LogBookEntry) {
	outbound.Timestamp = inbound.Body.Timestamp

	// Our loglevels must be normalized
	outbound.Severity, outbound.SeverityText = analyzeLogLevel(inbound.Body.Severity)

	// The message must be escaped in order to show xml content
	outbound.Message = html.EscapeString(inbound.Body.Message)

	outbound.Context = inbound.Body.Context
	outbound.Application = inbound.Origin.Application
	outbound.LoggerName = inbound.Origin.LoggerName
	outbound.RequestUri = inbound.Origin.RequestUri
	outbound.RequestId = inbound.Origin.RequestId
	outbound.logBookId = inbound.logBookId

	return
}
