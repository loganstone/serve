package main

import (
	"log"
	"net"

	"github.com/loganstone/serve/conf"
	"github.com/loganstone/serve/dir"
	"github.com/loganstone/serve/server"
)

func main() {
	opts := conf.Opts()

	watcher, err := dir.NewWatcher(opts.DirToServe)
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()
	go watcher.NowMyWatchBegins()

	ln, err := server.Listener(opts.PortToListen)
	if server.IsErrorAddressAlreadyInUse(err) {
		log.Println(err)
		ln, err = server.Listener(0)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Change port: [%d]\n", ln.Addr().(*net.TCPAddr).Port)
	}

	server.Run(watcher.VerifiedDir, ln)
}
