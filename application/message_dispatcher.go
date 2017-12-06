package application

type messageDispatcher struct {
	channels map[string]chan LogBookEntry
	incoming chan IncomingMessage
}

func (d *messageDispatcher) dispatch()  {
	for {
		inMsg :=<-d.incoming
		postMsg := processMessage(inMsg)

		if c, err := d.channels[postMsg.logBookId]; err {
			c <- postMsg
		}
	}
}
