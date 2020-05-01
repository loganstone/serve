package conf

import "flag"

const (
	// DefaultDir -d 옵션 값이 입력되지 않았을때 사용된다.
	DefaultDir = "."
	// DefaultPort -p 옵션 값이 입력되지 않았을대 사용된다.
	DefaultPort = 9000
)

// Options CLI 입력 옵션을 타입화.
type Options struct {
	DirToServe   string
	PortToListen int
}

var opts Options

func init() {
	flag.StringVar(&opts.DirToServe, "d", DefaultDir, "Directory to serve")
	flag.IntVar(&opts.PortToListen, "p", DefaultPort, "Port to listen on")
}

// Opts 'init' 에서 처리 된 -d, -p 값을 담은 type 을 반환하다.
func Opts() Options {
	flag.Parse()
	return opts
}
