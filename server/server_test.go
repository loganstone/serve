package server

import (
	"errors"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"testing"

	"gotest.tools/assert"

	"github.com/loganstone/serve/conf"
)

func TestRunServeAndReqeust(t *testing.T) {
	ln, err := Listener(conf.DefaultPort)
	if err != nil {
		t.Fatal(err)
	}
	srv := newServer(conf.DefaultDir, ln.Addr().(*net.TCPAddr).Port, false)
	errc := make(chan error, 1)

	go func() {
		retry := 0
		for {
			resp, err := http.Get("http://" + ln.Addr().String())
			if err != nil {
				if retry > 4 {
					t.Fatal(err)
				}
				retry++
				continue
			}

			if resp.StatusCode != 200 {
				t.Fatalf("received non-200 response: %d\n", resp.StatusCode)
			}

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Fatal(err)
			}

			actual := string(body)
			files, err := ioutil.ReadDir(conf.DefaultDir)
			if err != nil {
				t.Fatal(err)
			}

			for _, f := range files {
				expected := f.Name()
				assert.Assert(t, strings.Contains(actual, expected))
			}
			defer resp.Body.Close()
			defer ln.Close()
			defer srv.Close()
			break
		}
	}()

	go func() {
		errc <- srv.Serve(ln)
	}()
	<-errc
}

func TestIsErrorAddressAlreadyInUse(t *testing.T) {
	firstLn, err := Listener(conf.DefaultPort)
	defer firstLn.Close()
	if err != nil {
		t.Fatal(err)
	}
	secondLn, err := Listener(conf.DefaultPort)
	defer secondLn.Close()
	assert.Assert(t, IsErrorAddressAlreadyInUse(err))

	err = errors.New("some error")
	assert.Assert(t, !IsErrorAddressAlreadyInUse(err))
}
