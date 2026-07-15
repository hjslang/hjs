package html

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/js"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/plugin"
	"github.com/xjslang/xjs/scanner"
	"github.com/xjslang/xjs/token"
)

var (
	startTag       = token.RegisterType("<")
	endTag         = token.RegisterType("</")
	selfClosingTag = token.RegisterType("/>")
)

type Attr struct {
	Name  *js.Ident
	Value ast.Expr
}

type Tag struct {
	ast.BaseExpr
	Layout struct {
		StartTag       token.Token
		EndTag         token.Token
		SelfClosingTag token.Token
	}
	SelfClosing bool
	Name        *js.Ident
	Attrs       []Attr
	Children    []ast.Expr
}

func parseTag(p *parser.Parser) (node *Tag, err error) {
	node = &Tag{}
	if node.Layout.StartTag, err = p.Expect(startTag); err != nil {
		return
	}
	if node.Name, err = js.ParseIdent(p); err != nil {
		return
	}
	for {
		if typ := p.CurrentToken.Type; typ == token.GT || typ == selfClosingTag {
			break
		}
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
	if p.CurrentToken.Type == selfClosingTag {
		node.SelfClosing = true
		node.Layout.SelfClosingTag = p.CurrentToken
		p.AdvanceToken()
		return
	} else if _, err = p.Expect(token.GT); err != nil {
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
	return
}

// Plugin enriches the JavaScript parser, so that we can parse expressions that are not part of the JS standard.
func Plugin(b *plugin.Builder) {
	token.RegisterUnaryType(startTag)

	// now the parser can "scan" '<' and '</'
	b.UseScanner(func(sc *scanner.Scanner, next func() (token.Token, error)) (tok token.Token, err error) {
		if tok, err = next(); err != nil {
			return
		}
		switch tok.Type {
		case token.LT:
			c := sc.CurrentChar()
			switch {
			case scanner.IsLetter(c):
				tok.Type = startTag
			case c == '/':
				sc.AdvanceChar()
				tok.Type = endTag
				tok.Literal = "</"
			}
		case token.DIVIDE:
			if sc.CurrentChar() == '>' {
				sc.AdvanceChar()
				tok.Type = selfClosingTag
				tok.Literal = "/>"
			}
		}
		return
	})

	// now the parser can "parse" HTML tags
	b.UseUnaryParser(func(p *parser.Parser, next func() (ast.Expr, error)) (_ ast.Expr, err error) {
		if p.CurrentToken.Type == startTag {
			return parseTag(p)
		}
		return next()
	})
}
