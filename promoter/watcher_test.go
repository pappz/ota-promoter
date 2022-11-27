package promoter

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNewChangeWatcher_flat(t *testing.T) {
	onChangedChan := make(chan struct{})
	var onErrorResult error

	tmpDir, err := os.MkdirTemp("", "")
	if err != nil {
		t.Fatalf("error creating temp directory: %s", err.Error())
	}
	defer os.RemoveAll(tmpDir)

	onChangedFn := func() {
		onChangedChan <- struct{}{}
	}

	onErrorFn := func(err error) {
		onErrorResult = err
	}

	cw, err := NewChangeWatcher()
	if err != nil {
		t.Fatalf("%s", err.Error())
	}
	defer cw.CloseWatcher()

	err = cw.Watch(tmpDir, onChangedFn, onErrorFn)
	if err != nil {
		t.Fatalf("%s", err)
	}

	_, err = os.Create(tmpDir + "/empty.txt")
	if err != nil {
		t.Fatalf("failed to create sample file: %s", err)
	}

	timedOut := waitToChannel(onChangedChan)
	if timedOut {
		t.Errorf("changes not happend")
	}

	if onErrorResult != nil {
		t.Errorf("unexpected error result %s", onErrorResult)
	}
}

func TestNewChangeWatcher_subDir(t *testing.T) {
	subDirPath := "dir/sub"
	onChangedChan := make(chan struct{})
	var onErrorResult error

	tmpDir, err := prepareTestDirTree(subDirPath)
	if err != nil {
		t.Fatalf("failed to create tmpdir: %s", err)
	}
	defer os.RemoveAll(tmpDir)

	onChangedFn := func() {
		onChangedChan <- struct{}{}
	}

	onErrorFn := func(err error) {
		onErrorResult = err
	}

	cw, err := NewChangeWatcher()
	if err != nil {
		t.Fatalf("%s", err.Error())
	}
	defer cw.CloseWatcher()

	err = cw.Watch(tmpDir, onChangedFn, onErrorFn)
	if err != nil {
		t.Fatalf("%s", err)
	}

	_, err = os.Create(fmt.Sprintf("%s/%s/%s", tmpDir, subDirPath, "empty.txt"))
	if err != nil {
		t.Fatalf("failed to create sample file: %s", err)
	}

	timedOut := waitToChannel(onChangedChan)
	if timedOut {
		t.Errorf("changes not happend")
	}

	if onErrorResult != nil {
		t.Errorf("unexpected error result %s", onErrorResult)
	}
}

func TestNewChangeWatcher_deleteFile(t *testing.T) {
	tmpFileName := "empty.txt"
	onChangedChan := make(chan struct{})
	var onErrorResult error

	tmpDir, err := os.MkdirTemp("", "")
	if err != nil {
		t.Fatalf("error creating temp directory: %s", err.Error())
	}
	defer os.RemoveAll(tmpDir)

	_, err = os.Create(tmpDir + "/" + tmpFileName)
	if err != nil {
		t.Fatalf("failed to create sample file: %s", err)
	}

	onChangedFn := func() {
		onChangedChan <- struct{}{}
	}

	onErrorFn := func(err error) {
		onErrorResult = err
	}

	cw, err := NewChangeWatcher()
	if err != nil {
		t.Fatalf("%s", err.Error())
	}
	defer cw.CloseWatcher()

	err = cw.Watch(tmpDir, onChangedFn, onErrorFn)
	if err != nil {
		t.Fatalf("%s", err)
	}

	err = os.Remove(tmpDir + "/" + tmpFileName)
	if err != nil {
		t.Fatalf("failed to remove tmp file: %s", err)
	}

	timedOut := waitToChannel(onChangedChan)
	if timedOut {
		t.Errorf("changes not happend")
	}

	if onErrorResult != nil {
		t.Errorf("unexpected error result %s", onErrorResult)
	}
}

func waitToChannel(changedChan chan struct{}) bool {
	select {
	case <-changedChan:
		return false
	case <-time.After(1 * time.Second):
		return true
	}
}

func prepareTestDirTree(tree string) (string, error) {
	tmpDir, err := os.MkdirTemp("", "")
	if err != nil {
		return "", fmt.Errorf("error creating temp directory: %v\n", err)
	}

	err = os.MkdirAll(filepath.Join(tmpDir, tree), 0755)
	if err != nil {
		_ = os.RemoveAll(tmpDir)
		return "", err
	}

	return tmpDir, nil
}
