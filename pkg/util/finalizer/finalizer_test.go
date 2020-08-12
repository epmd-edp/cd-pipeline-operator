package finalizer
import (
   "github.com/stretchr/testify/assert"
   "testing"
)
func TestContainsString(t *testing.T) {
   e := ContainsString([]string{"a", "b"}, "b")
   assert.True(t, e)
}
