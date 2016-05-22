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

	return nil
}

// DestroyInbox handles the destruction of inbox with all its contents.
func (api) DestroyInbox(ctx context.Context, wq *app.ResponseRequest) error {

	return nil
}
