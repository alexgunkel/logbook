package application

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// This defines the default api root for incoming
// log messages.
const API_ROOT_PATH = "/api/v1/logbooks"

// seize of buffer for channels
const CHANNEL_BUFFER = 20

// This is the main backend application
// It contains the inbox and dispatcher
// and manages the processing of incoming
// log messages
type LogBookApplication struct {
	receiver   *inbox
	dispatcher *messageDispatcher
	engine     *gin.Engine
}

// Use this function to register the LogBookAppliccation
// to a given GIN-engine. It will
//
// * initialize inbox and dispatcher,
// * start the log inbox, and
// * start the LogBook API
func (app *LogBookApplication) AddApplicationToEngine(engine *gin.Engine) {
	app.engine = engine
	app.initAndStartDispatcher()

	app.startInbox()

	app.startLogBook()
}

// Creates the channel between the inbox and the message
// dispatcher
func (app *LogBookApplication) createChannelToDispatcher() {
	receiverToDispatcher := make(chan IncomingMessage, CHANNEL_BUFFER)
	app.receiver.chanelToMessageDispatcher = receiverToDispatcher
	app.dispatcher.incoming = receiverToDispatcher
}

// Initialize inbox and dispatcher, implement connection between
// them and start the dispatcher
func (app *LogBookApplication) initAndStartDispatcher() {
	app.receiver = &inbox{}
	app.dispatcher = &messageDispatcher{}

	// Create channel between inbox and messageDispatcher
	app.createChannelToDispatcher()

	app.dispatcher.channels = make(map[string]chan LogBookEntry)
	go app.dispatcher.dispatch()
}

// This starts the POST-receiver
func (app *LogBookApplication) startInbox() {
	app.engine.POST(API_ROOT_PATH+"/:client/logs", func(context *gin.Context) {
		logBookId := context.Param("client")
		app.receiver.submit(context, logBookId)
		context.Status(http.StatusAccepted)
	})
}

// This start the API for the LogBook
func (app *LogBookApplication) startLogBook() {
	app.engine.GET(API_ROOT_PATH+"/:client/logs", func(context *gin.Context) {
		logBookId := context.Param("client")
		ForceCookie(context, logBookId)

		outbound := make(chan LogBookEntry, CHANNEL_BUFFER)
		app.dispatcher.channels[logBookId] = outbound
		defer func() {
			delete(app.dispatcher.channels, logBookId)
		}()

		lb, _ := createLogBook(context.Writer, context.Request, outbound)
		lb.listen()
	})
}
