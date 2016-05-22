package api

import (
	"github.com/influx6/faux/context"
	"github.com/influx6/faux/web/app"
)

// API defines an handle through which the HttpInbox API can be retrieved from.
var API api

// api defines a struct which holds all controller methods for the HttpInbox API.
type api struct{}

// NewInbox handles the creation of a new inbox for the reception of http requests.
func (api) NewInbox(ctx context.Context, wq *app.ResponseRequest) error {
	wq.Write([]byte("New Inbox"))
	return nil
}

// GetInbox retrieves a giving box using the id it recieves.
func (api) GetInbox(ctx context.Context, wq *app.ResponseRequest) error {
	id, _ := wq.Params.Get("id")
	wq.Write([]byte("Inbox: ID[" + id + "]"))
	return nil
}

// DestroyInbox handles the destruction of inbox with all its contents.
func (api) DestroyInbox(ctx context.Context, wq *app.ResponseRequest) error {

	return nil
}
