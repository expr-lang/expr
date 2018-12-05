package expr_test

import (
	"log"
	"reflect"
	"testing"
	"regexp"

	"github.com/etcinitd/expr"
)

type DepthFirstVisitAll struct {}

var visitor DepthFirstVisitAll

func (v *DepthFirstVisitAll) Nil() (interface{}, error) {
	return nil, nil
}

func (v *DepthFirstVisitAll) Value(value interface{}) (interface{}, error) {
	log.Println("Value visited:", value)
	return nil, nil
}

func (v *DepthFirstVisitAll) Name(name string) (interface{}, error) {
	log.Println("Name visited:", name)
	return nil, nil
}

func (v *DepthFirstVisitAll) Unary(operator string, n expr.Node) (interface{}, error) {
	return n.Visit(v)
}

func (v *DepthFirstVisitAll) Len(n expr.Node) (interface{}, error) {
	return n.Visit(v)
}

func (v *DepthFirstVisitAll) Binary(operator string, left expr.Node, right expr.Node) (interface{}, error) {
	return v.Array([]expr.Node{left, right})
}

func (v *DepthFirstVisitAll) Conditional(cond expr.Node, left expr.Node, right expr.Node) (interface{}, error) {
	return v.Array([]expr.Node{cond, left, right})
}

func (v *DepthFirstVisitAll) Matches(r *regexp.Regexp, left expr.Node, right expr.Node) (interface{}, error) {
	return v.Array([]expr.Node{left, right})
}

func (v *DepthFirstVisitAll) Index(left expr.Node, right expr.Node) (interface{}, error) {
	return v.Array([]expr.Node{left, right})
}

func (v *DepthFirstVisitAll) Property(n expr.Node, name string) (interface{}, error) {
	return n.Visit(v)
}

func (v *DepthFirstVisitAll) Method(n expr.Node, method string, args []expr.Node) (interface{}, error) {
	return v.Array(append([]expr.Node{n},args...))
}

func (v *DepthFirstVisitAll) Function(name string, args []expr.Node) (interface{}, error) {
	log.Println("Function visited:", name)
	return v.Array(args)
}

func (v *DepthFirstVisitAll) Array(nodes []expr.Node) (interface{}, error) {
	for _, n := range nodes {
		_, err := n.Visit(v)
		if err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func (v *DepthFirstVisitAll) Map(nodes []expr.PairNode) (interface{}, error) {
	for _, pair := range nodes {
		_, err := pair.Key.Visit(v)
		if err != nil {
			return nil, err
		}
		_, err = pair.Value.Visit(v)
		if err != nil {
			return nil, err
		}
	}
	return nil, nil
}

type visitTest struct {
	input    string
	expected interface{}
}

type visitErrorTest struct {
	input string
	err   string
}

var visitTests = []visitTest{
	{
		"t || (a == b.c + d[i] * f(k) && !false)",
		nil,
	},
}

func TestVisit(t *testing.T) {
	for _, test := range visitTests {
		root, err := expr.Parse(test.input)
		if err != nil {
			t.Errorf("%s:\n%v", test.input, err)
			continue
		}
		actual, err := expr.Visit(root, &visitor)
		if err != nil {
			t.Errorf("%s:\n%v", test.input, err)
			continue
		}
		if !reflect.DeepEqual(actual, test.expected) {
			t.Errorf("%s:\ngot\n\t%#v\nexpected\n\t%#v", test.input, actual, test.expected)
		}
	}
}
