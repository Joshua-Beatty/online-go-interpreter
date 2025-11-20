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
}
type successResult struct {
	output any
}

func (e errorResult) ToValue() js.Value {
	return js.ValueOf(map[string]any{
		"success": false,
		"message": e.message,
	})
}

func (e successResult) ToValue() js.Value {
	return js.ValueOf(map[string]any{
		"success": true,
		"output":  e.output,
	})
}

func run(this js.Value, args []js.Value) (ret any) {
	if len(args) != 1 {
		return errorResult{
			message: fmt.Sprintf("Wrong number of arguments, expected: 1 got: %v", len(args)),
		}.ToValue()
	}

	defer func() {
		if r := recover(); r != nil {
			ret = errorResult{
				message: fmt.Sprintf("Failed to convert argument: %v", r),
			}.ToValue()
		}
	}()

	code := args[0].String()
	var buf bytes.Buffer

	i := interp.New(interp.Options{
		Stdout:       &buf,
		Stderr:       &buf,
		Unrestricted: true,
	})

	i.Use(stdlib.Symbols)

	_, err := i.Eval(code)
	if err != nil {
		panic(err)
	}

	return successResult{
		output: buf.String(),
	}.ToValue()
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
