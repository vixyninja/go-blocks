package main

import (
	"fmt"
	"os"
)

var (
	version = "dev"
	commit  = "unknown"
	date    = "unknown"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "version" {
		fmt.Printf("go-blocks version %s (commit: %s, built: %s)\n", version, commit, date)
		return
	}

	fmt.Println("go-blocks - My custom blocks for golang projects")
	fmt.Printf("Version: %s\n", version)
}
