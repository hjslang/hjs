> [!TIP]
> **Is it a library? Is it a framework? No, it's a programming language!**
> 
> This example is a simple web application, like one we could build manually with JavaScript and HTML, but instead we're using the HJS dialect, which allows us to use native HTML tags.

HJS (HTML for JavaScript) is an experimental JavaScript dialect with native support for HTML tags. The HJS compiler recognizes these tags and compiles them to standard JavaScript.

For example, the following code:
```jsx
export function button() {
	function handleClick() {
		alert('Button clicked!');
	}

    // **This is not JSX**, although it looks similar.
    // This is an extension of JavaScript expressions.
	return (
		<button type="button" onClick=handleClick> // {} is not necessary
			"Click me!" // strings must be quoted
		</button>
	);
}
```
is compiled to:
```js
export function button() {
  function handleClick() {
    alert("Button clicked!");
  }
  return (function () {
    const elem = document.createElement("button");
    elem.setAttribute("type", "button");
    elem.addEventListener("click", handleClick);
    elem.append("Click me!");
    return elem;
  })();
}
```

## Install and run
To run the application, you must install the `hjs` compiler:

```bash
# <repo-dir> is the repository directory, not the example directory
cd <repo-dir>
mage install
```

Once installed, run the `npm run dev` command:
```bash
cd example
npm run dev

Ready for changes
Change detected dist/src/main.js
Change detected dist/src/components/button.js
Change detected dist/src
Change detected dist/src/components
Serving "./dist" at http://127.0.0.1:8080
```

And that's it!
