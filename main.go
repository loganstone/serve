package main

import (
	"flag"
	"log"
	"net"
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

	watcher, err := dir.NewWatcher(absPath)
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	go func() {
		log.Printf("verified directory [%s], and now my watch begins", absPath)
		for {
			select {
			case event := <-watcher.Events:
				log.Printf("Watcher - %s\n", event)

			case err := <-watcher.Errors:
				log.Println("Watcher ERROR!", err)
			}
		}
	}()

	ln, err := server.Listener(*portToListen)
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
