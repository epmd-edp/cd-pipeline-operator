package finalizer

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestContainsStringHelloWorld(t *testing.T) {
	e := ContainsStringHelloWorld([]string{"a", "b"}, "b")
	assert.True(t, e)
}

func TestContainsStringHelloWorldTest(t *testing.T) {
	e := ContainsStringHelloWorldTest([]string{"a", "b"}, "b")
	assert.True(t, e)
}
