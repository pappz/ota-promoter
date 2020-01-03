package main

import (
	"github.com/fsnotify/fsnotify"
	"os"
	"path/filepath"
)

var (
	watcher *fsnotify.Watcher
)

func watch() (err error) {
	watcher, err = fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	if err := filepath.Walk(promotedFolder, watchDir); err != nil {
		_ = watcher.Close()
		return err
	}

	go func() {
		for {
			select {
			case _ = <-watcher.Events:
				if err := readFiles(); err != nil {
					log.Errorf("failed to read promoted files: %v", err)
				}
			case err := <-watcher.Errors:
				if err != nil {
					log.Errorf("Watcher error: %s", err.Error())
				}
			}
		}
	}()

	return
}

func watchDir(path string, fi os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	if fi.Mode().IsDir() {
		return watcher.Add(path)
	}

	return nil
}

func closeWatcher() {
	if watcher != nil {
		_ = watcher.Close()
	}
}
