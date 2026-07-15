<div align="center">

![HJS Logo](assets/logo.svg)

# HJS &ndash; HTML for JavaScript

A JavaScript dialect with native HTML support.

[Playground](https://xjslang.github.io/hjs/)

</div>

## How to use it

```go
package main

import (
	"fmt"

	"github.com/xjslang/hjs"
)

func main() {
	input := `let p = <p>
  "Hello, "
  <strong>"World!"</strong>
</p>`

	// transform the input to AST
	result, err := hjs.Parse([]byte(input))
	if err != nil {
		panic(err)
	}

	// transform the AST to valid JS code
	jsCode, err := hjs.Compile(result, hjs.WithRuntime())
	if err != nil {
		panic(err)
	}

	fmt.Println(jsCode)
	// Output:
	// let p = (function () {
	//   const elem = document.createElement("p");
	//   elem.append("Hello, ");
	//   elem.append(
	//     (function () {
	//       const elem = document.createElement("strong");
	//       elem.append("Word!");
	//       return elem;
	//     })(),
	//   );
	//   return elem;
	// })();
}
```

You'll find more examples in [./hjs_test.go](./hjs_test.go).

## Dev

**Architecture:**

```
./compiler.go  -- transform AST to valid JS code
./formatter.go -- transform AST to propertly formatted code
./plugin.go    -- enrich the XJS parser
./hjs.go       -- entry point
./hjs_test.go  -- integration tests
```

**Main mage commands:**

```bash
mage lint # check source code for errors
mage test # execute tests
mage -l   # complete list of commands
```
