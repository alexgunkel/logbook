package application

type messageDispatcher struct {
	channels map[string]chan Message
	incoming chan Message
}

func (d *messageDispatcher) dispatch()  {
	for {
		postMsg :=<- d.incoming

		if c, err := d.channels[postMsg.logBookId]; err {
			c <- postMsg
		}
	}
}
