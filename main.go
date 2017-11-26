package main

import (
	"github.com/alexgunkel/logbook/application"
	"path/filepath"
	"github.com/alexgunkel/logbook/frontend"
	"github.com/gin-gonic/gin"
)

func main()  {
	engine := gin.Default()
	path, _ := filepath.Abs("./resources/private/template")
	application.AddDispatcher(engine)
	frontend.AddFrontend(engine, path)
	engine.Run()
}
