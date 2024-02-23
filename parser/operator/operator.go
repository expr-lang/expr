package operator

type Associativity int

const (
	Left Associativity = iota + 1
	Right
)

type Operator struct {
	Precedence    int
	Associativity Associativity
	Method        bool
}

func Less(a, b string) bool {
	return Binary[a].Precedence < Binary[b].Precedence
}

func IsBoolean(op string) bool {
	return op == "and" || op == "or" || op == "&&" || op == "||"
}

var Unary = map[string]Operator{
	"not": {50, Left, false},
	"!":   {50, Left, false},
	"-":   {90, Left, false},
	"+":   {90, Left, false},
}

var Binary = map[string]Operator{
	"|":          {0, Left, false},
	"or":         {10, Left, false},
	"||":         {10, Left, false},
	"and":        {15, Left, false},
	"&&":         {15, Left, false},
	"==":         {20, Left, false},
	"!=":         {20, Left, false},
	"<":          {20, Left, false},
	">":          {20, Left, false},
	">=":         {20, Left, false},
	"<=":         {20, Left, false},
	"in":         {20, Left, true},
	"matches":    {20, Left, true},
	"contains":   {20, Left, true},
	"startsWith": {20, Left, true},
	"endsWith":   {20, Left, true},
	"..":         {25, Left, false},
	"+":          {30, Left, false},
	"-":          {30, Left, false},
	"*":          {60, Left, false},
	"/":          {60, Left, false},
	"%":          {60, Left, false},
	"**":         {100, Right, false},
	"^":          {100, Right, false},
	"??":         {500, Left, false},
}
