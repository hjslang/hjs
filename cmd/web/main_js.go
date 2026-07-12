package main

import (
	"syscall/js"

	"github.com/xjslang/hjs"
	"github.com/xjslang/xjs/printer"
)

func main() {
	js.Global().Set("hjs", js.ValueOf(map[string]any{
		"compile": js.FuncOf(func(this js.Value, args []js.Value) any {
			input := args[0].String()
			result, err := hjs.Parse([]byte(input))
			if err != nil {
				return map[string]any{"error": err.Error()}
			}
			code, err := hjs.Compile(result)
			if err != nil {
				return map[string]any{"error": err.Error()}
			}
			return map[string]any{"code": code}
		}),
		"format": js.FuncOf(func(this js.Value, args []js.Value) any {
			input := args[0].String()
			result, err := hjs.Parse([]byte(input))
			if err != nil {
				return map[string]any{"error": err.Error()}
			}
			code, err := hjs.Format(result, printer.WithIndent("\t"))
			if err != nil {
				return map[string]any{"error": err.Error()}
			}
			return map[string]any{"code": code}
		}),
	}))

	// keep the program running
	select {}
}
