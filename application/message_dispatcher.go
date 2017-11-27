package application

type messageDispatcher struct {
	channels map[string]chan Message
	incoming chan NewMessage
}

func (d *messageDispatcher) dispatch()  {
	for {
		postMsg := processMessage(<-d.incoming)

		if c, err := d.channels[postMsg.logBookId]; err {
			c <- postMsg
		}
	}
}
