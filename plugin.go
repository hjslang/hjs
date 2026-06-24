package hjs

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/builder"
	"github.com/xjslang/xjs/js"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/scanner"
	"github.com/xjslang/xjs/token"
)

var (
	startTag = token.RegisterType("start-tag")
	endTag   = token.RegisterType("end-tag")
)

type Attr struct {
	Name  *js.Ident
	Value ast.Expr
}

type Tag struct {
	ast.BaseExpr
	Layout struct {
		StartTag token.Token
		EndTag   token.Token
	}
	Name     *js.Ident
	Attrs    []Attr
	Children []ast.Expr
}

type RxVariable struct {
	ast.BaseExpr
	Variable *js.Variable
}

type RxLetStmt struct {
	ast.BaseStmt
	LetStmt *js.LetStmt
}

type RxAssignStmt struct {
	ast.BaseStmt
	AssignStmt *js.AssignStmt
}

func ParseTag(p *parser.Parser) (_ *Tag, err error) {
	node := &Tag{}
	if node.Layout.StartTag, err = p.Expect(startTag); err != nil {
		return
	}
	if node.Name, err = js.ParseIdent(p); err != nil {
		return
	}
	for p.CurrentToken.Type != token.GT {
		var attr Attr
		if attr.Name, err = js.ParseIdent(p); err != nil {
			return
		}
		if _, err = p.Expect(token.ASSIGN); err != nil {
			return
		}
		if attr.Value, err = js.ParseValue(p); err != nil {
			return
		}
		node.Attrs = append(node.Attrs, attr)
	}
	if _, err = p.Expect(token.GT); err != nil {
		return
	}
	for p.CurrentToken.Type != endTag {
		var child ast.Expr
		if child, err = p.ParseExpr(); err != nil {
			return
		}
		node.Children = append(node.Children, child)
	}
	if node.Layout.EndTag, err = p.Expect(endTag); err != nil {
		return
	}
	var ident *js.Ident
	if ident, err = js.ParseIdent(p); err != nil {
		return
	}
	if ident.Literal != node.Name.Literal {
		return nil, p.ErrorAt(
			ident.Token,
			"expected closing tag </"+node.Name.Literal+">",
		)
	}
	if _, err = p.Expect(token.GT); err != nil {
		return
	}
	return node, nil
}

// Plugin enriches the JavaScript parser, so that we can parse expressions that are not part of the JS standard.
func Plugin(b *builder.Builder) {
	token.RegisterUnaryType(startTag)

	// now the parser can "scan" '<' and '</'
	b.UseScanner(func(sc *scanner.Scanner, next func() (token.Token, error)) (tok token.Token, err error) {
		if tok, err = next(); err != nil {
			return
		}
		if tok.Type == token.LT {
			c := sc.CurrentChar()
			switch {
			case scanner.IsLetter(c):
				tok.Type = startTag
			case c == '/':
				sc.AdvanceChar()
				tok.Type = endTag
				tok.Literal = "</"
			}
		}
		return
	})

	// now the parser can "parse" HTML tags
	b.UseUnaryParser(func(p *parser.Parser, next func() (ast.Expr, error)) (_ ast.Expr, err error) {
		if p.CurrentToken.Type == startTag {
			return ParseTag(p)
		}
		return next()
	})
}

func RxPlugin(b *builder.Builder) {
	b.UseStmtParser(func(p *parser.Parser, next func() (ast.Stmt, error)) (node ast.Stmt, err error) {
		if node, err = next(); err != nil {
			return
		}
		switch v := node.(type) {
		case *js.LetStmt:
			if v.Name.Literal[0] == '$' {
				w := &RxLetStmt{}
				w.LetStmt = v
				return w, nil
			}
		case *js.AssignStmt:
			if v.Name.Literal[0] == '$' {
				w := &RxAssignStmt{}
				w.AssignStmt = v
				return w, nil
			}
		}
		return
	})

	b.UseExprParser(func(p *parser.Parser, next func() (ast.Expr, error)) (node ast.Expr, err error) {
		if node, err = next(); err != nil {
			return
		}
		if node, ok := node.(*js.Variable); ok {
			if node.Name.Literal[0] == '$' {
				v := &RxVariable{}
				v.Variable = node
				return v, nil
			}
		}
		return
	})
}
