package main

import (
	"fmt"
	"glock/src/gl"
	"time"
)

func main() {
	nanosecond := int32(time.Now().Nanosecond())
	fmt.Println(nanosecond)
	gl.Start()
}
