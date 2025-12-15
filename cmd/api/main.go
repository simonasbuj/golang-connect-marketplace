// Package main is the entry point for the application.
package main

import (
	"log"
	"time"
)

const timeOutSeconds = 2

func main() {
	count := 0
	for {
		log.Printf("this is going to be awesome: %d\n", count)
		count++

		time.Sleep(time.Second * timeOutSeconds)
	}
}
