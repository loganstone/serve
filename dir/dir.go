package dir

import "path/filepath"

// Abs ...
func Abs(dir string) (string, error) {
	if filepath.IsAbs(dir) {
		return dir, nil
	}
	return filepath.Abs(dir)
}
