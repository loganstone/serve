package dir

import (
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
