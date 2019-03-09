package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"gotest.tools/assert"
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

func TestAbs(t *testing.T) {
	absPath, err := abs(dirToServe)
	if err != nil {
		t.Fatal(err)
	}

	expected, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, absPath, expected)

	testPath := "/abs/path"
	_, err = abs(testPath)
	assert.Assert(t, err == nil)
}

func TestRunServeAndReqeust(t *testing.T) {
	ts, err := newTestServer(dirToServe, portToListen)
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
	firstSrv, err := newTestServer(dirToServe, portToListen)
	defer firstSrv.Close()
	if err != nil {
		t.Fatal(err)
	}
	firstSrv.Start()

	secondSrv, err := newTestServer(dirToServe, portToListen)
	defer secondSrv.Close()
	assert.Assert(t, isErrorAddressAlreadyInUse(err))
}
