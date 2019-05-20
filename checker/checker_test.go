package checker_test

import (
	"fmt"
	"github.com/antonmedv/expr/checker"
	"github.com/antonmedv/expr/internal/helper"
	"github.com/antonmedv/expr/parser"
	"testing"
)

func TestCheck(t *testing.T) {
	node, _ := parser.Parse("1 + foo")
	_, err := checker.Check(node, nil, helper.NewSource("1 + foo"))
	fmt.Printf("%v\n", err)
}
