.PHONY = clean docker

BINARY  = go-serve
SOURCES = log.go main.go response.go server.go utils.go

go-serve: $(SOURCES)
	go build -o $(BINARY) *.go

docker:
	docker build -t $(BINARY):0.1.0 .

clean:
	rm $(BINARY)

