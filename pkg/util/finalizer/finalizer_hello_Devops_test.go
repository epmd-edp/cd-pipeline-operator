package finalizer

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestContainsStringHelloDevops(t *testing.T) {
	e := ContainsStringHelloDevops([]string{"a", "b"}, "b")
	assert.True(t, e)
}

func TestContainsStringHelloDevopsTest(t *testing.T) {
	e := ContainsStringHelloDevopsTest([]string{"a", "b"}, "b")
	assert.True(t, e)
}

