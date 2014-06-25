package asycache

import (
	. "time"
)

type cacheEntity struct {
	data interface {}
	last_updated Time
}

type setChanParam struct{k string; v interface{}}
type getChanParam struct{k string; clbk chan <- interface{};}

type Cache struct {
	entities map[string]*cacheEntity
	set_chan chan setChanParam
	get_chan chan getChanParam
}

func MakeCache()(*Cache) {
	c := &Cache{set_chan: make(chan setChanParam, 100),
		get_chan: make(chan getChanParam, 100)}

	go func (){
		for {
			select {
			case gcp := <- c.get_chan:
				v, e := c.entities[gcp.k]
				if !e {
					gcp.clbk <- nil
				}
				gcp.clbk <- v
			case scp := <- c.set_chan:
				c.entities[scp.k] = &cacheEntity{data: scp.v, last_updated: Now()}
			}
		}
	}()

	return c
}

func (c *Cache) Set(k string, v interface{}) {
	c.set_chan <- setChanParam{k, v}
	return
}

func (c *Cache) Get(k string, timeout Duration) interface{} {
	clbk := make(chan interface{})
	c.get_chan <- getChanParam{k, clbk}
	ticker := NewTicker(timeout)
	var v interface{}
	v = nil
	for {
		select {
		case v = <- clbk:
		case <- ticker.C:
			break
		}
	}
	return v
}
