package main

//go:generate enumer -type=APIResponseStatus -json enums
//go:generate mkdir -p embedded
//go:generate esc -o embedded/assets.go -pkg embedded -prefix "dist/" dist

import (
	"fmt"
	"os"

	"github.com/andrexus/cloud-initer/cmd"
)

func main() {
	if err := cmd.RootCmd().Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to run command: %v\n", err)
		os.Exit(1)
	}
}
