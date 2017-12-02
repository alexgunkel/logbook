package frontend

import (
	"github.com/gin-gonic/gin"
	"html/template"
	"os"
	"github.com/alexgunkel/logbook/application"
)

type User struct {
	Identifier string
	Uri        string
	BaseHref   string
}

type WebApplication struct {
	templateFolder string
}

func (a *WebApplication) SetTemplateDirPath(path string) {
	a.templateFolder = path
}

func (a *WebApplication) InitLogBookClientApplication(c *gin.Context, gen *IdGenerator) {
	identifier, err := c.Cookie("logbook")
	if nil != err {
		identifier = gen.getNewIdentifier()
		c.SetCookie("logbook", identifier, 0, "", "", false, false)
	}

	user := User{}
	user.Identifier = identifier
	user.Uri = "ws://" + getHost() + ":" + getPort() + application.API_ROOT_PATH + "/" + identifier + "/logs"
	user.BaseHref = "/logbook"

	t := template.New("Index.html")
	t, err = t.ParseFiles(a.getEntryPoint())
	if err != nil {
		panic(err)
	}
	err = t.ExecuteTemplate(c.Writer, "Index.html", user)
}

func (a *WebApplication) getEntryPoint() string {
	if _, err := os.Stat(a.templateFolder + "/Index.html"); nil == err {
		return a.templateFolder + "/Index.html"
	}

	if _, err := os.Stat(a.templateFolder + "/index.html"); nil == err {
		return a.templateFolder + "/index.html"
	}
	panic("Entry point template not found in " + a.templateFolder)
}

func getPort() string {
	if port := os.Getenv("PORT"); port != "" {
		return port
	}

	return "8080"
}

func getHost() string {
	if host := os.Getenv("HOST"); host != "" {
		return host
	}

	return "localhost"
}
