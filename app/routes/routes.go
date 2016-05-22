package routes

import (
	"github.com/influx6/faux/web/app"
	"github.com/influx6/httpinbox/app/api"
)

// InitRoutes registers the routes for this app with the root runner.
func InitRoutes(root *app.App) {
	app.PageRoute(root, "GET", "/new", api.API.NewInbox)
}
