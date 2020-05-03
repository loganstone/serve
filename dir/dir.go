package dir

import (
	"errors"
	"log"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

// Watcher .
type Watcher struct {
	*fsnotify.Watcher
	VerifiedDir string
}

// Abs .
func Abs(dir string) (string, error) {
	if filepath.IsAbs(dir) {
		return dir, nil
	}
	return filepath.Abs(dir)
}

// NewWatcher .
func NewWatcher(dir string) (*Watcher, error) {
	absPath, err := Abs(dir)
	if err != nil {
		return nil, err
	}

	dirInfo, err := os.Stat(absPath)
	if err != nil {
		return nil, err
	}

	if !dirInfo.IsDir() {
		return nil, errors.New("-d option value must be directory")
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	if err := watcher.Add(dir); err != nil {
		return nil, err
	}
	return &Watcher{watcher, dir}, nil
}

// NowMyWatchBegins .
func (w *Watcher) NowMyWatchBegins() {
	log.Printf("verified directory [%s], and now my watch begins", w.VerifiedDir)
	for {
		select {
		case event := <-w.Events:
			log.Printf("Watcher - %s\n", event)

		case err := <-w.Errors:
			log.Println("Watcher ERROR!", err)
		}
	}
}
