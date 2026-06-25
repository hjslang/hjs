package hjs

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/printer"
)

func Formatter(pr *printer.Printer, node ast.Node, next func(node ast.Node) error) error {
	switch v := node.(type) {
	case *tag:
		pr.SpPrint(v.Layout.StartTag).Print(v.Name)
		for _, a := range v.Attrs {
			pr.SpPrint(a.Name).Print("=", a.Value)
		}
		pr.Print(">")
		pr.IncreaseIndent()
		for _, child := range v.Children {
			pr.Print(child)
		}
		pr.DecreaseIndent()
		pr.Print(v.Layout.EndTag, v.Name, ">")
	default:
		return next(node)
	}
	return nil
}
