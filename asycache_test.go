package asycache

import "testing"
import (
	"time"
)

func TestICanAddAndGet(t *testing.T) {
	c := MakeCache(1*time.Minute)
	data := map[string]string{"k1": "v1", "k2": "v2", "k3": "v3", "k4": "v4", "k5": "v5"}
	fchans := make([] chan bool, len(data))
	i := 0
	for k, v := range data {
		fchans[i] = c.Set(k, v, 1*time.Minute)
		i++
	}
	//waiting for setting operation be accomplished
	for i = 0; i < len(fchans); i++ {
		<- fchans[i]
	}
	for k, v := range data {
		cv, ok := c.Get(k, 1*time.Second)
		var s string
		if ok {
			s = cv.(string)
		}
		if s != v {
			t.Errorf("Map value %v != Cache value %v with key %v", v, s, k)
		}
	}
}

func TestICanAddAndReplaceElement(t *testing.T) {
	c := MakeCache(1*time.Minute)
	<- c.Set("a", "hello", 1*time.Minute)
	<- c.Set("a", "world", 1*time.Minute)
	v, ok := c.Get("a", 1*time.Second)
	if !ok {
		t.Error("Element not found")
		return
	}
	s := v.(string)
	if s != "world" {
		t.Errorf("string has value = %v, but should world", s)
	}
}

//test requires some time about 5 seconds to execute
func TestThatOutdatedItemsWillBeRemovedInTime(t *testing.T) {
	c := MakeCache(1*time.Second)
	<- c.Set("a", "hello", 3*time.Second)
	<- c.Set("b", "world", 10*time.Second)
	time.Sleep(5*time.Second)
	_, ok := c.Get("a", 1*time.Second)
	if ok {
		t.Error("Element should be deleted")
		return
	}
	_, ok = c.Get("b", 1*time.Second)
	if !ok {
		t.Error("Element should not be deleted")
		return
	}
}
