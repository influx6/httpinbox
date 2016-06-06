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

```

We can run our tests during development easily by calling `go test -race -cpu 10`.
The value we give to `-cpu` in truth won't affect much as only the physically
available CPU's will be used, and it's always adviced to run your tests and even
execution of app during development with the race flag `-race` to help catch any
race conditions lying about.

These tests will allow us to us to ensure functional correctness and with the
below scripts added to our Harrow deployment system we can ensure these test
is runned using `go test` everytime we commit to our git repository.

```sh


```
---Pictures demostrating the addition of test script to arrow here----------
