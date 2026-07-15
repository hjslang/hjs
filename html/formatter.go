package html

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/printer"
)

func Formatter(pr *printer.Printer, node ast.Node, next func(node ast.Node) error) error {
	switch v := node.(type) {
	case *Tag:
		pr.Space().Print(v.Layout.StartTag, v.Name)
		for _, a := range v.Attrs {
			pr.Space().Print(a.Name, "=", a.Value)
		}
		if v.SelfClosing {
			pr.Space().Print(v.Layout.SelfClosingTag)
		} else {
			pr.Print(">")
			pr.IncreaseIndent()
			for _, child := range v.Children {
				pr.Print(child)
			}
			pr.DecreaseIndent()
			pr.Print(v.Layout.EndTag, v.Name, ">")
		}
	default:
		return next(node)
	}
	return nil
}
