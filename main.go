package main

import (
	"flag"
	"log"
	"os"

	"github.com/loganstone/serve/conf"
	"github.com/loganstone/serve/dir"
	"github.com/loganstone/serve/server"
)

var (
	dirToServe   = flag.String("d", conf.DefaultDir, "directory to serve")
	portToListen = flag.Int("p", conf.DefaultPort, "port to listen on")
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
