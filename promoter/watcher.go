package promoter

import (
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

// ChangeWatcher continuously visit the changes in the directory and sub path. In case of changes call the callback fn
type ChangeWatcher struct {
	watcher        *fsnotify.Watcher
	promotedFolder string
	onChangedFn    func()
	onErrFn        func(error)
}

// NewChangeWatcher create new ChangeWatcher instance
func NewChangeWatcher(promotedFolder string, onChanged func(), onError func(error)) (ChangeWatcher, error) {
	var err error
	cw := ChangeWatcher{
		promotedFolder: promotedFolder,
		onChangedFn:    onChanged,
		onErrFn:        onError,
	}
	cw.watcher, err = fsnotify.NewWatcher()
	return cw, err
}

// Watch start the watch service
func (cw *ChangeWatcher) Watch() (err error) {
	if err := filepath.Walk(cw.promotedFolder, cw.watchDir); err != nil {
		return err
	}
	go func() {
		for {
			select {
			case _ = <-cw.watcher.Events:
				cw.onChangedFn()
			case err := <-cw.watcher.Errors:
				cw.onErrFn(err)
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
