package html

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/js"
	"github.com/xjslang/xjs/jsextended"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/scanner"
	"github.com/xjslang/xjs/token"
)

var (
	START_TAG        = token.RegisterType("START_TAG", "<")
	END_TAG          = token.RegisterType("END_TAG", "</")
	SELF_CLOSING_TAG = token.RegisterType("SELF_CLOSING_TAG", "/>")
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
	if node.Layout.StartTag, err = p.Expect(START_TAG); err != nil {
		return
	}
	if node.Name, err = js.ParseIdent(p); err != nil {
		return
	}
	for {
		if typ := p.CurrentToken.Type; typ == token.GT || typ == SELF_CLOSING_TAG {
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
	if p.CurrentToken.Type == SELF_CLOSING_TAG {
		node.SelfClosing = true
		node.Layout.SelfClosingTag = p.CurrentToken
		p.AdvanceToken()
		return
	} else if _, err = p.Expect(token.GT); err != nil {
		return
	}
	for p.CurrentToken.Type != END_TAG {
		var child ast.Expr
		if child, err = p.ParseExpr(); err != nil {
			return
		}
		node.Children = append(node.Children, child)
	}
	if node.Layout.EndTag, err = p.Expect(END_TAG); err != nil {
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
func Parse(input string) (*js.Program, error) {
	token.RegisterPrefixOp(START_TAG)

	sb := jsextended.ScannerBuilder()
	// now the parser can "scan" '<' and '</'
	sb.UseScanner(func(sc *scanner.Scanner, next func(*scanner.Scanner) (token.Token, error)) (tok token.Token, err error) {
		if tok, err = next(sc); err != nil {
			return
		}
		switch tok.Type {
		case token.LT:
			c := sc.CurrentChar()
			switch {
			case scanner.IsLetter(c):
				tok.Type = START_TAG
			case c == '/':
				sc.AdvanceChar()
				tok.Type = END_TAG
				tok.Literal = "</"
			}
		case token.DIVIDE:
			if sc.CurrentChar() == '>' {
				sc.AdvanceChar()
				tok.Type = SELF_CLOSING_TAG
				tok.Literal = "/>"
			}
		}
		return
	})
	pb := jsextended.ParserBuilder()
	pb.UsePrefixOpParser(func(p *parser.Parser, next func(*parser.Parser) (ast.Expr, error)) (_ ast.Expr, err error) {
		if p.CurrentToken.Type == START_TAG {
			return parseTag(p)
		}
		return next(p)
	})
	s := sb.Build(input)
	p := pb.Build(s)
	return js.ParseProgram(p)
}
