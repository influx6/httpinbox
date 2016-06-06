package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"

	"github.com/influx6/httpinbox/app"
)

const (
	get  = "get"
	post = "post"
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

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		boxLen := len("/inbox")

		if len(r.URL.Path) < boxLen {
			inbox.GetAllInbox(w, r, nil)
			return
		}

		var paramsPieces []string

		params := r.URL.Path[boxLen:]
		params = strings.TrimSpace(strings.TrimPrefix(strings.TrimSuffix(params, "/"), "/"))

		if params != "" {
			paramsPieces = strings.Split(params, "/")
		}

		switch len(paramsPieces) {
		case 0:
			switch strings.ToLower(r.Method) {
			case get:
				inbox.GetAllInbox(w, r, nil)
				return
			case post:
				inbox.NewInbox(w, r, nil)
				return
			default:
				w.WriteHeader(http.StatusMethodNotAllowed)
				return
			}

		case 1:
			id := paramsPieces[0]

			switch strings.ToLower(r.Method) {
			case get:
				inbox.GetInbox(w, r, map[string]string{"id": id})
				return
			default:
				inbox.AddToInbox(w, r, map[string]string{"id": id})
				return
			}

		case 2:
			id := paramsPieces[0]
			reqid := paramsPieces[1]

			switch strings.ToLower(r.Method) {
			case get:
				inbox.GetInboxItem(w, r, map[string]string{"id": id, "reqid": reqid})
				return
			default:
				w.WriteHeader(http.StatusMethodNotAllowed)
				return
			}
		default:
			w.WriteHeader(http.StatusBadRequest)
			return

		}
	})

	go func() {
		if err := http.ListenAndServe(addr, nil); err != nil {
			fmt.Printf("Server Error: %s\n", err)
		}
	}()

	// Listen for an interrupt signal from the OS.
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	<-sigChan
}
