.PHONY: all build build-linux build-darwin build-windows build-all clean docker-build

BINARY := sdk
VERSION ?= dev

all: build

build:
	go build -o $(BINARY) ./

build-linux:
	GOOS=linux GOARCH=amd64 go build -o $(BINARY)-linux-amd64 ./

build-darwin:
	GOOS=darwin GOARCH=amd64 go build -o $(BINARY)-darwin-amd64 ./

build-windows:
	GOOS=windows GOARCH=amd64 go build -o $(BINARY)-windows-amd64.exe ./

build-all: build-linux build-darwin build-windows

clean:
	rm -f $(BINARY) $(BINARY)-* *.exe

# Build a Docker image (optional)
docker-build:
	docker build -t sudoku-solver:latest .
