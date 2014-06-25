Asycache
========

Intro
------

Asycache is a asynchronous cache that does not uses mutual locks. Instead it uses goroutine that manages access to the stored data.

Features
--------

* No locks
* Uses goroutines

Usage
-----

To create object just create Cache type instance.
Interface:
	* set(key, object) - adds new object to cache with given key. If object exists replaces it.
		Function returns immediatlly.
	* get(key) : object, exists - returns object from cache. If object does not exists returns nil and exists to true.
		Function waits until cache processing goroutine returns.

Example
-------

c := Cache()
c.set("mykey", &struct{s string}{"world"})
obj, exists := c.get("mykey")
if exists {
	fmt.Println(obj.s)
}

Copyright, License & Contributors
=================================

Use it for free with no restrictions.