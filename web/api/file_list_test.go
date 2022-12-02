package api

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gorilla/mux"
	"github.com/pappz/ota-promoter/promoter"
)

func TestRegisterFileListHandler(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "")
	if err != nil {
		t.Fatalf("error creating temp directory: %s", err.Error())
	}
	defer os.RemoveAll(tmpDir)
	_, err = os.CreateTemp(tmpDir, "example.txt")
	if err != nil {
		t.Fatalf("error creating temp directory: %s", err.Error())
	}

	router := mux.NewRouter()
	p := promoter.NewPromoter(tmpDir)
	err = p.ReadFiles()
	if err != nil {
		t.Fatalf("can not read promoted files: %s", err)
	}

	RegisterFileListHandler(router, p)
	testServer := httptest.NewServer(router)
	defer testServer.Close()

	resp, err := http.Get(testServer.URL + "/files")
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("received non-200 response: %d\n", resp.StatusCode)
	}
	actual, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	var response ResponseFileList
	err = json.Unmarshal(actual, &response)
	if err != nil {
		t.Errorf("json unmarshal error: %s", err)
	}

	if response.Version != p.Version() {
		t.Errorf("invalid response, want: %s, get: %s", p.Version(), response.Version)
	}

	if len(response.Files) != 1 {
		t.Fatalf("invalid response")
	}

	if response.Files[0].Checksum != p.PromotedFiles()[0].Checksum {
		t.Errorf("invalid response, want: %s, get: %s", p.PromotedFiles()[0].Checksum, response.Files[0].Checksum)
	}
}
