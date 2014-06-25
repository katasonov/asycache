package asycache

import "testing"
import (
	"time"
)

func TestAddAndGet(t *testing.T) {
	c := MakeCache()
	c.Set("a", interface {}(string("hello")))
	c.Set("b", "world")
	s := c.Get("a", 1*time.Second)
	if s == nil || s != "a" {
		t.Error("Cant get value")
	}
}
