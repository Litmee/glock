package main

import (
	"glock/src/gl"
	"glock/src/model"
)

func main() {
	l := model.GLock{C: make(chan interface{}, 1)}
	l.C <- true
	model.GlobalLockSet[100] = &l
	gl.Start()
}
