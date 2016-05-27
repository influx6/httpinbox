# Building your own HttpInbox
This article is about showing how you can build your own application using Go and
take it from completion to deployment using Harrow to simplify your life,
this allows you a insight into how you can combine different technology with Harrow
to seamlessly move from product iteration to product continuous deployment and
integration.

We will be building a simple application that logs http requests coming
to it by allocating `inboxes` with URLs that allows you capture and store which
requests are hitting your servers, this example will be very contrived but with
this as a base more complex versions can indeed be built.

The concept is simple: A service that allows you to create a http request inbox
where clients could hit and we can capture that requests to allow us perform
whatever metric or process we wish to do so with that data.

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

We have a `app` directory where the controller code for our application with their
view templates in `app/views` and also a `main.go` file which will be used to both
run our application and build our binary for the project when deploying.
In accordance with the new go1.6 approach, we also have a `vendor` directory
where all source code based dependencies are stored, this allows us to lock down
our dependencies with other libraries.

Our application `main.go` file is rather simple, it loads off specific configuration
from the environment which allows us to set different configurations depending on
the platform be it for `development-testing` or `production` without touching the
codebase.

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

 I personally believe that the use of panics should be done sparingly and only
 during the initial, start up phase of an application, because using it during
 the runtime phase simplify kills your application, loosing a lot of context,
 data and providing adequate headaches for future runs.

```go
	inbox := app.New(dataDir, viewsDir)
```

	The core of our small application which instance we created as above is found
	within the `app` folder which contains the controller code that accepts requests
	processing them accordingly.

	It's made up of two structures which encapsulate the intended behaviour and
	response we need. The first is the `HttpInbox` struct and the second the `DataMan`
	struct which handles reads and writes requests for stored inboxes.

	- HttpInbox  **/app/api.go**
	The `HttpInbox` implements the logic for handling the different route behaviors
	we need, our app has 3 basic functions:

	1. To create a new inbox when receiving a `POST /inbox`, these lets us
	generate new inboxes and redirects the client to a new URL `/inbox/:id`, which
	requests of non-GET methods could be made to.

	2. To capture all incoming requests using the `NON-GET /inbox/:id`
		URL(mostly all Non-Get: HEAD, DELETE, POST, PATCH,..etc), and store their
		internal details in a persistent file within the appropriate directory,
		catalogue by a increasing index number to allow viewing of a specific requests
		using the `/inbox/:id/:index` URL.

	3. To allow viewing all inboxes with the `GET /inbox/:id` and their
	associated requests using go templates by allowing use of

	Underneath `HttpInbox` has a instance of the `DataMan` which encapsulates the
	IO operation which are needed by our app. This allows us a sweet separation
	between our controller code and the service part which handles the low-level
	logic needed by our application.

	```go
package app

import (
	"bytes"
	"fmt"
	"html/template"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

type viewDir string

// For returns a new path joined with the giving path.
func (v viewDir) For(path string) string {
	return filepath.Join(string(v), path)
}

// HTTPInbox defines a struct which holds all controller methods for the HTTPInbox API.
type HTTPInbox struct {
	man   *DataMan
	mbl   sync.RWMutex
	inbox map[string]int
	views viewDir
}

// New returns a new API instance.
func New(dataDir string, views string) *HTTPInbox {
	api := HTTPInbox{
		views: viewDir(views),
		man:   NewDataMan(dataDir),
		inbox: make(map[string]int),
	}

	var new bool
	files, err := api.man.ReadAllInbox()
	if err != nil {
		os.MkdirAll(dataDir, 0755)
		new = true
	}

	if !new {
		api.mbl.Lock()
		for _, file := range files {
			if strings.HasPrefix(file.Name(), ".") {
				continue
			}

			fmt.Printf("Adding Inbox: %s\n", file.Name())
			api.inbox[file.Name()] = 0
			if fsl, err := api.man.ReadInbox(file.Name()); err == nil {
				fmt.Printf("Adding Inbox Item: Box[%s] : Items[%d]\n", file.Name(), len(fsl))
				api.inbox[file.Name()] = len(fsl)
			}
		}
		api.mbl.Unlock()
	}

	return &api
}

// NewInbox handles the creation of a new inbox for the reception of http requests.
func (h *HTTPInbox) NewInbox(res http.ResponseWriter, req *http.Request, param map[string]string) {
	uuid := randString(10)

	h.mbl.RLock()
	if _, ok := h.inbox[uuid]; ok {
		uuid = randString(10)
	}
	h.mbl.RUnlock()

	if err := h.man.PrepareInbox(uuid); err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte(err.Error()))
		return
	}

	h.mbl.Lock()
	defer h.mbl.Unlock()
	h.inbox[uuid] = 0

	http.Redirect(res, req, fmt.Sprintf("/inbox/%s", uuid), 301)
}

// AddToInbox adds the needed requests into the inbox lists of requests.
func (h *HTTPInbox) AddToInbox(res http.ResponseWriter, req *http.Request, param map[string]string) {
	inboxID := param["id"]

	var ok bool
	var count int

	h.mbl.RLock()
	count, ok = h.inbox[inboxID]
	h.mbl.RUnlock()

	if !ok {
		res.WriteHeader(http.StatusNotFound)
		return
	}

	if err := h.man.WriteInbox(inboxID, req, count); err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte(err.Error()))
		return
	}

	// Increment the received count.
	count++

	h.mbl.Lock()
	defer h.mbl.Unlock()
	h.inbox[inboxID] = count
}

// GetAllInbox retrieves all inbox in the system.
func (h *HTTPInbox) GetAllInbox(res http.ResponseWriter, req *http.Request, param map[string]string) {
	tm, err := template.ParseFiles(h.views.For("layout.tml"), h.views.For("all.tml"))
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte(err.Error()))
		return
	}

	box := make(map[string]int)

	h.mbl.RLock()
	// fmt.Printf("Allboxes: %+s\n", h.inbox)
	for id, c := range h.inbox {
		box[id] = c
	}
	h.mbl.RUnlock()

	var buf bytes.Buffer

	data := struct {
		Total int
		Items map[string]int
	}{Total: len(box), Items: box}

	if err := tm.ExecuteTemplate(&buf, "layout", data); err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte(err.Error()))
		return
	}

	res.WriteHeader(http.StatusOK)
	res.Write(buf.Bytes())
}

// GetInbox retrieves a giving box using the id it recieves.
func (h *HTTPInbox) GetInbox(res http.ResponseWriter, req *http.Request, param map[string]string) {
	inboxID := param["id"]

	var ok bool
	h.mbl.RLock()
	_, ok = h.inbox[inboxID]
	h.mbl.RUnlock()

	if !ok {
		res.WriteHeader(http.StatusNotFound)
		res.Write([]byte("InboxID['" + inboxID + "'] not Found."))
		return
	}

	tm, err := template.ParseFiles(h.views.For("layout.tml"), h.views.For("list.tml"))
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte(err.Error()))
		return
	}

	files, err := h.man.ReadInbox(inboxID)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte(err.Error()))
		return
	}

	var buf bytes.Buffer

	data := struct {
		Inbox string
		Items []os.FileInfo
	}{Inbox: inboxID, Items: files}

	if err := tm.ExecuteTemplate(&buf, "layout", data); err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte(err.Error()))
		return
	}

	res.WriteHeader(http.StatusOK)
	res.Write(buf.Bytes())
}

// GetInboxItem retrieves a giving box using the id it recieves.
func (h *HTTPInbox) GetInboxItem(res http.ResponseWriter, req *http.Request, param map[string]string) {
	inboxID := param["id"]

	var ok bool
	h.mbl.RLock()
	_, ok = h.inbox[inboxID]
	h.mbl.RUnlock()

	if !ok {
		res.WriteHeader(http.StatusNotFound)
		return
	}

	itemID, err := strconv.Atoi(param["reqid"])
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte("ItemID must be an int: " + err.Error()))
		return
	}

	tm, err := template.ParseFiles(h.views.For("layout.tml"), h.views.For("single.tml"))
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte(err.Error()))
		return
	}

	data, err := h.man.ReadInboxItem(inboxID, itemID)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte(err.Error()))
		return
	}

	var buf bytes.Buffer

	mdata := struct {
		Inbox string
		Item  int
		Data  string
	}{Inbox: inboxID, Item: itemID, Data: string(data)}

	if err := tm.ExecuteTemplate(&buf, "layout", mdata); err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte(err.Error()))
		return
	}

	res.WriteHeader(http.StatusOK)
	res.Write(buf.Bytes())
}

// DestroyInbox handles the destruction of inbox with all its contents.
func (h *HTTPInbox) DestroyInbox(res http.ResponseWriter, req *http.Request, param map[string]string) {
}

//==============================================================================

// RandString generates a set of random numbers of a set length
func randString(n int) string {
	const alphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	var bytes = make([]byte, n)
	rand.Read(bytes)
	for i, b := range bytes {
		bytes[i] = alphanum[b%byte(len(alphanum))]
	}
	return string(bytes)
}

	```

	- DataWriter

```go
package app

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"os"
	"sort"
)

//==============================================================================

// WriteRequest defines a requests for storing a http.Requests by the
// data manager.
type WriteRequest struct {
	ID     string
	rindex int
	rq     *http.Request
	Done   chan error
}

// NewWriteRequest returns a new instance of a request with the
// giving ID and reqest object.
func NewWriteRequest(id string, rq *http.Request, rc int) *WriteRequest {
	wr := WriteRequest{ID: id, rq: rq, rindex: rc, Done: make(chan error)}
	return &wr
}

// ToJSON returns the binary representation for the requests to be
// written to file. It just encapsulates the transformations for us.
func (w *WriteRequest) ToJSON() ([]byte, error) {
	return httputil.DumpRequest(w.rq, true)
}

//==============================================================================

// DataMan defines a data write and reads manager which allows
// central control of the incoming write requests for saving specific
// data.
type DataMan struct {
	dataDir   string
	newWrites chan *WriteRequest
}

// NewDataMan returns a new instance of a DataMan.
func NewDataMan(dir string) *DataMan {
	dm := DataMan{
		dataDir:   dir,
		newWrites: make(chan *WriteRequest),
	}

	go dm.begin()

	return &dm
}

// ByName Implement sort.Interface for []os.FileInfo based on Name()
type ByName []os.FileInfo

func (v ByName) Len() int           { return len(v) }
func (v ByName) Swap(i, j int)      { v[i], v[j] = v[j], v[i] }
func (v ByName) Less(i, j int) bool { return v[i].Name() < v[j].Name() }

// ReadAllInbox gets a inbox and a specific item from that inbox.
func (dm *DataMan) ReadAllInbox() ([]os.FileInfo, error) {
	dir, err := os.Open(dm.dataDir)
	if err != nil {
		return nil, err
	}

	// Read only one level deep.
	files, err := dir.Readdir(0)
	if err != nil {
		return nil, err
	}

	// Remove any non-directories or unwanted files.
	count := len(files) - 1

	for ind, info := range files {
		if info.Name() == ".DS_Store" || !info.IsDir() {
			files[ind] = files[count]
			count--
			continue
		}
	}

	files = files[:count+1]

	// Sort in ascending order by name.
	sort.Sort(ByName(files))

	return files, nil
}

// PrepareInbox prepares the inbox folder for storage and allows us
// to persist that inbox.
func (dm *DataMan) PrepareInbox(inboxID string) error {
	dirfile := fmt.Sprintf("%s/%s/", dm.dataDir, inboxID)

	if err := os.MkdirAll(dirfile, 0755); err != nil {
		return err
	}

	return nil
}

// ReadInbox gets a inbox and all items from that inbox.
func (dm *DataMan) ReadInbox(inboxID string) ([]os.FileInfo, error) {
	datafile := fmt.Sprintf("%s/%s/", dm.dataDir, inboxID)
	dir, err := os.Open(datafile)
	if err != nil {
		return nil, err
	}

	// Read only one level deep.
	files, err := dir.Readdir(0)
	if err != nil {
		return nil, err
	}

	// Remove any directories or unwanted files.
	count := len(files) - 1

	for ind, info := range files {
		if info.Name() == ".DS_Store" || info.IsDir() {
			files[ind] = files[count]
			count--
			continue
		}
	}

	files = files[:count+1]

	// Sort in ascending order by name.
	sort.Sort(ByName(files))

	return files, nil
}

// ReadInboxItem gets a inbox and a specific item from that inbox.
func (dm *DataMan) ReadInboxItem(inboxID string, rc int) ([]byte, error) {
	datafile := fmt.Sprintf("%s/%s/%d", dm.dataDir, inboxID, rc)
	file, err := os.Open(datafile)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	// Generally its not adviced to do this especially if your
	// files are very large but rather use a io.Reader when
	// larger files are coming.
	return ioutil.ReadAll(file)
}

// WriteInbox sends a write requests for a specific inbox to be written into
// the inbox store.
func (dm *DataMan) WriteInbox(inboxID string, req *http.Request, rc int) error {
	wq := NewWriteRequest(inboxID, req, rc)
	dm.newWrites <- wq
	return <-wq.Done
}

// begin contains the write management routine for the storage of requests
// into a proper file.
func (dm *DataMan) begin() {
	for {
		select {
		case wq, ok := <-dm.newWrites:
			if !ok {
				break
			}

			data, err := wq.ToJSON()
			if err != nil {
				wq.Done <- err
				continue
			}

			dirfile := fmt.Sprintf("%s/%s", dm.dataDir, wq.ID)

			if err = os.MkdirAll(dirfile, 0755); err != nil {
				wq.Done <- err
				continue
			}

			// Create the needed file within the appropriate path.
			datafile := fmt.Sprintf("%s/%d", dirfile, wq.rindex)
			file, err := os.Create(datafile)
			if err != nil {
				wq.Done <- err
				continue
			}

			// Attempt to write the request into our file.
			total, err := file.Write(data)
			if err != nil {
				wq.Done <- err
				file.Close()
				continue
			}

			// If this was not a complete write then close and
			// delete te file, we don't want half written mess.
			if total != len(data) {
				wq.Done <- errors.New("Invalid Write")
				file.Close()
				os.Remove(datafile)
				continue
			}

			file.Close()
			close(wq.Done)
		}
	}
}

//==============================================================================
```
