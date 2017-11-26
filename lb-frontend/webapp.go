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

func InitLogBookClientApplication(c *gin.Context, gen *IdGenerator)  {
	identifier, err := c.Cookie("logbook")
	if nil != err {
		identifier = gen.getNewIdentifier()
		c.SetCookie("logbook", identifier, 0, "", "", false, false)
	}

	t := template.New("index.tmp")
	t, err = t.ParseFiles( "./../lb-logbook/index.tmp" )
	if err != nil {
		panic(err)
	}
	err = t.ExecuteTemplate(c.Writer, "index.tmp", identifier)
}