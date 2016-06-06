package app_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/influx6/httpinbox/app"
)

var (
	inbox   string
	success = "\u2713"
	failed  = "\u2717"
	client  = http.Client{Timeout: 30 * time.Second}
)

// TestAPI validates thee behaviour of the API.
func TestAPI(t *testing.T) {
	t.Logf("Given the need to use the HttpInbox API")
	{
		t.Logf("\tWhen given an HttpInbox server")
		{

			// Delete the data directory used for the test when done.
			defer func() {
				os.RemoveAll("./data")
			}()

			inbox := app.New("./data", "./views")
			serv := httptest.NewServer(inbox)
			defer serv.Close()

			testIndexRoute(serv, t)
			testAddInbox(serv, t)
			testAddToInbox(serv, t)
			testGetInbox(serv, t)
			testGetInboxItem(serv, t)
			testGetInboxItemData(serv, t)
		}
	}
}

func testIndexRoute(server *httptest.Server, t *testing.T) {
	t.Logf("\tWhen requests for the inbox list '/'")
	{

		res, err := client.Get(server.URL + "/")
		if err != nil {
			t.Fatalf("\t%s\tShould have successfully completed http request: %s", failed, err)
		}
		t.Logf("\t%s\tShould have successfully completed http request", success)

		if res.StatusCode != 200 {
			t.Fatalf("\t%s\tShould have successfully received a 200 http status code", failed)
		}
		t.Logf("\t%s\tShould have successfully received a 200 http status code", success)

	}
}

func testAddInbox(server *httptest.Server, t *testing.T) {
	t.Logf("\tWhen adding a new inbox with a 'POST /inbox'")
	{
		res, err := client.Post(server.URL+"/inbox", "text/plain", nil)
		if err != nil {
			t.Fatalf("\t%s\tShould have successfully completed http request: %s", failed, err)
		}
		t.Logf("\t%s\tShould have successfully completed http request", success)

		if res.StatusCode != 301 {
			t.Fatalf("\t%s\tShould have successfully received a 301 http status code", failed)
		}
		t.Logf("\t%s\tShould have successfully received a 301 http status code", success)

		loc, err := res.Location()
		if err != nil {
			t.Fatalf("\t%s\tShould have successfully recieved new location for http request: %s", failed, err)
		}
		t.Logf("\t%s\tShould have successfully recieved new location for http request", success)

		if loc.Path == "/inbox" {
			t.Logf("\t\tRecieved Location: %s\n", loc.Path)
			t.Fatalf("\t%s\tShould have successfully recieved new location different from post", failed)
		}
		t.Logf("\t%s\tShould have successfully recieved new location different from post", success)

		inbox = loc.Path
	}
}

func testAddToInbox(server *httptest.Server, t *testing.T) {
	t.Logf("\tWhen adding a requests to a inbox with 'ANY /inbox/:id'")
	{
		res, err := client.Head(server.URL + inbox)
		if err != nil {
			t.Fatalf("\t%s\tShould have successfully completed http request: %s", failed, err)
		}
		t.Logf("\t%s\tShould have successfully completed http request", success)

		if res.StatusCode != 200 {
			t.Fatalf("\t%s\tShould have successfully received a 200 http status code for %s", failed, inbox)
		}
		t.Logf("\t%s\tShould have successfully received a 200 http status code for %s", success, inbox)

	}
}

func testGetInbox(server *httptest.Server, t *testing.T) {
	t.Logf("\tWhen getting a inbox with its requests 'GET /inbox/:id'")
	{

		res, err := client.Get(server.URL + inbox)
		if err != nil {
			t.Fatalf("\t%s\tShould have successfully completed http request: %s", failed, err)
		}
		t.Logf("\t%s\tShould have successfully completed http request", success)

		if res.StatusCode != 200 {
			t.Fatalf("\t%s\tShould have successfully received a 200 http status code for %s", failed, inbox)
		}
		t.Logf("\t%s\tShould have successfully received a 200 http status code for %s", success, inbox)

	}
}

func testGetInboxItem(server *httptest.Server, t *testing.T) {
	t.Logf("\tWhen getting a inbox requests details with 'GET /inbox/:id/:reqid'")
	{
		inboxURL := inbox + "/0"
		res, err := client.Get(server.URL + inboxURL)
		if err != nil {
			t.Fatalf("\t%s\tShould have successfully completed http request: %s", failed, err)
		}
		t.Logf("\t%s\tShould have successfully completed http request", success)

		if res.StatusCode != 200 {
			t.Fatalf("\t%s\tShould have successfully received a 200 http status code for %s", failed, inboxURL)
		}
		t.Logf("\t%s\tShould have successfully received a 200 http status code for %s", success, inboxURL)

		content := res.Header.Get("Content-Type")

		if content != "text/html" {
			t.Fatalf("\t%s\tShould have successfully recieved response with Content-Type[%q]", failed, "text/html")
		}
		t.Logf("\t%s\tShould have successfully recieved response with Content-Type[%q]", success, "text/html")
	}
}

func testGetInboxItemData(server *httptest.Server, t *testing.T) {
	t.Logf("\tWhen getting a inbox requests data with 'GET /inbox/:id/:reqid'")
	{
		inboxURL := inbox + "/0"
		req, err := http.NewRequest("GET", server.URL+inboxURL, nil)
		if err != nil {
			t.Fatalf("\t%s\tShould have successfully created http request: %s", failed, err)
		}
		t.Logf("\t%s\tShould have successfully created http request", success)

		req.Header.Set("Accepts", "application/data")

		res, err := client.Do(req)
		if err != nil {
			t.Fatalf("\t%s\tShould have successfully completed http request: %s", failed, err)
		}
		t.Logf("\t%s\tShould have successfully completed http request", success)

		content := res.Header.Get("Content-Type")

		if content != "application/data" {
			t.Fatalf("\t%s\tShould have successfully recieved response with Content-Type[%q]", failed, "application/data")
		}
		t.Logf("\t%s\tShould have successfully recieved response with Content-Type[%q]", success, "application/data")

		if res.StatusCode != 200 {
			t.Fatalf("\t%s\tShould have successfully received a 200 http status code for %s", failed, inboxURL)
		}
		t.Logf("\t%s\tShould have successfully received a 200 http status code for %s", success, inboxURL)
	}
}
