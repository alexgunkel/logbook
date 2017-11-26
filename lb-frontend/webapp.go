package lb_frontend

import (
	"strconv"
	"net/http"
	"github.com/gin-gonic/gin"
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

	indexAction(c.Writer, c.Request)
}

func indexAction(w http.ResponseWriter, r *http.Request)  {
	
}