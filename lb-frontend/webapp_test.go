package lb_frontend

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"github.com/gin-gonic/gin"
	"net/http"
)

func TestGetNewIdentifierSetsDifferentIdds(t *testing.T) {
	generator := &IdGenerator{}

	assert.NotEqual(t, generator.getNewIdentifier(), generator.getNewIdentifier())
}

func TestInitLogBookClientApplication(t *testing.T) {
	router := gin.Default()
	generator := &IdGenerator{}
	router.GET("/logbook", func(context *gin.Context) {
		InitLogBookClientApplication(context, generator)
	})
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/logbook", nil)
	cookie := &http.Cookie{Name: "logbook", Value: "1234"}
	request.AddCookie(cookie)
	router.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "<body>")
	assert.Contains(t, recorder.Body.String(), "1234")
	assert.Contains(t, recorder.Body.String(), "Hello World!")
	assert.Contains(t, recorder.Body.String(), "ws:\\/\\/localhost\\/logbook\\/1234\\/logs")
}