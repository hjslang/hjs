package hjs

import (
	"regexp"
	"strings"

	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/printer"
)

var onAttrRe = regexp.MustCompile(`^on([A-Z]\w*)`)

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

const RxRuntime = `function rx(initVal) {
  let _val = initVal;
  let _listeners = [];
  return {
    get: function () {
      return _val;
    },
    set: function (val) {
      _val = val;
      for (let listener of _listeners) {
        listener();
      }
    },
  };
}`

func RxCompiler(pr *printer.Printer, node ast.Node, next func(node ast.Node) error) error {
	switch v := node.(type) {
	case *RxVariable:
		pr.Print(v.Variable.Name, ".get()")
	case *RxLetStmt:
		pr.Print("let ", v.LetStmt.Name, " = rx(", v.LetStmt.Value, ")")
	case *RxAssignStmt:
		pr.Print(v.AssignStmt.Name, ".set(", v.AssignStmt.Value, ")")
	default:
		return next(node)
	}
	return nil
}
