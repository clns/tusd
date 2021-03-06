package tusd

import (
	"net/http"
	"os"
	"testing"
)

type headStore struct {
	zeroStore
}

func (s headStore) GetInfo(id string) (FileInfo, error) {
	if id != "yes" {
		return FileInfo{}, os.ErrNotExist
	}

	return FileInfo{
		Offset: 11,
		Size:   44,
		MetaData: map[string]string{
			"name": "lunrjs.png",
			"type": "image/png",
		},
	}, nil
}

func TestHead(t *testing.T) {
	handler, _ := NewHandler(Config{
		BasePath:  "https://buy.art/",
		DataStore: headStore{},
	})

	(&httpTest{
		Name:   "Successful request",
		Method: "HEAD",
		URL:    "yes",
		ReqHeader: map[string]string{
			"TUS-Resumable": "1.0.0",
		},
		Code: http.StatusNoContent,
		ResHeader: map[string]string{
			"Offset":        "11",
			"Entity-Length": "44",
			"Metadata":      "name bHVucmpzLnBuZw==,type aW1hZ2UvcG5n",
		},
	}).Run(handler, t)

	(&httpTest{
		Name:   "Non-existing file",
		Method: "HEAD",
		URL:    "no",
		ReqHeader: map[string]string{
			"TUS-Resumable": "1.0.0",
		},
		Code: http.StatusNotFound,
	}).Run(handler, t)
}
