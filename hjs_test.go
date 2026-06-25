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
		<strong>"World!"</strong>
	</p>`
	result, err := hjs.Parse([]byte(input))
	if err != nil {
		panic(err)
	}

	// transform the AST to valid JS code
	jsCode, err := hjs.Compile(result)
	if err != nil {
		panic(err)
	}
	fmt.Println(jsCode)
	// Output: let p = (function(){const elem = document.createElement('p');elem.append("Hello, ");elem.append((function(){const elem = document.createElement('strong');elem.append("World!");return elem})());return elem})();
}

func ExampleFormat() {
	// transform the input to AST
	input := `let p = <p>
		"Hello, "<strong>
		"World!"
		</strong>
		</p>`
	result, err := hjs.Parse([]byte(input))
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
	//   "Hello, " <strong>
	//     "World!"
	//   </strong>
	// </p>;
}

func TestParse(t *testing.T) {
	t.Run("errors", func(t *testing.T) {
		input := "<div>'Hello, World!'</p>"
		_, err := hjs.Parse([]byte(input))
		require.Error(t, err)
		require.Equal(t, err.Error(), "[line:0, col:22] expected closing tag </div>")
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
	</button>`
	result, err := hjs.Parse([]byte(input))
	require.NoError(t, err)
	code, err := hjs.Compile(result)
	require.NoError(t, err)
	golden.Assert(t, []byte(code))

	t.Run("empty tag", func(t *testing.T) {
		input := `let p = <p></p>`
		result, err := hjs.Parse([]byte(input))
		require.NoError(t, err)
		out, err := hjs.Compile(result)
		require.NoError(t, err)
		require.Equal(t, "let p = (function(){const elem = document.createElement('p');return elem})();", out)
	})
}

func TestFormat(t *testing.T) {
	input := `
	let msg = 'Hello, World!'
	let handleClick = function () {
		console.log('Hello, Mars!')
	}
	let btn = <button type="button" onClick=handleClick>
		msg
	</button>`
	result, err := hjs.Parse([]byte(input))
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
			result, err := hjs.Parse([]byte(test.input))
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
		result, err := hjs.Parse([]byte(input))
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
		</strong>/* c4 */</p>`
		result, err := hjs.Parse([]byte(input))
		require.NoError(t, err)

		// transform the AST to properly formatted code
		code, err := hjs.Format(result, printer.WithIndent("\t"))
		require.NoError(t, err)
		expectedCode := "let p = <p>\n\t// c1\n\t\"Hello, \"\n\t<strong>\n\t\t/* c2 */\n\t\t\"World!\" // c3\n\t</strong> /* c4 */</p>;"
		require.Equal(t, expectedCode, code)
	})
}
