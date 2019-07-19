package conf

import "flag"

const (
	// DefaultDir ...
	DefaultDir = "."
	// DefaultPort ...
	DefaultPort = 9000
)

// Options .
type Options struct {
	DirToServe   string
	PortToListen int
}

var opts Options

func init() {
	flag.StringVar(&opts.DirToServe, "d", DefaultDir, "Directory to serve")
	flag.IntVar(&opts.PortToListen, "p", DefaultPort, "Port to listen on")
	flag.Parse()
}

// Opts .
func Opts() Options {
	return opts
}
