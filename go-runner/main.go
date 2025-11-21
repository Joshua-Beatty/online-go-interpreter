// main.go
// Automatically handled with Promise rejects when returning an error!
package main

import (
	"bytes"
	"fmt"
	"sync/atomic"
	"syscall/js"

	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
)

var (
	// currentRunID tracks the latest execution generation.
	// Old runs will check this and exit/stop streaming if they don't match.
	currentRunID int64
)

type errorResult struct {
	message string
	output  string
}
type successResult struct {
	output string
}

func (e errorResult) ToValue() js.Value {
	return js.ValueOf(map[string]any{
		"success": false,
		"message": e.message,
		"output":  e.output,
	})
}

func (e successResult) ToValue() js.Value {
	return js.ValueOf(map[string]any{
		"success": true,
		"output":  e.output,
	})
}

// streamingWriter writes output to both a buffer and streams it to JavaScript
type streamingWriter struct {
	buf     *bytes.Buffer
	jsFunc  js.Value
	hasFunc bool
	runID   int64
}

func newStreamingWriter(runID int64) *streamingWriter {
	sw := &streamingWriter{
		buf:   &bytes.Buffer{},
		runID: runID,
	}

	// Try to get the JavaScript streaming function
	jsGlobal := js.Global()
	streamFunc := jsGlobal.Get("__goStreamOutput")
	if !streamFunc.IsUndefined() && streamFunc.Type() == js.TypeFunction {
		sw.jsFunc = streamFunc
		sw.hasFunc = true
	}

	return sw
}

func (sw *streamingWriter) Write(p []byte) (n int, err error) {
	// Check if this writer belongs to the current run
	if atomic.LoadInt64(&currentRunID) != sw.runID {
		// This is an old run, discard output silently
		return len(p), nil
	}

	// Write to buffer for final output
	n, err = sw.buf.Write(p)

	// Also stream to JavaScript if function is available
	if sw.hasFunc && n > 0 {
		// Convert bytes to string and call JavaScript function
		data := string(p[:n])
		sw.jsFunc.Invoke(data)
	}

	return n, err
}

func (sw *streamingWriter) String() string {
	return sw.buf.String()
}

func run(this js.Value, args []js.Value) (ret any) {
	if len(args) < 1 || len(args) > 2 {
		return errorResult{
			message: fmt.Sprintf("Wrong number of arguments, expected: 1-2 got: %v", len(args)),
			output:  "",
		}.ToValue()
	}

	// Increment run ID for this new execution
	runID := atomic.AddInt64(&currentRunID, 1)

	code := args[0].String()

	// Get the onComplete callback if provided
	var onComplete js.Value
	if len(args) >= 2 && args[1].Type() == js.TypeFunction {
		onComplete = args[1]
	}

	// Run asynchronously in a goroutine
	go func() {
		defer func() {
			if r := recover(); r != nil {
				// Check if we are still the current run
				if atomic.LoadInt64(&currentRunID) != runID {
					return
				}

				result := errorResult{
					message: fmt.Sprintf("Panic: %v", r),
					output:  "",
				}
				if !onComplete.IsUndefined() {
					onComplete.Invoke(result.ToValue())
				}
			}
		}()

		streamWriter := newStreamingWriter(runID)

		i := interp.New(interp.Options{
			Stdout:       streamWriter,
			Stderr:       streamWriter,
			Unrestricted: true,
		})

		i.Use(stdlib.Symbols)

		prog, err := i.Compile(code)
		if err != nil {
			// Check if we are still the current run
			if atomic.LoadInt64(&currentRunID) != runID {
				return
			}

			result := errorResult{
				message: fmt.Sprintf("failed to compile code: %v", err),
				output:  streamWriter.String(),
			}
			if !onComplete.IsUndefined() {
				onComplete.Invoke(result.ToValue())
			}
			return
		}

		_, err = i.Execute(prog)

		// Check if we are still the current run
		if atomic.LoadInt64(&currentRunID) != runID {
			return
		}

		if err != nil {
			result := errorResult{
				message: "code exited with error",
				output:  streamWriter.String(),
			}
			if !onComplete.IsUndefined() {
				onComplete.Invoke(result.ToValue())
			}
			return
		}

		result := successResult{
			output: streamWriter.String(),
		}
		if !onComplete.IsUndefined() {
			onComplete.Invoke(result.ToValue())
		}
	}()

	// Return immediately to avoid blocking
	return js.ValueOf(map[string]any{
		"async": true,
	})
}

func registerCallbacks() {
	js.Global().Set("run", js.FuncOf(run))
}

func main() {
	c := make(chan struct{}, 0)

	// register functions
	registerCallbacks()
	println("WASM Go Initialized!")

	<-c
}
