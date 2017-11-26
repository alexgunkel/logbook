package lb_frontend

import (
	"strconv"
	"github.com/gin-gonic/gin"
	"html/template"
)

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

func InitLogBookClientApplication(c *gin.Context, gen *IdGenerator)  {
	identifier, err := c.Cookie("logbook")
	if nil != err {
		identifier = gen.getNewIdentifier()
		c.SetCookie("logbook", identifier, 0, "", "", false, false)
	}

	user := User{}
	user.Identifier = identifier
	user.Uri = "ws://localhost/logbook/" + identifier + "/logs"

	t := template.New("Index.html")
	t, err = t.ParseFiles( "./../resources/private/template/Index.html" )
	if err != nil {
		panic(err)
	}
	err = t.ExecuteTemplate(c.Writer, "Index.html", user)
}