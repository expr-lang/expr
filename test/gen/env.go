package main

var Env = map[string]any{
	"ok":    true,
	"i":     1,
	"str":   "str",
	"f64":   .5,
	"array": []int{1, 2, 3, 4, 5},
	"foo":   Foo{"foo"},
	"list":  []Foo{{"bar"}, {"baz"}},
	"add":   func(a, b int) int { return a + b },
	"greet": func(name string) string { return "Hello, " + name },
}

type Foo struct {
	Bar string
}

func (f Foo) String() string {
	return "foo"
}
