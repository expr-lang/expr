package fuzz

func NewEnv() map[string]any {
	return map[string]any{
		"ok":    true,
		"f64":   .5,
		"f32":   float32(.5),
		"i":     1,
		"i64":   int64(1),
		"i32":   int32(1),
		"array": []int{1, 2, 3, 4, 5},
		"list":  []Foo{{"bar"}, {"baz"}},
		"foo":   Foo{"bar"},
		"add":   func(a, b int) int { return a + b },
		"div":   func(a, b int) int { return a / b },
		"half":  func(a float64) float64 { return a / 2 },
		"score": func(a int, x ...int) int {
			s := a
			for _, n := range x {
				s += n
			}
			return s
		},
		"greet": func(name string) string { return "Hello, " + name },
	}
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
