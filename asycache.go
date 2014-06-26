package asycache

import (
	. "time"
)

// Wrapper type that instances stored in cache
type cacheEntity struct {
	//object that was asked to be stored in cache
	data interface{}
	//life_time after that object will be removed from cache
	life_time Duration
	//last time when object was requested
	last_updated Time
}

type setChanParam struct {
	k           string
	v           interface{}
	finish_chan chan bool
	life_time   Duration
}
type getChanParam struct {
	k    string
	clbk chan<- interface{}
}

//Cache type that contains data used by cache
type Cache struct {
	//storage object for given data to store
	//stores pointer to entities for greater speed up access
	entities map[string]*cacheEntity
	//channel to set new data object for key
	set_chan chan setChanParam
	//channel to get object by key
	get_chan chan getChanParam
}

// Returns pointer to the new object of type Cache
// Each cleanup_dt nanoseconds cache service will
// remove outdated objects from the storage.
func MakeCache(cleanup_dt Duration) *Cache {
	c := &Cache{entities: make(map[string]*cacheEntity),
		set_chan: make(chan setChanParam, 100),
		get_chan: make(chan getChanParam, 100)}

	go func() {
		ticker := NewTicker(cleanup_dt)
		for {
			select {
			case gcp := <-c.get_chan:
				v, ok := c.entities[gcp.k]
				if !ok {
					gcp.clbk <- nil
					break
				}
				gcp.clbk <- v.data
				c.entities[gcp.k].last_updated = Now()
			case scp := <-c.set_chan:
				c.entities[scp.k] = &cacheEntity{data: scp.v, life_time: scp.life_time, last_updated: Now()}
				scp.finish_chan <- true
			case <-ticker.C:
				for k, v := range c.entities {
					//remove element if it was not updated longer than it's life time
					if Now().Sub(v.last_updated) >= v.life_time {
						delete(c.entities, k)
					}
				}
			}
		}
	}()

	return c
}

//Cache function that requests cache to set new element v
//with key k to the cache storage or replace existing element
//with the same key.
//dt - life time for the element, if element was not requested
//during that time it will be removed
func (c *Cache) Set(k string, v interface{}, dt Duration) chan bool {
	finish_chan := make(chan bool)
	c.set_chan <- setChanParam{k, v, finish_chan, dt}
	return finish_chan
}

func (c *Cache) Get(k string, timeout Duration) (interface{}, bool) {
	clbk := make(chan interface{}, 2)
	c.get_chan <- getChanParam{k, clbk}
	ticker := NewTicker(timeout)
	var v interface{}
	v = nil
	for {
		select {
		case v = <-clbk:
			if v != nil {
				return v, true
			}
			return nil, false
		case <-ticker.C:
			return nil, false
		}
	}
	return v, false
}
