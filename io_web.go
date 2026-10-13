//go:build wasm

package main

import (
	"syscall/js"
)

// jsWriter captures data written by Go and pushes it to xterm.js
type jsWriter struct {
	terminal js.Value
}

// writes via xterm.js write method
func (w *jsWriter) Write(p []byte) (n int, err error) {
	// call xterm.js terminal.write() method
	w.terminal.Call("write", string(p))
	return len(p), nil
}

// jsReader blocks and waits for input chunks emitted by xterm.js
type jsReader struct {
	ch chan byte
}

// blocks until a byte is available, then reads
func (r *jsReader) Read(p []byte) (n int, err error) {
	if len(p) == 0 {
		return 0, nil
	}
	// block until at least one byte is received
	b := <-r.ch
	p[0] = b
	n = 1

	// get any other available bytes in the channel up to length
	for n < len(p) {
		select {
		case b := <-r.ch:
			p[n] = b
			n++
		default:
			return n, nil
		}
	}
	return n, nil
}

// set up IO for using xterm.js
func InitIO() {
	// fetch the global xterm instance from the browser window context.
	// assumes terminal exposed in JS via window.term = new Terminal();
	term := js.Global().Get("term")

	// setup the custom reader and channel
	inputChan := make(chan byte, 1024)
	InputReader = &jsReader{ch: inputChan}

	// setup the custom writer
	OutputWriter = &jsWriter{terminal: term}

	// bind xterm.js data event to push into Go channel
	onDataCallback := js.FuncOf(func(this js.Value, args []js.Value) any {
		inputStr := args[0].String()
		for i := 0; i < len(inputStr); i++ {
			inputChan <- inputStr[i]
		}
		return nil
	})
	term.Call("onData", onDataCallback)
}
