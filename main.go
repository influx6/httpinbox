package main

import (
	"net/http"
	"os"
	"os/signal"

	"github.com/ardanlabs/kit/cfg"
	"github.com/influx6/faux/db/mongo"
	"github.com/influx6/faux/web/app"
	"github.com/influx6/faux/web/middleware"
	"github.com/influx6/httpinbox/app/routes"
)

func main() {

	// Initialize the configuration system to retrieve environment variavles.
	cfg.Init(cfg.EnvProvider{Namespace: "HTTPINBOX"})

	// Set the base level headers all response must contain.
	baseHeaders := map[string]string{"X-App": "HttpInbox"}

	// Make a new mongodb session middleware to be used globally.
	dbm := middleware.MongoDB(mongo.Config{
		Host:     cfg.MustString("MONGO_HOST"),
		AuthDB:   cfg.MustString("MONGO_AUTHDB"),
		DB:       cfg.MustString("MONGO_DB"),
		User:     cfg.MustString("MONGO_USER"),
		Password: cfg.MustString("MONGO_PASS"),
	}, nil)

	// Create the server http.Handler conforming instance.
	inbox := app.New(nil, true, baseHeaders, middleware.Log, dbm)

	// Register all the routes.
	routes.InitRoutes(inbox)

	// Retrieve the address we wish to use for the app.
	addr := cfg.MustString("HOST_ADDR")

	go func() {
		http.ListenAndServe(addr, inbox)
	}()

	// Listen for an interrupt signal from the OS.
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	<-sigChan
}
