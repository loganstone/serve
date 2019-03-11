package server

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"runtime"
	"syscall"
)

func isErrorAddressAlreadyInUse(err error) bool {
	errOpError, ok := err.(*net.OpError)
	if !ok {
		return false
	}
	errSyscallError, ok := errOpError.Err.(*os.SyscallError)
	if !ok {
		return false
	}
	errErrno, ok := errSyscallError.Err.(syscall.Errno)
	if !ok {
		return false
	}
	if errErrno == syscall.EADDRINUSE {
		return true
	}
	const WSAEADDRINUSE = 10048
	if runtime.GOOS == "windows" && errErrno == WSAEADDRINUSE {
		return true
	}
	return false
}

func newServer(dir string, port int) *http.Server {
	return &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: fileServerHandlerWithLogging(dir),
	}
}

// Run ...
func Run(dir string, port int) {
	srv := newServer(dir, port)
	log.Printf("Serving [%s] on HTTP [%s]\n", dir, srv.Addr)

	err := srv.ListenAndServe()
	if err != nil {
		if isErrorAddressAlreadyInUse(err) {
			log.Println(err)
			log.Printf("Change port: [%d]\n", port+1)
			Run(dir, port+1)
		}
		log.Fatal(err)
	}
}
