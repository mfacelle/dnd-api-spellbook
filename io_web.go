//go:build wasm

package main

import (
	"log"
	"strings"
	"syscall/js"
)

// jsWriter captures data written by Go and pushes it to xterm.js
type jsWriter struct {
	terminal js.Value
}

// writes via xterm.js write method
func (w *jsWriter) Write(p []byte) (n int, err error) {
	// replace newlines with windows-style newlines
	outputStr := strings.ReplaceAll(string(p), "\n", "\r\n")
	// call xterm.js terminal.write() method
	w.terminal.Call("write", outputStr)
	return len(p), nil
}

// jsReader blocks and waits for input chunks emitted by xterm.js
type jsReader struct {
	ch      chan []byte
	pending []byte
}

// blocks until a byte is available, then reads
func (r *jsReader) Read(p []byte) (n int, err error) {
	if len(p) == 0 {
		return 0, nil
	}
	if len(r.pending) == 0 {
		r.pending = <-r.ch
	}
	n = copy(p, r.pending)
	r.pending = r.pending[n:]
	log.Println("received " + string(p[:n]))
	return n, nil
}

// set up IO for using xterm.js
func InitIO() {
	// fetch the global xterm instance from the browser window context.
	// assumes terminal exposed in JS via window.term = new Terminal();
	term := js.Global().Get("term")

	// setup the custom reader and channel
	inputChan := make(chan []byte, 16)
	reader := &jsReader{ch: inputChan}
	InputReader = reader
	js.Global().Set("handleTerminalInput", js.FuncOf(func(this js.Value, args []js.Value) any {
		log.Println("handleTerminalInput")
		if len(args) == 0 {
			return false
		}

		inputChan <- []byte(args[0].String() + "\n")
		return true
	}))

	// setup the custom writer
	OutputWriter = &jsWriter{terminal: term}
}
