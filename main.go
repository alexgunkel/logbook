package main

import (
	"github.com/alexgunkel/logbook/application"
	"path/filepath"
)

func main()  {
	path, _ := filepath.Abs("./resources/private/template")
	logBook := application.Application(path)
	logBook.Run()
}
