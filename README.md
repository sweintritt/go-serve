go serve
=========

Simple Web UI to upload and excute Go code. Select a source code file on your
machine and get the output in return.

# Build

```sh
$ go get -u github.com/gorilla/mux
$ go build -o go-serve *.go
```

# Docker image

To build a Docker image, the built binary including all required libraries
must be copied into `bin/`.

```sh
$ sudo docker build -t go-serve:0.1.0 .
```

> NOTE: Currently the web ui is available on port 8081 and cannot be changed.
