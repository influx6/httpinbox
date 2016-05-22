package routes

import (
	"github.com/influx6/faux/web/app"
	"github.com/influx6/httpinbox/app/api"
)

// InitRoutes registers the routes for this app with the root runner.
func InitRoutes(root *app.App) {
	root.Handle(nil, "GET", "/inbox/:id", api.API.GetInbox)
	root.Handle(nil, "POST", "/inbox", api.API.NewInbox)
	root.Handle(nil, "DELETE", "/inbox", api.API.DestroyInbox)
}
