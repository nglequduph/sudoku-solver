# Minimal Go test project

This workspace contains a single-file Go program to experiment quickly.

Files:

- `main.go`: sample program with an `Add` function and a runnable `main`.

Run:

```sh
go run main.go
go run main.go -a 5 -b 7
```

Build (multi-platform)

You can build locally or produce cross-platform binaries using the provided `Makefile`.

Quick local build:

```sh
go build -o sdk ./
./sdk
```

Use `make` targets to produce platform-specific binaries:

```sh
# build for current platform
make build

# build specific platforms
make build-linux    # linux/amd64 -> sdk-linux-amd64
make build-darwin   # macOS darwin/amd64 -> sdk-darwin-amd64
make build-windows  # windows/amd64 -> sdk-windows-amd64.exe

# build all three
make build-all

# remove binaries
make clean

# optional: build a Docker image
make docker-build
```
