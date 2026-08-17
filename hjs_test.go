package hjs_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xjslang/hjs"
	"github.com/xjslang/xjs/printer"
	"github.com/xorcare/golden"
)

func ExampleCompile() {
	// transform the input to AST
	input := `let p = <p>
		"Hello, "
		<span />
		<strong>"World!"</strong>
	</p>`
	result, err := hjs.Parse(input)
	if err != nil {
		panic(err)
	}

	// transform the AST to valid JS code
	jsCode, err := hjs.Compile(result)
	if err != nil {
		panic(err)
	}
	fmt.Println(jsCode)
	// Output: let p = (function(){const elem = document.createElement('p');elem.append("Hello, ");elem.append((function(){const elem = document.createElement('span');return elem})());elem.append((function(){const elem = document.createElement('strong');elem.append("World!");return elem})());return elem})();
}

func ExampleFormat() {
	// transform the input to AST
	input := `let p = <p>
		"Hello, "
		<span /><strong>
		"World!"
		</strong>
		</p>`
	result, err := hjs.Parse(input)
	if err != nil {
		panic(err)
	}

	// transform the AST to properly formatted XJS code
	xjsCode, err := hjs.Format(result)
	if err != nil {
		panic(err)
	}
	fmt.Println(xjsCode)
	// Output:
	// let p = <p>
	//   "Hello, "
	//   <span /> <strong>
	//     "World!"
	//   </strong>
	// </p>;
}

func TestParse(t *testing.T) {
	t.Run("errors", func(t *testing.T) {
		input := "<div>'Hello, World!'</p>"
		_, err := hjs.Parse(input)
		require.Error(t, err)
		require.ErrorContains(t, err, "expected closing tag </div>")
	})
}

func TestCompile(t *testing.T) {
	input := `
	let msg = 'Hello, World!'
	let handleClick = function () {
		console.log('Hello, Mars!')
	}
	let btn = <button type="button" onClick=handleClick>
		msg
	</button>
	
	let div1 = <div></div>
	let div2 = <div /> // self closing`
	result, err := hjs.Parse(input)
	require.NoError(t, err)
	code, err := hjs.Compile(result)
	require.NoError(t, err)
	golden.Assert(t, []byte(code))

	tests := []struct {
		name, input, expected string
	}{
		{"empty tag", `<p></p>`, `(function(){const elem = document.createElement('p');return elem})();`},
		{"self closing tag", `<div />`, `(function(){const elem = document.createElement('div');return elem})();`},
		{"dataset", `<span dataUserId='123'>"Hello, World!"</span>`, `(function(){const elem = document.createElement('span');elem.dataset.userId = '123';elem.append("Hello, World!");return elem})();`},
		{"handler", `<button onMouseDown=(handleClick)>"Click me!"</button>`, `(function(){const elem = document.createElement('button');elem.addEventListener('mousedown', (handleClick));elem.append("Click me!");return elem})();`},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := hjs.Parse(test.input)
			require.NoError(t, err)
			code, err := hjs.Compile(result)
			require.NoError(t, err)
			require.Equal(t, test.expected, code)
		})
	}
}

func TestFormat(t *testing.T) {
	input := `
	let msg = 'Hello, World!'
	let handleClick = function () {
		console.log('Hello, Mars!')
	}
	let btn = <button type="button" onClick=handleClick>
		msg
	</button>
	
	let div1 = <div></div>
	let div2 = <div /> // self closing`
	result, err := hjs.Parse(input)
	require.NoError(t, err)
	code, err := hjs.Format(result)
	require.NoError(t, err)
	golden.Assert(t, []byte(code))

	t.Run("preserve newlines", func(t *testing.T) {
		tests := []struct {
			input, expected string
		}{
			{"<p>'Hello, World!'</p>", "<p>'Hello, World!'</p>;"},
			{"<p>\n'Hello, World!'</p>", "<p>\n\t'Hello, World!'</p>;"},
			{"<p>\n'Hello, World!'\n</p>", "<p>\n\t'Hello, World!'\n</p>;"},
			{"<p>\n'Hello, '<b>'World!'</b>\n</p>", "<p>\n\t'Hello, ' <b>'World!'</b>\n</p>;"},
			{"<p>\n'Hello, '\n<b>'World!'</b>\n</p>", "<p>\n\t'Hello, '\n\t<b>'World!'</b>\n</p>;"},
		}
		for _, test := range tests {
			result, err := hjs.Parse(test.input)
			if !assert.NoError(t, err) {
				continue
			}
			code, err := hjs.Format(result, printer.WithIndent("\t"))
			if !assert.NoError(t, err) {
				continue
			}
			assert.Equal(t, test.expected, code)
		}
	})

	t.Run("empty tags", func(t *testing.T) {
		input := `let p = <p></p>`
		result, err := hjs.Parse(input)
		require.NoError(t, err)
		out, err := hjs.Format(result)
		require.NoError(t, err)
		require.Equal(t, "let p = <p></p>;", out)
	})

	t.Run("comments", func(t *testing.T) {
		input := `let p = <p>
		// c1
		"Hello, "
		<strong>
		/* c2 */
		"World!" // c3
		</strong>/* c4 */</p>
		
		let div = <div // c
		/>`
		result, err := hjs.Parse(input)
		require.NoError(t, err)

		// transform the AST to properly formatted code
		code, err := hjs.Format(result, printer.WithIndent("\t"))
		require.NoError(t, err)
		expectedCode := "let p = <p>\n\t// c1\n\t\"Hello, \"\n\t<strong>\n\t\t/* c2 */\n\t\t\"World!\" // c3\n\t</strong> /* c4 */</p>;\n\nlet div = <div // c\n/>;"
		require.Equal(t, expectedCode, code)
	})
}
