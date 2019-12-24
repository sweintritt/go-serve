.PHONY = all clean docker prepare-docker

NAME      = go-serve
SOURCES   = log.go main.go response.go server.go utils.go
BUILD_DIR = ./build
VERSION   = $(shell cat VERSION)

all: docker clean

mux:
	go get -u github.com/gorilla/mux

$(NAME): $(SOURCES) mux
	go build -o $(NAME) *.go

prepare-docker: $(NAME)
	mkdir -p $(BUILD_DIR)/tmp
	cp $(NAME) $(BUILD_DIR)
	cp -r public/ $(BUILD_DIR)/

docker: $(NAME) prepare-docker
	sudo docker build -t $(NAME):$(VERSION) .

clean:
	rm -f $(NAME)
	rm -rf $(BUILD_DIR)
