package server

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"runtime"
	"syscall"
	"time"
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

// IsPortInUse ...
func IsPortInUse(port int) bool {
	host := fmt.Sprintf(":%d", port)
	timeoutSecs := 5
	conn, err := net.DialTimeout(
		"tcp", host, time.Duration(timeoutSecs)*time.Second)

	if conn != nil {
		defer conn.Close()
	}

	if err != nil {
		return false
	}

	return true
}

// FreePort ...
func FreePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}

// Run ...
func Run(dir string, port int) {
	srv := newServer(dir, port)
	log.Printf("Serving [%s] on HTTP [%s]\n", dir, srv.Addr)

	err := srv.ListenAndServe()
	if err != nil {
		if isErrorAddressAlreadyInUse(err) {
			log.Println(err)
			freePort, err := FreePort()
			if err != nil {
				log.Fatal(err)
			}
			log.Printf("Change port: [%d]\n", freePort)
			Run(dir, freePort)
		}
		log.Fatal(err)
	}
}
