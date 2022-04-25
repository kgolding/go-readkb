package main

import (
	"fmt"
	"github/kgolding/go-readkb"
	"os"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s <path>\ne.g. %s /dev/input/event0\n", os.Args[0], os.Args[0])
		return
	}

	k, err := readkb.NewFromPath(os.Args[1])
	if err != nil {
		println(err.Error())
		return
	}
	for e := range k.C {
		fmt.Printf(">>> '%s' h%X %d (Scancode %X)\n", string(e.Char), e.Char, e.Char, e.Scancode)
	}
	k.Close()
	if k.Err != nil {
		println(k.Err.Error())
	}
	time.Sleep(time.Second)
}
