package main

import (
	"fmt"
	"time"
)

func main() {
	count := 0
	for {
		fmt.Printf("this is going to be awesome: %d\n", count)
		count ++
		time.Sleep(time.Second * 2)
	}
}