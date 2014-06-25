package asycache

import "testing"
import (
	"time"
)

func TestICanAddAndGet(t *testing.T) {
	c := MakeCache()
	c.Set("a", "hello")
	c.Set("b", "world")
	v, ok := c.Get("a", 1*time.Second)
	var s string
	if ok {
		s = v.(string)
	}
	if s != "hello" {
		t.Error("Cant get value")
	}
}
