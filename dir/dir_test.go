package dir

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAbs(t *testing.T) {
	absPath, err := Abs(".")
	if err != nil {
		t.Fatal(err)
	}

	expected, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, absPath, expected)

	testPath := "/abs/path"
	_, err = Abs(testPath)
	assert.Nil(t, err)
}
