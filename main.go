package main

import (
	"net/http"
	"os"
	"os/signal"

	"github.com/dimfeld/httptreemux"
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

	mux := httptreemux.New()
	mux.GET("/", inbox.GetAllInbox)
	mux.POST("/inbox", inbox.NewInbox)
	mux.GET("/inbox/:id", inbox.GetInbox)

	for _, method := range []string{"POST", "DELETE", "PUT", "PATCH", "HEAD"} {
		mux.Handle(method, "/inbox/:id", inbox.AddToInbox)
	}

	mux.GET("/inbox/:id/:reqid", inbox.GetInboxItem)

	go func() {
		http.ListenAndServe(addr, mux)
	}()

	// Listen for an interrupt signal from the OS.
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	<-sigChan
}
