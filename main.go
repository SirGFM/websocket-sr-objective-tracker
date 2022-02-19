package main

import (
	"log"
	"os"
	"os/signal"
)

// startServer and configure its signal handler.
func startServer() {
	args := parseArgs()
	_ = args

	intHndlr := make(chan os.Signal, 1)
	signal.Notify(intHndlr, os.Interrupt)

	// TODO: Actually start the server

	<-intHndlr
	log.Printf("Exiting...")
}

func main() {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime)

	defer func() {
		if r := recover(); r != nil {
			log.Fatalf("Application panicked! %+v", r)
		}
	} ()

	startServer()
}
