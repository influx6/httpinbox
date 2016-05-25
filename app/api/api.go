package api

import (
	"bytes"
	"fmt"
	"html/template"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
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
		os.MkdirAll(dataDir, 0777)
		new = true
	}

	if !new {
		api.mbl.Lock()
		for _, file := range files {
			if fsl, err := api.man.ReadInbox(file.Name()); err == nil {
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

	h.mbl.Lock()
	defer h.mbl.Unlock()
	h.inbox[uuid] = 0

	http.Redirect(res, req, fmt.Sprintf("/inbox/%s", uuid), 301)
}

// AddToInbox adds the needed requests into the inbox lists of requests.
func (h *HTTPInbox) AddToInbox(res http.ResponseWriter, req *http.Request, param map[string]string) {
	inboxID := param[":id"]

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
	for id, c := range h.inbox {
		box[id] = c
	}
	h.mbl.RUnlock()

	var buf bytes.Buffer

	if err := tm.ExecuteTemplate(&buf, "layout", map[string]interface{}{
		"total": len(box),
		"items": box,
	}); err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte(err.Error()))
		return
	}

	res.WriteHeader(http.StatusOK)
	res.Write(buf.Bytes())
}

// GetInbox retrieves a giving box using the id it recieves.
func (h *HTTPInbox) GetInbox(res http.ResponseWriter, req *http.Request, param map[string]string) {
	inboxID := param[":id"]

	var ok bool
	h.mbl.RLock()
	_, ok = h.inbox[inboxID]
	h.mbl.RUnlock()

	if !ok {
		res.WriteHeader(http.StatusNotFound)
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

	if err := tm.ExecuteTemplate(&buf, "layout", map[string]interface{}{
		"inbox": inboxID,
		"items": files,
	}); err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte(err.Error()))
		return
	}

	res.WriteHeader(http.StatusOK)
	res.Write(buf.Bytes())
}

// GetInboxItem retrieves a giving box using the id it recieves.
func (h *HTTPInbox) GetInboxItem(res http.ResponseWriter, req *http.Request, param map[string]string) {
	inboxID := param[":id"]

	var ok bool
	h.mbl.RLock()
	_, ok = h.inbox[inboxID]
	h.mbl.RUnlock()

	if !ok {
		res.WriteHeader(http.StatusNotFound)
		return
	}

	itemID, err := strconv.Atoi(param[":reqid"])
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

	if err := tm.ExecuteTemplate(&buf, "layout", map[string]interface{}{
		"inbox": inboxID,
		"item":  itemID,
		"data":  string(data),
	}); err != nil {
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
