package main

import (
	"github.com/alexgunkel/logbook/application"
	"path/filepath"
	"github.com/alexgunkel/logbook/frontend"
	"github.com/gin-gonic/gin"
	"os"
)

func main()  {
	engine := gin.Default()

	// Orchestrate the dispatcher stuff
	dispatcher := &application.LogBookApplication{}
	dispatcher.AddApplicationToEngine(engine)

	// Orchestrate the frontend stuff
	if "" != os.Getenv(frontend.STATIC_APP_DIR_ENV) {
		frontend.SetStaticApp(engine)
	} else {
		path, _ := filepath.Abs("./resources/private/template")
		frontend.AddFrontend(engine, path)
	}

	// Start the engine
	engine.Run()
}
