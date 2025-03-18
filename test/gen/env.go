package main

var Env = map[string]any{
	"ok":        true,
	"i":         1,
	"str":       "str",
	"f64":       .5,
	"f32":       float32(.5),
	"i64":       int64(1),
	"i32":       int32(2),
	"u64":       uint64(4),
	"u32":       uint32(5),
	"u16":       uint16(6),
	"u8":        uint8(7),
	"array":     []int{1, 2, 3, 4, 5},
	"foo":       Foo{"foo"},
	"bar":       Bar{42, "bar"},
	"listOfFoo": []Foo{{"bar"}, {"baz"}},
	"add":       func(a, b int) int { return a + b },
	"div":       func(a, b int) int { return a / b },
	"half":      func(a float64) float64 { return a / 2 },
	"sumUp": func(a int, x ...int) int {
		s := a
		for _, n := range x {
			s += n
		}
		return s
	},
	"greet": func(name string) string { return "Hello, " + name },
}

type Foo struct {
	Bar string
}

func (f Foo) String() string {
	return "foo"
}

func (f Foo) Qux(s string) string {
	return f.Bar + s
}

type Bar struct {
	I   int    `expr:"i"`
	Str string `expr:"str"`
}
