package main

import (
	"os"
	"path/filepath"

	"github.com/alexgunkel/logbook/application"
	"github.com/alexgunkel/logbook/frontend"
	"github.com/gin-gonic/gin"
)

func main() {
	engine := gin.Default()

	// Orchestrate the dispatcher stuff
	dispatcher := &application.LogBookApplication{}
	dispatcher.AddApplicationToEngine(engine)

	// Orchestrate the frontend stuff
	frontend.AddFrontend(engine, getAppDirEnv())

	// Start the engine
	engine.Run()
}

func getAppDirEnv() string {
	var path string
	if "" != os.Getenv(frontend.STATIC_APP_DIR_ENV) {
		path, _ = filepath.Abs(os.Getenv(frontend.STATIC_APP_DIR_ENV))
	} else {
		path, _ = filepath.Abs("./public")
		os.Setenv(frontend.STATIC_APP_DIR_ENV, path)
	}
	return path
}
