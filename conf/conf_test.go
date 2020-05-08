package conf

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOpts(t *testing.T) {
	opts := Opts()
	assert.Equal(t, DefaultDir, opts.DirToServe)
	assert.Equal(t, DefaultPort, opts.PortToListen)
}
