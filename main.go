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

	// Listener port 사용 검사는 생성 시 에러로 확인하는 것이
	// 코드가 명료하고 간단하다.
	// prot 확인과 생성 사이의 시간에 다른 process 에서 port 를 점유한다면,
	// 에러가 발생하기는 마찬가지이고, 코드가 더 복잡해진다.
	ln, err := server.Listener(opts.PortToListen)
	if server.IsErrorAddressAlreadyInUse(err) {
		log.Println(err)
		// port 자동할당
		ln, err = server.Listener(0)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Change port: [%d]\n", ln.Addr().(*net.TCPAddr).Port)
	}

	server.Run(watcher.VerifiedDir, ln)
}
