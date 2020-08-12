package finalizer

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestContainsStringHelloSonar(t *testing.T) {
	e := ContainsStringHelloSonar([]string{"a", "b"}, "b")
	assert.True(t, e)
}

