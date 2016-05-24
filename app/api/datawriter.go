package api

import (
	"net/http"
)

//==============================================================================

// WriteRequest defines a requests for storing a http.Requests by the
// data manager.
type WriteRequest struct {
	ID   string
	rq   *http.Request
	Done chan error
}

// NewWriteRequest returns a new instance of a request with the
// giving ID and reqest object.
func NewWriteRequest(id string, rq *http.Request) *WriteRequest {
	wr := WriteRequest{ID: id, rq: rq, Done: make(chan error)}
	return &wr
}

// ToJSON returns the binary representation for the requests to be
// written to file. It just encapsulates the transformations for us.
func (w *WriteRequest) ToJSON() ([]byte, error) {
	return nil, nil
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

// begin contains the write management routine for the storage of requests
// into a proper file.
func (dm *DataMan) begin() {
	for {
		select {
		case req, ok := <-dm.newWrites:
			if !ok {
				break
			}

			// datafile := fmt.Sprintf("%s/%s/%d",dm.dataDir,req.ID)
			// file, err := os.Open(datafile)
			_ = req
		}
	}
}

//==============================================================================
