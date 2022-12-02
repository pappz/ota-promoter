package api

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"strconv"
	"testing"

	log "github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
	"github.com/pappz/ota-promoter/promoter"
)

func TestMain(m *testing.M) {
	log.StandardLogger().SetLevel(log.FatalLevel)
	os.Exit(m.Run())
}

func TestRegisterDownloadHandler_invalidFile(t *testing.T) {
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

	RegisterDownloadHandler(router, p)
	testServer := httptest.NewServer(router)
	defer testServer.Close()

	resp, err := http.Get(testServer.URL + "/files/invalidchecksum")
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 400 {
		t.Errorf("invalid response code, want: %d, get: %d", 400, resp.StatusCode)
	}
}

func TestRegisterDownloadHandler(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "")
	if err != nil {
		t.Fatalf("error creating temp directory: %s", err.Error())
	}
	defer os.RemoveAll(tmpDir)
	tf, err := os.CreateTemp(tmpDir, "example.txt")
	if err != nil {
		t.Fatalf("error creating temp directory: %s", err.Error())
	}

	router := mux.NewRouter()
	p := promoter.NewPromoter(tmpDir)
	err = p.ReadFiles()
	if err != nil {
		t.Fatalf("can not read promoted files: %s", err)
	}

	RegisterDownloadHandler(router, p)
	testServer := httptest.NewServer(router)
	defer testServer.Close()

	resp, err := http.Get(testServer.URL + "/files/" + p.PromotedFiles()[0].Checksum)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("received non-200 response: %d", resp.StatusCode)
	}
	_, err = io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	if resp.Header.Get("X-target-path") != path.Base(tf.Name()) {
		t.Errorf("invalid target-path, want: %s, get: %s", tf.Name(), resp.Header.Get("X-target-path"))
	}

	size, err := fileSize(tf)
	if err != nil {
		t.Fatalf("faild to determinen file size: %s", err)
	}

	if resp.Header.Get("Content-Length") != strconv.FormatInt(size, 10) {
		t.Errorf("invalid target-path, want: %s, get: %s", tf.Name(), resp.Header.Get("X-target-path"))
	}

	if resp.Header.Get("Content-Disposition") != wantContentDisposition(tf) {
		t.Errorf("invalid Content-Disposition, want: %s, get: %s", wantContentDisposition(tf), resp.Header.Get("Content-Disposition"))
	}
}

func fileSize(f *os.File) (int64, error) {
	fi, err := f.Stat()
	if err != nil {
		return 0, err
	}
	return fi.Size(), nil
}

func wantContentDisposition(f *os.File) string {
	return fmt.Sprintf("attachment; filename=%s", path.Base(f.Name()))
}
