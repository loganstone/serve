package main

import (
	"os"
	"testing"

	"gotest.tools/assert"
)

func TestAbsPath(t *testing.T) {
	actual, err := absPath(".")
	if err != nil {
		t.Fatal(err)
	}

	expected, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, actual, expected)
}
