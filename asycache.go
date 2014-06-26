package asycache

import (
	. "time"
)

type cacheEntity struct {
	data interface {}
	life_time Duration
	last_updated Time
}

type setChanParam struct{k string; v interface{}; finish_chan chan bool;life_time Duration}
type getChanParam struct{k string; clbk chan <- interface{};}

type Cache struct {
	entities map[string]*cacheEntity
	set_chan chan setChanParam
	get_chan chan getChanParam
	//revision_dt Duration //time delta for doing revision of Cache (e.g. removing outdated elements)
}

func MakeCache(cleanup_dt Duration)(*Cache) {
	c := &Cache{entities: make(map[string]*cacheEntity),
		set_chan: make(chan setChanParam, 100),
		get_chan: make(chan getChanParam, 100)}

	go func (){
		ticker := NewTicker(cleanup_dt)
		for {
			select {
			case gcp := <- c.get_chan:
				v, ok := c.entities[gcp.k]
				if !ok {
					gcp.clbk <- nil
					break
				}
				gcp.clbk <- v.data
			case scp := <- c.set_chan:
				c.entities[scp.k] = &cacheEntity{data: scp.v, life_time: scp.life_time, last_updated: Now()}
				scp.finish_chan <- true
			case <- ticker.C:
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

//k is a key
//v is a value of any type - could be also a pointer
//dt - life time for the element
func (c *Cache) Set(k string, v interface{}, dt Duration) chan bool {
	finish_chan := make(chan bool)
	c.set_chan <- setChanParam{k, v, finish_chan, dt}
	return finish_chan
}

func (c *Cache) Get(k string, timeout Duration) (interface{}, bool) {
	clbk := make(chan interface{})
	c.get_chan <- getChanParam{k, clbk}
	ticker := NewTicker(timeout)
	var v interface{}
	v = nil
	for {
		select {
		case v = <- clbk:
			if v != nil {
				return v, true
			}
			return nil, false
		case <- ticker.C:
			return nil, false
		}
	}
	return v, false
}

