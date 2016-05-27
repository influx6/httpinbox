# Setting up Testing
Ensuring the operating status of our application is most important and these means
we need to have tests we can run every time a new update is added to our repository
to ensure the functionality of our application.

To achieve these we need to have at leasts unit tests ready to be runned by our
deployment scripts either at commits to our repositories or before deployments
of new updates.

We will be adding basic level controller tests which allows us to validate the
behaviour of our endpoints and underline structure of our apps  with the code below:

```go

package app_test

import (
	"net/http/httptest"
	"testing"

	"github.com/dimfeld/httptreemux"
	"github.com/influx6/httpinbox/app"
)

// TestAPI validates thee behaviour of the API.
func TestAPI(t *testing.T) {
	t.Logf("Given the need to use the HttpInbox API")
	{
		t.Logf("\tWhen given an HttpInbox server")
		{

			inbox := app.New("../data", "./views")

			mux := httptreemux.New()
			mux.GET("/", inbox.GetAllInbox)
			mux.POST("/inbox", inbox.NewInbox)
			mux.GET("/inbox/:id", inbox.GetInbox)
			mux.GET("/inbox/:id/:reqid", inbox.GetInboxItem)

			for _, method := range []string{"POST", "DELETE", "PUT", "PATCH", "HEAD"} {
				mux.Handle(method, "/inbox/:id", inbox.AddToInbox)
			}

			serv := httptest.NewServer(mux)
			defer serv.Close()

			testIndexRoute(serv, t)
			testAddInbox(serv, t)
			testAddToInbox(serv, t)
			testGetInbox(serv, t)
			testGetInboxItem(serv, t)
		}
	}
}

func testIndexRoute(server *httptest.Server, t *testing.T) {
	t.Logf("\tWhen requests for the inbox list '/'")
	{
	}
}

func testAddInbox(server *httptest.Server, t *testing.T) {
	t.Logf("\tWhen adding a new inbox with a 'POST /inbox'")
	{
	}
}

func testAddToInbox(server *httptest.Server, t *testing.T) {
	t.Logf("\tWhen adding a requests to a inbox with 'ANY /inbox/:id'")
	{
	}
}

func testGetInbox(server *httptest.Server, t *testing.T) {
	t.Logf("\tWhen getting a inbox with its requests 'GET /inbox/:id'")
	{
	}
}

func testGetInboxItem(server *httptest.Server, t *testing.T) {
	t.Logf("\tWhen getting a inbox requests details with 'GET /inbox/:id/:reqid'")
	{
	}
}

```

These test will allow us to us to ensure functional correctness and with the
below scripts added to our Harrow deployment system we can ensure these test
is runned using `go test` everytime we commit to our git repository.

```sh


```
---Pictures demostrating the addition of test script to arrow here----------
