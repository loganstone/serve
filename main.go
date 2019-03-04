package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"syscall"
)

const defaultDir = "."
const defaultPort = 9000

var dirToServe string
var portToListen int

type hasStatusCodeResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *hasStatusCodeResponseWriter) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

func absPath(path string) (string, error) {
	if filepath.IsAbs(path) {
		return path, nil
	}

	path, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	return path, nil
}

func wrapHandlerWithLogging(wrappedHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("-> %s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		hrw := &hasStatusCodeResponseWriter{w, http.StatusOK}
		wrappedHandler.ServeHTTP(hrw, r)

		statusCode := hrw.statusCode
		log.Printf("<- %d %s\n", statusCode, http.StatusText(statusCode))
	})
}

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

func run() {
	addr := fmt.Sprintf(":%d", portToListen)

	fileServerHandler := http.FileServer(http.Dir(dirToServe))

	log.Printf("Serving [%s] on HTTP [%s]\n", dirToServe, addr)
	err := http.ListenAndServe(addr, wrapHandlerWithLogging(fileServerHandler))

	if err != nil {
		if isErrorAddressAlreadyInUse(err) {
			log.Println(err)
			portToListen++
			log.Printf("Change port: [%d]\n", portToListen)
			run()
		}
		log.Fatal(err)
	}
}

func init() {
	flag.StringVar(&dirToServe, "d", defaultDir, "directory to serve")
	flag.IntVar(&portToListen, "p", defaultPort, "port to listen on")
}

func main() {
	flag.Parse()

	dirToServe, err := absPath(dirToServe)
	if err != nil {
		log.Fatal(err)
	}

	if _, err := os.Stat(dirToServe); err != nil {
		log.Fatal(err)
	}

	run()
}
