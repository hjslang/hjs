package hjs

import (
	"github.com/xjslang/xjs"
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/js"
	"github.com/xjslang/xjs/printer"
)

type compileConfig struct {
	withRuntime bool
}

type CompileOption func(*compileConfig)

func WithRuntime() CompileOption {
	return func(cfg *compileConfig) {
		cfg.withRuntime = true
	}
}

func Parse(input []byte) (*js.Program, error) {
	b := xjs.NewBuilder().
		Install(Plugin).
		Install(RxPlugin)
	p := b.Build(input)
	return js.ParseProgram(p)
}

// hjs.Compile(result, hjs.WithRuntime())
func Compile(result ast.Node, opts ...CompileOption) (string, error) {
	cfg := &compileConfig{
		withRuntime: false,
	}
	for _, opt := range opts {
		opt(cfg)
	}
	pr := xjs.NewPrinter(printer.Compact())
	pr.UsePrinter(Compiler)
	pr.UsePrinter(RxCompiler)
	if cfg.withRuntime {
		pr.Print(RxRuntime, "\n\n")
	}
	pr.Print(result)
	return pr.Output()
}

func Format(result ast.Node, opts ...printer.Option) (string, error) {
	pr := xjs.NewPrinter(opts...)
	pr.UsePrinter(Formatter)
	pr.UsePrinter(RxFormatter)
	pr.Print(result)
	return pr.Output()
}
