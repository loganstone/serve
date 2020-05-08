package dir

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAbs(t *testing.T) {
	absPath, err := Abs(".")
	assert.NoError(t, err)

	expected, err := os.Getwd()
	assert.NoError(t, err)
	assert.Equal(t, absPath, expected)

	testPath := "/abs/path"
	_, err = Abs(testPath)
	assert.NoError(t, err)
}

func TestNewWatcher(t *testing.T) {
	_, err := NewWatcher(".")
	assert.NoError(t, err)

	_, err = NewWatcher("bad_dir")
	assert.Error(t, errors.New("lstat bad_dir: no such file or directory"), err)

	_, err = NewWatcher("./dir.go")
	assert.Error(t, errors.New("-d option value must be directory"), err)
}
