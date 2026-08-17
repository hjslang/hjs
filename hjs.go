package hjs

import (
	"github.com/xjslang/hjs/html"
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/js"
	"github.com/xjslang/xjs/printer"
)

func Parse(input string) (*js.Program, error) {
	return html.Parse(input)
}

func Compile(result ast.Node) (string, error) {
	return html.Compile(result)
}

func Format(result ast.Node, opts ...printer.Option) (string, error) {
	return html.Format(result, opts...)
}
