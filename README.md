go serve
=========

Simple Web UI to upload and excute Go code. Select a source code file on your
machine and get the output in return.

# Build

## Preparation

*go-serve* requires `gorilla/mux`. Install it with

```sh
$ go get -u github.com/gorilla/mux
```

## Build

Build the project with

```sh
$ make go-serve
```

## Docker image

Build the docker image with

```sh
$ make docker
```

and start a new container with

```sh
```

> NOTE: Currently the web ui is available on port 8081 and cannot be changed.
