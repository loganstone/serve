package dir

import (
	"os"
	"testing"

	"gotest.tools/assert"
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
	assert.Assert(t, err == nil)
}
