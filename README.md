go-base91
=========

This is a work-in-progress implementaiton of Base91 encoding for Go.

Example Usage
-------------
```package main

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
```
