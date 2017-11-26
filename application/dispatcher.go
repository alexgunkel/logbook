package application

type dispatcher struct {
	channels map[string]chan PostMessage
	incoming chan PostMessage
}

func (d *dispatcher) dispatch()  {
	for {
		postMsg :=<- d.incoming

		if c, err := d.channels[postMsg.logBookId]; err {
			c <- postMsg
		}
	}
}
