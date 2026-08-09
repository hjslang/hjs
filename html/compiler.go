package html

import (
	"regexp"
	"strings"

	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/jsextended"
	"github.com/xjslang/xjs/printer"
)

var (
	onAttrRe   = regexp.MustCompile(`^on([A-Z]\w*)`)
	dataAttrRe = regexp.MustCompile(`^data([A-Z]\w*)`)
)

// Compiler transforms the code to valid JS code.
func Compiler(pr *printer.Printer, node ast.Node, next func(node ast.Node) error) error {
	switch v := node.(type) {
	case *Tag:
		pr.Print("(function(){")
		pr.Print("const elem = document.createElement('", v.Name, "');")
		for _, attr := range v.Attrs {
			if matches := onAttrRe.FindStringSubmatch(attr.Name.Literal); matches != nil {
				name := strings.ToLower(matches[1])
				pr.Print("elem.addEventListener('", name, "', ", attr.Value, ");")
			} else if matches := dataAttrRe.FindStringSubmatch(attr.Name.Literal); matches != nil {
				name := matches[1]
				name = strings.ToLower(name[:1]) + name[1:]
				pr.Print("elem.dataset.", name, " = ", attr.Value, ";")
			} else {
				pr.Print("elem.setAttribute('", attr.Name, "', ", attr.Value, ");")
			}
		}
		for _, child := range v.Children {
			pr.Print("elem.append(", child, ");")
		}
		pr.Print("return elem})()")
	default:
		return next(node)
	}
	return nil
}

func Compile(result ast.Node) (string, error) {
	pr := jsextended.PrinterBuilder().
		UsePrinter(Compiler).
		Build(printer.Compact())
	pr.Print(result)
	return pr.Output()
}
