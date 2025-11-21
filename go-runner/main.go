// main.go
// Automatically handled with Promise rejects when returning an error!
package main

import (
	"bytes"
	"fmt"
	"syscall/js"

	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
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
}

func newStreamingWriter() *streamingWriter {
	sw := &streamingWriter{
		buf: &bytes.Buffer{},
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
				result := errorResult{
					message: fmt.Sprintf("Panic: %v", r),
					output:  "",
				}
				if !onComplete.IsUndefined() {
					onComplete.Invoke(result.ToValue())
				}
			}
		}()

		streamWriter := newStreamingWriter()

		i := interp.New(interp.Options{
			Stdout:       streamWriter,
			Stderr:       streamWriter,
			Unrestricted: true,
		})

		i.Use(stdlib.Symbols)

		prog, err := i.Compile(code)
		if err != nil {
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
	println("WASM Go Initialized")

	<-c
}
