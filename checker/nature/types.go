package nature

type SubType interface {
	String() string
}

type Array struct {
	Of Nature
}

func (a Array) String() string {
	return "[]" + a.Of.String()
}
