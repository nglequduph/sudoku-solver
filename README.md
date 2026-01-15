# Minimal Go test project

This workspace contains a single-file Go program to experiment quickly.

Files:

- `main.go`: sample program with an `Add` function and a runnable `main`.

Run:

```sh
go run main.go
go run main.go -a 5 -b 7
```

Build:

```sh
go build -o sdk
./sdk -a 4 -b 6
```
