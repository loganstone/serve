# Serve

[![Go Report Card](https://goreportcard.com/badge/github.com/loganstone/serve)](https://goreportcard.com/report/github.com/loganstone/serve)

Alternative to "python -m http.server --bind 0 9000"

## Getting started

### Prerequisites

* [Install golang](https://golang.org/doc/install) ;)

### Install and run

* Install `serve`

```shell
$ go install github.com/loganstone/serve@latest
```

* Run

```shell
$ serve
```

### First contact

Click! http://localhost:9000

### Usage of Serve

```shell
$ serve -h
usage of serve:
  -d string
    	directory to serve (default ".")
  -p int
    	port to listen on (default 9000)
```

## Running Tests

```shell
$ go test -v -count=1 ./...  # no cached
```

## Key Features

- Serving directory to http.
- Logging.

## To-do Features

- Add something if need.
