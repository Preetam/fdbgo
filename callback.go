package main

import "C"
import "unsafe"

//export GoCallback
func GoCallback(f unsafe.Pointer, p unsafe.Pointer) {
	c := *(*chan bool)(p)
	c <- true
}
