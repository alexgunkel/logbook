package lb_frontend

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestGetNewIdentifierSetsDifferentIdds(t *testing.T) {
	generator := &IdGenerator{}

	assert.NotEqual(t, generator.getNewIdentifier(), generator.getNewIdentifier())
}