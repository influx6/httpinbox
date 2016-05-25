package api

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

	// Sort in ascending order by name.
	sort.Sort(ByName(files))

	return files, nil
}

// ReadInbox gets a inbox and a specific item from that inbox.
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

			// Create the needed file within the appropriate path.
			datafile := fmt.Sprintf("%s/%s/%d", dm.dataDir, wq.ID, wq.rindex)
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
