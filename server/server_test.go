package server

import (
	"errors"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/loganstone/serve/conf"
)

var testLn net.Listener
var testSrv *http.Server

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func setup() {
	ln, err := Listener(conf.DefaultPort)
	if err != nil {
		log.Fatal(err)
	}
	testLn = ln
	testSrv = newServer(conf.DefaultDir, testLn.Addr().(*net.TCPAddr).Port, false)
	go testSrv.Serve(ln)
}

func teardown() {
	testLn.Close()
	testSrv.Close()
}

func TestRequest(t *testing.T) {
	retry := 0
	for {
		resp, err := http.Get("http://" + testLn.Addr().String())
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
			assert.Contains(t, actual, expected)
		}
		defer resp.Body.Close()
		break
	}
}

func TestIsErrorAddressAlreadyInUse(t *testing.T) {
	ln, err := Listener(conf.DefaultPort)
	defer ln.Close()
	assert.True(t, IsErrorAddressAlreadyInUse(err))

	err = errors.New("some error")
	assert.False(t, IsErrorAddressAlreadyInUse(err))
}
