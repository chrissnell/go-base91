go-base91
=========

This is an implementaiton of Base91 encoding for Go.  It is functional but it's a work in progress.  Error handling needs to be improved and there are many optimizations to be made.

Example
-------
	package main
	
	import (
		"fmt"
		"github.com/chrissnell/base91"
	)
	
	func main() {
		str := base91.StdEncoding.EncodeToString([]byte("Hi Evan! Daddy loves you!"))
		fmt.Println(str)
	
		dstr, _ := base91.StdEncoding.DecodeString("KagZQ^]@?F/FaxbjeBC=a^JT%y+&lQE")
		fmt.Println(string(dstr))
	}
