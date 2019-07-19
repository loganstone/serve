package main

import (
	"log"
	"net"
	"os"

	"github.com/loganstone/serve/conf"
	"github.com/loganstone/serve/dir"
	"github.com/loganstone/serve/server"
)

func main() {
	opts := conf.Opts()

	absPath, err := dir.Abs(opts.DirToServe)
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

	watcher, err := dir.NewWatcher(absPath)
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()
	go dir.NowMyWatchBegins(absPath, watcher)

	ln, err := server.Listener(opts.PortToListen)
	if server.IsErrorAddressAlreadyInUse(err) {
		log.Println(err)
		ln, err = server.Listener(0)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Change port: [%d]\n", ln.Addr().(*net.TCPAddr).Port)
	}

	server.Run(absPath, ln)
}
