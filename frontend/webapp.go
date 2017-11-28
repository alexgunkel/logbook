package frontend

import (
	"strconv"
	"github.com/gin-gonic/gin"
	"html/template"
	"os"
)

func AddFrontend(engine *gin.Engine, templateDir string) {
	fe := &WebApplication{}
	fe.SetTemplateDirPath(templateDir)
	gen := &IdGenerator{}

	engine.GET("/logbook", func(context *gin.Context) {
		fe.InitLogBookClientApplication(context, gen)
	})
}

type WebApplication struct {
	templateFolder string
}

type IdGenerator struct {
	lastIdentifier int64
}

func (app *IdGenerator) getNewIdentifier() string {
	app.lastIdentifier++
	return strconv.FormatInt(app.lastIdentifier, 10)
}

type User struct {
	Identifier string
	Uri string
}

func (a *WebApplication) SetTemplateDirPath(path string)  {
	a.templateFolder = path
}

func (a *WebApplication) InitLogBookClientApplication(c *gin.Context, gen *IdGenerator)  {
	identifier, err := c.Cookie("logbook")
	if nil != err {
		identifier = gen.getNewIdentifier()
		c.SetCookie("logbook", identifier, 0, "", "", false, false)
	}

	user := User{}
	user.Identifier = identifier
	user.Uri = "ws://" + getHost() + ":" + getPort() + "/logbook/" + identifier + "/logs"

	t := template.New("Index.html")
	t, err = t.ParseFiles( a.templateFolder + "/Index.html" )
	if err != nil {
		panic(err)
	}
	err = t.ExecuteTemplate(c.Writer, "Index.html", user)
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
