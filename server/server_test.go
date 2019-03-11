package server

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http/httptest"
	"strings"
	"testing"

	"gotest.tools/assert"

	"github.com/loganstone/serve/conf"
)

var (
	dirToServe   = "."
	portToListen = 9000
)

func newTestServer(dir string, port int) (*httptest.Server, error) {
	handler := fileServerHandlerWithLogging(dir)
	ts := httptest.NewUnstartedServer(handler)
	err := ts.Listener.Close()
	if err != nil {
		return ts, err
	}

	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return ts, err
	}

	ts.Listener = ln
	return ts, nil
}

func TestRunServeAndReqeust(t *testing.T) {
	ts, err := newTestServer(conf.DefaultDir, conf.DefaultPort)
	defer ts.Close()
	if err != nil {
		t.Fatal(err)
	}

	ts.Start()
	client := ts.Client()

	resp, err := client.Get(ts.URL)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("received non-200 response: %d\n", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	actual := string(body)
	files, err := ioutil.ReadDir(dirToServe)
	if err != nil {
		t.Fatal(err)
	}

	for _, f := range files {
		expected := f.Name()
		assert.Assert(t, strings.Contains(actual, expected))
	}
}

func TestIsErrorAddressAlreadyInUse(t *testing.T) {
	firstSrv, err := newTestServer(conf.DefaultDir, conf.DefaultPort)
	defer firstSrv.Close()
	if err != nil {
		t.Fatal(err)
	}
	firstSrv.Start()

	secondSrv, err := newTestServer(conf.DefaultDir, conf.DefaultPort)
	defer secondSrv.Close()
	assert.Assert(t, isErrorAddressAlreadyInUse(err))

	err = errors.New("some error")
	assert.Assert(t, !isErrorAddressAlreadyInUse(err))
}
