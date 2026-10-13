//go:build !wasm

package main

import (
	"os"
)

// standard io for a regular cli application
func InitIO() {
	InputReader = os.Stdin
	OutputWriter = os.Stdout
}
