package promoter

import (
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

type ChangeWatcher struct {
	watcher *fsnotify.Watcher
}

func NewChangeWatcher() (ChangeWatcher, error) {
	var err error
	cw := ChangeWatcher{}
	cw.watcher, err = fsnotify.NewWatcher()
	return cw, err
}

func (cw *ChangeWatcher) Watch(promotedFolder string, onChanged func(), onError func(error)) (err error) {
	if err := filepath.Walk(promotedFolder, cw.watchDir); err != nil {
		return err
	}
	go func() {
		for {
			select {
			case _ = <-cw.watcher.Events:
				onChanged()
			case err := <-cw.watcher.Errors:
				onError(err)
			}
		}
	}()

	return
}

func (cw *ChangeWatcher) CloseWatcher() {
	if cw.watcher == nil {
		return
	}
	_ = cw.watcher.Close()
}

func (cw *ChangeWatcher) watchDir(path string, fi os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	if fi.Mode().IsDir() {
		return cw.watcher.Add(path)
	}

	return nil
}
