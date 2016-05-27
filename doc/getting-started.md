# Building your own HttpInbox
This article is about showing how you can build your own application using Go and take it from completion to deployment using Harrow to simplify your life, this allows you a insight into how you can combine different technology with Harrow to seamlessly move from product iteration to product continuous deployment and integration.

We will be building a simple application that logs http requests coming
to it by allocating `inboxes` with URLs that allows you capture and store which requests are hitting your servers, this example will be very contrived but with this as a base more complex versions can indeed be built.

The concept is simple: A service that allows you to create a http request inbox where clients could hit and we can capture that requests
to allow us perform whatever metric or process we wish to do so with that data.

Our application really has a very simple structure and this is very
important, because the structure of our code has a massive effect on
our deployment strategy on a longer term.

```bash
~/.../httpinbox > tree -d .
.
├── app
│   ├── api.go
│   ├── datawriter.go
│   └── views
│       ├── all.tml
│       ├── layout.tml
│       ├── list.tml
│       └── single.tml
├── main.go
├── readme.md
└── vendor
```

We have a `app` directory where the controller code for our application with their view templates in `app/views` and also a `main.go` file which will be used to both run our application and build our binary for the project when deploying. In accordance with the new go1.6 approach, we also have a `vendor` directory where all source code based dependencies are stored, this allows us to lock down our dependencies with other libraries.

Our application `main.go` file is rather simple, it loads off specific configuration from the environment which allows us to set different configurations depending on the platform be it for `development-testing` or `production` without touching the codebase.

```go
// main.go

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
	mux.GET("/inbox/:id/:reqid", inbox.GetInboxItem)

	for _, method := range []string{"POST", "DELETE", "PUT", "PATCH", "HEAD"} {
		mux.Handle(method, "/inbox/:id", inbox.AddToInbox)
	}


	go func() {
		http.ListenAndServe(addr, mux)
	}()

	// Listen for an interrupt signal from the OS.
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	<-sigChan
}
```

Our configurations are loaded from the environment
using the `os` native package.

```go
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
```

 I personally believe that the use of panics should be done sparingly and only during the initial, start up phase of an application, because using it during the runtime phase simplify kills your application, loosing a lot of context, data and providing adequate headaches for future runs.

```go
	mux := httptreemux.New()
```

 By using the [HttpTreemux](github.com/dimfeld/httptreemux) library which provides a better router than the default `net/http` router, we register for routes which sets up our applications process as a
 whole.

 - We set up the home("/") which loads up our application's main page where all inboxes are displayed and accessed.
```go
	mux.GET("/", inbox.GetAllInbox)
```

- We set up the route which when it recieves a `POST` requests, generates a inbox with a randomly giving `ID` and redirects to
the URL for the inbox to which requests could be made to save http
requests to.
```go
	mux.POST("/inbox", inbox.NewInbox)
```

- We also set up a route which receives requests for inbox with their unique `ID`s.
```go
	mux.GET("/inbox/:id", inbox.GetInbox)
```

- And lastly setup a route through which requests already saved can be used.
```go
	mux.GET("/inbox/:id/:reqid", inbox.GetInboxItem)
```
