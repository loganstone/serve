package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"gotest.tools/assert"
)

func newTestServer() *httptest.Server {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", portToListen))
	if err != nil {
		log.Fatal(err)
	}

	handler := wrapHandlerWithLogging(http.FileServer(http.Dir(dirToServe)))
	ts := httptest.NewUnstartedServer(handler)
	ts.Listener.Close()
	ts.Listener = ln
	return ts
}

func TestAbsPath(t *testing.T) {
	actual, err := absPath(".")
	if err != nil {
		t.Fatal(err)
	}

	expected, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, actual, expected)
}

func TestRunServeAndReqeust(t *testing.T) {
	dirToServe = defaultDir
	portToListen = defaultPort

	ts := newTestServer()

	ts.Start()
	client := ts.Client()
	defer ts.Close()

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
