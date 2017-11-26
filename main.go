package main

import (
	"github.com/alexgunkel/logbook/lb-logbook"
	"path/filepath"
)

func main()  {
	path, _ := filepath.Abs("./resources/private/template")
	logBook := lb_logbook.Application(path)
	logBook.Run()
}
