package main

import (
	"flag"
	"log"
	"os"

	"github.com/loganstone/serve/dir"
	"github.com/loganstone/serve/server"
)

const (
	defaultDir  = "."
	defaultPort = 9000
)

var (
	dirToServe   = flag.String("d", defaultDir, "directory to serve")
	portToListen = flag.Int("p", defaultPort, "port to listen on")
)

func main() {
	flag.Parse()

	absPath, err := dir.Abs(*dirToServe)
	if err != nil {
		log.Fatal(err)
	}

	dirInfo, err := os.Stat(absPath)
	if err != nil {
		log.Fatal(err)
	}

	if !dirInfo.IsDir() {
		log.Fatal("-d option value must be directory")
	}

	server.Run(absPath, *portToListen)
}
