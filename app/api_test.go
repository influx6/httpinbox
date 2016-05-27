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
