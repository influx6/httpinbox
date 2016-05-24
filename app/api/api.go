package api

import (
	"encoding/json"
	"net/http"
	"sync"
)

// HTTPInbox defines a struct which holds all controller methods for the HTTPInbox API.
type HTTPInbox struct {
	man   *DataMan
	mbl   sync.RWMutex
	inbox map[string]int
}

// New returns a new API instance.
func New(dataDir string) *HTTPInbox {
	api := HTTPInbox{
		man:   NewDataMan(dataDir),
		inbox: make(map[string]int),
	}

	return &api
}

// NewInbox handles the creation of a new inbox for the reception of http requests.
func (h *HTTPInbox) NewInbox(res http.ResponseWriter, req *http.Request, param map[string]string) {
	d, _ := json.Marshal(req)
	res.Write(d)
}

// AddToInbox adds the needed requests into the inbox lists of requests.
func (h *HTTPInbox) AddToInbox(res http.ResponseWriter, req *http.Request, param map[string]string) {
}

// GetInbox retrieves a giving box using the id it recieves.
func (h *HTTPInbox) GetInbox(res http.ResponseWriter, req *http.Request, param map[string]string) {
}

// DestroyInbox handles the destruction of inbox with all its contents.
func (h *HTTPInbox) DestroyInbox(res http.ResponseWriter, req *http.Request, param map[string]string) {
}
