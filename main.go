package main

import "github.com/alexgunkel/logbook/lb-logbook"

func main()  {
	logBook := lb_logbook.Application()
	logBook.Run()
}
