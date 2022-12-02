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

func TestRegisterVersionHandler(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "")
	if err != nil {
		t.Fatalf("error creating temp directory: %s", err.Error())
	}
	defer os.RemoveAll(tmpDir)

	router := mux.NewRouter()
	p := promoter.NewPromoter(tmpDir)
	err = p.ReadFiles()
	if err != nil {
		t.Fatalf("can not read promoted files: %s", err)
	}

	RegisterVersionHandler(router, p)
	testServer := httptest.NewServer(router)
	defer testServer.Close()

	resp, err := http.Get(testServer.URL + "/files/version")
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

	var response map[string]string
	err = json.Unmarshal(actual, &response)
	if err != nil {
		t.Errorf("json unmarshal error: %s", err)
	}

	v := response["version"]
	if v != p.Version() {
		t.Errorf("invalid response, want: %s, get: %s", p.Version(), v)
	}
}
