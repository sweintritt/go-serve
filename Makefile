.PHONY = all clean docker prepare-docker

NAME      = go-serve
SOURCES   = src/*.go
BUILD_DIR = ./build
VERSION   = $(shell cat VERSION)

all: docker clean

mux:
	go get -u github.com/gorilla/mux

prepare:
	mkdir -f $(BUILD_DIR)

$(NAME): $(SOURCES) mux
	go build -o $(BUILD_DIR)/$(NAME) $(SOURCES)

prepare-docker: $(NAME)
	cp -r public/ $(BUILD_DIR)/

docker: $(NAME) prepare-docker
	sudo docker build -t $(NAME):$(VERSION) .

clean:
	rm -f $(NAME)
	rm -rf $(BUILD_DIR)
