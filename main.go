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

	// 지정된 port 가 이미 점유되었는지 확인은
	// Listener 생성 시 error 로 확인하도록 한다.
	// Listener 생성 전 port 점유 확인 코드가 있다면
	// 코드만 복잡하고, 아래와 같은 상황에 같은 error 를 발생시킨다.
	// 1. Port 점유 확인 코드가 실행되어 port 미 점유 확인
	// 2. 다른 process 에서 port 점유
	// 3. Listener 생성 코드 실행 시 에러
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
