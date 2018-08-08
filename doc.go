/*
Package expr is an engine that can evaluate expressions.

	// Evaluate expression on data.
	result, err := expr.Eval("expression", data)

	// Or precompile expression to ast first.
	node, err := expr.Parse("expression")

	// And run later.
	result, err := expr.Run(node, data)


Passing in Variables

You can pass variables into the expression, which can be of any valid Go type (including structs):

	// Maps
	data := map[string]interface{}{
		"Foo": ...
		"Bar": ...
	}

	// Structs
	data := Payload{
		Foo: ...
		Bar: ...
	}

	// Pass object
	result, err := expr.Eval("Foo == Bar", data)

Expr uses reflection for accessing and iterating passed data.
For example you can pass nested structures without any modification or preparation:

	type Cookie struct {
		Key   string
		Value string
	}
	type User struct {
		UserAgent string
		Cookies   []Cookie
	}
	type Request struct {
		User *user
	}

	req := Request{&User{
		Cookies:   []Cookie{{"origin", "www"}},
		UserAgent: "Firefox",
	}}

	ok, err := expr.Eval(`User.UserAgent matches "Firefox" and User.Cookies[0].Value == "www"`, req)


Passing in Functions

You can also pass functions into the expression:

	data := map[string]interface{}{
		"Request": req,
		"Values": func(xs []Cookie) []string {
			vs := make([]string, 0)
			for _, x := range xs {
				vs = append(vs, x.Value)
			}
			return vs
		},
	}

	ok, err := expr.Eval(`"www" in Values(Request.User.Cookies)`, data)


Parsing and caching

If you planning to execute some expression lots times, it's good to parse it first and only one time:

	// Parse expression to AST.
	ast, err := expr.Parse(expression)

	// Run given AST
	ok, err := expr.Run(ast, data)


Strict mode

Expr package support strict parse mode in which some type checks performed during parsing.
To parse expression in strict mode, define all of used variables:


	expression := `Request.User.UserAgent matches "Firefox"`
	node, err := expr.Parse(expression, expr.Define("Request", request{}))


Parse function will check used variables, accessed filed, logical operators and some other type checks.

If you try to use some undeclared variables, or access unknown field, an error will be returned during paring:

	expression := `Request.User.Cookies[0].Timestamp`
	node, err := expr.Parse(expression, expr.Define("Request", request{}))

	// err: Request.User.Cookies[0].Timestamp undefined (type expr_test.cookie has no field Timestamp)

Also it's possible to define all used variables and functions using expr.With and struct:

	type payload struct {
		Request *Request
		Values  func(xs []Cookie) []string
	}

	node, err := expr.Parse(expression, expr.With(payload{}))

Or with map:

	data := map[string]interface{}{
		"Request": req,
		"Values": func(xs []Cookie) []string {...},
	}

	node, err := expr.Parse(expression, expr.With(data))


Printing

Compiled ast can be compiled back to string expression using stringer fmt.Stringer interface:

	node, err := expr.Parse(expression)
	code := fmt.Sprintf("%v", node)


Number type

Inside Expr engine there is no distinguish between int, uint and float types (as in JavaScript).
All numbers inside Expr engine represented as `float64`.
You should remember about it if you use any of binary operators (`+`, `-`, `/`, `*`, etc).
Otherwise type remain unchanged.

	data := map[string]int{
		"Foo": 1,
		"Bar": 2,
	}

	out, err := expr.Eval(`Foo`, data) // int

	out, err := expr.Eval(`Foo + Bar`, data) // float64


*/
package expr
