package lb_frontend

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"errors"
)

type webContextMock struct {
	cookie string
	err error
	status int
	location string
}

func (c *webContextMock) SetCookie(name, value string, maxAge int, path, domain string, secure, httpOnly bool) {
	c.cookie = value
}

func (c *webContextMock) Cookie(name string) (string, error) {
	return c.cookie, c.err
}

func (c *webContextMock) Redirect(code int, location string)  {
	c.status = code
	c.location = location
}

func TestInitLogBookClientApplicationSetsDifferentCookies(t *testing.T) {
	generator := &IdGenerator{}
	context := webContextMock{}
	context.err = errors.New("asd")
	contextTwo := webContextMock{}
	contextTwo.err = errors.New("asd")
	InitLogBookClientApplication(&context, generator)
	InitLogBookClientApplication(&contextTwo, generator)

	assert.NotEqual(t, context.cookie, contextTwo.cookie)
}