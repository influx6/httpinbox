package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"

	"github.com/influx6/httpinbox/app"
)

func main() {

	// Retrieve the environment vairiables needed for the
	// address and storage path for our inbox.
	addr := os.Getenv("HTTPINBOX_LISTEN")
	dataDir := os.Getenv("HTTPINBOX_DATA")
	viewsDir := os.Getenv("HTTPINBOX_VIEWS")

	if dataDir == "" {
		panic("Require valid path for inbox store")
	}

	if viewsDir == "" {
		panic("Require valid path for app views")
	}

	inbox := app.New(dataDir, viewsDir)

	go func() {
		if err := http.ListenAndServe(addr, inbox); err != nil {
			fmt.Printf("Server Error: %s\n", err)
		}
	}()

	// Listen for an interrupt signal from the OS.
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	<-sigChan
}
