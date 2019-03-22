package dir

import (
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

// Abs ...
func Abs(dir string) (string, error) {
	if filepath.IsAbs(dir) {
		return dir, nil
	}
	return filepath.Abs(dir)
}

// NewWatcher ...
func NewWatcher(dir string) (*fsnotify.Watcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	if err := watcher.Add(dir); err != nil {
		return nil, err
	}

	return watcher, nil
}
