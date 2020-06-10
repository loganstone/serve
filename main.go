package main

import (
	"flag"
	"log"
	"net"

	"github.com/loganstone/serve/conf"
	"github.com/loganstone/serve/dir"
	"github.com/loganstone/serve/server"
)

func main() {
	opts := conf.Opts()

	if opts.DirToServe == "" {
		flag.PrintDefaults()
		return
	}

	watcher, err := dir.NewWatcher(opts.DirToServe)
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()
	go watcher.NowMyWatchBegins()

	// Listener port 가 이미 사용되는지 확인은
	// 생성 시 에러로 확인하는 것이 코드가 명료하고 간단하다.
	// prot 확인과 생성 코드가 분리되어있다면,
	// 실행 사이의 시간에 다른 process 에서 port 를 점유할 수 있어,
	// 에러가 발생하기는 마찬가지이고, 코드가 더 복잡해진다.
	ln, err := server.Listener(opts.PortToListen)
	if server.IsErrorAlreadyInUse(err) {
		log.Println(err)
		// port 자동할당
		ln, err = server.Listener(0)
		if err != nil {
			log.Println(err)
			log.Fatal("failed port auto-assignment")
		}
		message := "port [%d] is already in use. change port to : [%d]\n"
		log.Printf(message, opts.PortToListen, ln.Addr().(*net.TCPAddr).Port)
	}

	server.Run(watcher.VerifiedDir, ln)
}
