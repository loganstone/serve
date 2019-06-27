package dir

import (
	"log"
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

// NewWatcher .
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

// NowMyWatchBegins ...
func NowMyWatchBegins(dir string, w *fsnotify.Watcher) {
	log.Printf("verified directory [%s], and now my watch begins", dir)
	for {
		select {
		case event := <-w.Events:
			log.Printf("Watcher - %s\n", event)

		case err := <-w.Errors:
			log.Println("Watcher ERROR!", err)
		}
	}
}
