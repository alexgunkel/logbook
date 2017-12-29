package frontend

import (
	"html/template"
	"os"

	"github.com/alexgunkel/logbook/application"
	"github.com/gin-gonic/gin"
)

type IndexTemplateData struct {
	Identifier   string
	Uri          string
	PathToStatic string
	Port         string
	EndPoint     string
}

type WebApplication struct {
	templateFolder string
}

func SetApplication(engine *gin.Engine, templateDir string) {
	fe := &WebApplication{}
	fe.SetTemplateDirPath(templateDir)

	engine.GET(STATIC_RELATIVE_PATH, func(context *gin.Context) {
		fe.InitLogBookClientApplication(context)
	})
}

func (a *WebApplication) SetTemplateDirPath(path string) {
	a.templateFolder = path
}

func (a *WebApplication) InitLogBookClientApplication(c *gin.Context) {
	identifier, err := c.Cookie("logbook")
	if nil != err {
		identifier = getNewIdentifier()
		c.SetCookie("logbook", identifier, 0, "", "", false, false)
	}

	templateData := IndexTemplateData{}
	templateData.Identifier = identifier
	templateData.Uri = "ws://" + getHost() + ":" + getPort() + application.API_ROOT_PATH + "/" + identifier + "/logs"
	templateData.Port = getPort()
	templateData.EndPoint = application.API_ROOT_PATH + "/" + identifier + "/logs"
	templateData.PathToStatic = STATIC_BASE_HREF

	t := template.New("Index.html")
	t, err = t.ParseFiles(a.getEntryPoint())
	if err != nil {
		panic(err)
	}
	err = t.ExecuteTemplate(c.Writer, "Index.html", templateData)
	if nil != err {
		panic(err)
	}
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
