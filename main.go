package main

import (
	"fmt"
	"os"

	"example.com/sdk/server"
)

func main() {
	port := "8081"
	if len(os.Args) > 1 {
		port = os.Args[1]
	}

	fmt.Printf("Starting Sudoku Solver Service on port %s...\n", port)
	server.Start(port)
}
