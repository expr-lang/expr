package checker

import (
	"fmt"
	"reflect"
	"regexp"
	"time"

	"github.com/expr-lang/expr/ast"
	"github.com/expr-lang/expr/builtin"
	. "github.com/expr-lang/expr/checker/nature"
	"github.com/expr-lang/expr/conf"
	"github.com/expr-lang/expr/file"
	"github.com/expr-lang/expr/parser"
)

var (
	anyType = reflect.TypeOf(new(any)).Elem()

	unknown        = Nature{}
	nilNature      = Nature{Nil: true}
	boolNature     = Nature{Type: reflect.TypeOf(true)}
	integerNature  = Nature{Type: reflect.TypeOf(0)}
	floatNature    = Nature{Type: reflect.TypeOf(float64(0))}
	stringNature   = Nature{Type: reflect.TypeOf("")}
	arrayNature    = Nature{Type: reflect.TypeOf([]any{})}
	mapNature      = Nature{Type: reflect.TypeOf(map[string]any{})}
	timeNature     = Nature{Type: reflect.TypeOf(time.Time{})}
	durationNature = Nature{Type: reflect.TypeOf(time.Duration(0))}
)

// ParseCheck parses input expression and checks its types. Also, it applies
// all provided patchers. In case of error, it returns error with a tree.
func ParseCheck(input string, config *conf.Config) (*parser.Tree, error) {
	tree, err := parser.ParseWithConfig(input, config)
	if err != nil {
		return tree, err
	}

	_, err = new(Checker).PatchAndCheck(tree, config)
	if err != nil {
		return tree, err
	}

	return tree, nil
}

// Check calls Check on a disposable Checker.
func Check(tree *parser.Tree, config *conf.Config) (reflect.Type, error) {
	return new(Checker).Check(tree, config)
}

type Checker struct {
	config          *conf.Config
	predicateScopes []predicateScope
	varScopes       []varScope
	err             *file.Error
	needsReset      bool
}

type predicateScope struct {
	collection Nature
	vars       []varScope
}

type varScope struct {
	name   string
	nature Nature
}

// PatchAndCheck applies all patchers and checks the tree.
func (c *Checker) PatchAndCheck(tree *parser.Tree, config *conf.Config) (reflect.Type, error) {
	c.reset(config)
	if len(config.Visitors) > 0 {
		// Run all patchers that dont support being run repeatedly first
		c.runVisitors(tree, false)

		// Run patchers that require multiple passes next (currently only Operator patching)
		c.runVisitors(tree, true)
	}
	return c.Check(tree, config)
}

// Check checks types of the expression tree. It returns type of the expression
// and error if any. If config is nil, then default configuration will be used.
func (c *Checker) Check(tree *parser.Tree, config *conf.Config) (reflect.Type, error) {
	c.reset(config)
	return c.check(tree)
}

// Run visitors in a given config over the given tree
// runRepeatable controls whether to filter for only vistors that require multiple passes or not
func (c *Checker) runVisitors(tree *parser.Tree, runRepeatable bool) {
	for {
		more := false
		for _, v := range c.config.Visitors {
			// We need to perform types check, because some visitors may rely on
			// types information available in the tree.
			_, _ = c.Check(tree, c.config)

			r, repeatable := v.(interface {
				Reset()
				ShouldRepeat() bool
			})

			if repeatable {
				if runRepeatable {
					r.Reset()
					ast.Walk(&tree.Node, v)
					more = more || r.ShouldRepeat()
				}
			} else {
				if !runRepeatable {
					ast.Walk(&tree.Node, v)
				}
			}
		}

		if !more {
			break
		}
	}
}

func (c *Checker) check(tree *parser.Tree) (reflect.Type, error) {
	nt := c.visit(tree.Node)

	// To keep compatibility with previous versions, we should return any, if nature is unknown.
	t := nt.Type
	if t == nil {
		t = anyType
	}

	if c.err != nil {
		return t, c.err.Bind(tree.Source)
	}

	if c.config.Expect != reflect.Invalid {
		if c.config.ExpectAny {
			if nt.IsUnknown() {
				return t, nil
			}
		}

		switch c.config.Expect {
		case reflect.Int, reflect.Int64, reflect.Float64:
			if !nt.IsNumber() {
				return nil, fmt.Errorf("expected %v, but got %s", c.config.Expect, nt.String())
			}
		default:
			if nt.Kind() != c.config.Expect {
				return nil, fmt.Errorf("expected %v, but got %s", c.config.Expect, nt.String())
			}
		}
	}

	return t, nil
}

func (c *Checker) reset(config *conf.Config) {
	if c.needsReset {
		clearSlice(c.predicateScopes)
		clearSlice(c.varScopes)
		c.predicateScopes = c.predicateScopes[:0]
		c.varScopes = c.varScopes[:0]
		c.err = nil
	}
	c.needsReset = true

	if config == nil {
		config = conf.New(nil)
	}
	c.config = config
}

func clearSlice[S ~[]E, E any](s S) {
	var zero E
	for i := range s {
		s[i] = zero
	}
}

func (v *Checker) visit(node ast.Node) Nature {
	var nt Nature
	switch n := node.(type) {
	case *ast.NilNode:
		nt = v.nilNode(n)
	case *ast.IdentifierNode:
		nt = v.identifierNode(n)
	case *ast.IntegerNode:
		nt = v.integerNode(n)
	case *ast.FloatNode:
		nt = v.floatNode(n)
	case *ast.BoolNode:
		nt = v.boolNode(n)
	case *ast.StringNode:
		nt = v.stringNode(n)
	case *ast.ConstantNode:
		nt = v.constantNode(n)
	case *ast.UnaryNode:
		nt = v.unaryNode(n)
	case *ast.BinaryNode:
		nt = v.binaryNode(n)
	case *ast.ChainNode:
		nt = v.chainNode(n)
	case *ast.MemberNode:
		nt = v.memberNode(n)
	case *ast.SliceNode:
		nt = v.sliceNode(n)
	case *ast.CallNode:
		nt = v.callNode(n)
	case *ast.BuiltinNode:
		nt = v.builtinNode(n)
	case *ast.PredicateNode:
		nt = v.predicateNode(n)
	case *ast.PointerNode:
		nt = v.pointerNode(n)
	case *ast.VariableDeclaratorNode:
		nt = v.variableDeclaratorNode(n)
	case *ast.SequenceNode:
		nt = v.sequenceNode(n)
	case *ast.ConditionalNode:
		nt = v.conditionalNode(n)
	case *ast.ArrayNode:
		nt = v.arrayNode(n)
	case *ast.MapNode:
		nt = v.mapNode(n)
	case *ast.PairNode:
		nt = v.pairNode(n)
	default:
		panic(fmt.Sprintf("undefined node type (%T)", node))
	}
	node.SetNature(nt)
	return nt
}

func (v *Checker) error(node ast.Node, format string, args ...any) Nature {
	if v.err == nil { // show first error
		v.err = &file.Error{
			Location: node.Location(),
			Message:  fmt.Sprintf(format, args...),
		}
	}
	return unknown
}

func (v *Checker) nilNode(*ast.NilNode) Nature {
	return nilNature
}

func (v *Checker) identifierNode(node *ast.IdentifierNode) Nature {
	if variable, ok := v.lookupVariable(node.Value); ok {
		return variable.nature
	}
	if node.Value == "$env" {
		return unknown
	}

	return v.ident(node, node.Value, v.config.Strict, true)
}

// ident method returns type of environment variable, builtin or function.
func (v *Checker) ident(node ast.Node, name string, strict, builtins bool) Nature {
	if nt, ok := v.config.Env.Get(&v.config.NtCache, name); ok {
		return nt
	}
	if builtins {
		if fn, ok := v.config.Functions[name]; ok {
			return Nature{Type: fn.Type(), Func: fn}
		}
		if fn, ok := v.config.Builtins[name]; ok {
			return Nature{Type: fn.Type(), Func: fn}
		}
	}
	if v.config.Strict && strict {
		return v.error(node, "unknown name %s", name)
	}
	return unknown
}

func (v *Checker) integerNode(*ast.IntegerNode) Nature {
	return integerNature
}

func (v *Checker) floatNode(*ast.FloatNode) Nature {
	return floatNature
}

func (v *Checker) boolNode(*ast.BoolNode) Nature {
	return boolNature
}

func (v *Checker) stringNode(*ast.StringNode) Nature {
	return stringNature
}

func (v *Checker) constantNode(node *ast.ConstantNode) Nature {
	return Nature{Type: reflect.TypeOf(node.Value)}
}

func (v *Checker) unaryNode(node *ast.UnaryNode) Nature {
	nt := v.visit(node.Node)
	nt = nt.Deref()

	switch node.Operator {

	case "!", "not":
		if nt.IsBool() {
			return boolNature
		}
		if nt.IsUnknown() {
			return boolNature
		}

	case "+", "-":
		if nt.IsNumber() {
			return nt
		}
		if nt.IsUnknown() {
			return unknown
		}

	default:
		return v.error(node, "unknown operator (%s)", node.Operator)
	}

	return v.error(node, `invalid operation: %s (mismatched type %s)`, node.Operator, nt.String())
}

func (v *Checker) binaryNode(node *ast.BinaryNode) Nature {
	l := v.visit(node.Left)
	r := v.visit(node.Right)

	l = l.Deref()
	r = r.Deref()

	switch node.Operator {
	case "==", "!=":
		if l.ComparableTo(r) {
			return boolNature
		}

	case "or", "||", "and", "&&":
		if l.IsBool() && r.IsBool() {
			return boolNature
		}
		if l.MaybeCompatible(r, BoolCheck) {
			return boolNature
		}

	case "<", ">", ">=", "<=":
		if l.IsNumber() && r.IsNumber() {
			return boolNature
		}
		if l.IsString() && r.IsString() {
			return boolNature
		}
		if l.IsTime() && r.IsTime() {
			return boolNature
		}
		if l.IsDuration() && r.IsDuration() {
			return boolNature
		}
		if l.MaybeCompatible(r, NumberCheck, StringCheck, TimeCheck, DurationCheck) {
			return boolNature
		}

	case "-":
		if l.IsNumber() && r.IsNumber() {
			return l.PromoteNumericNature(r)
		}
		if l.IsTime() && r.IsTime() {
			return durationNature
		}
		if l.IsTime() && r.IsDuration() {
			return timeNature
		}
		if l.IsDuration() && r.IsDuration() {
			return durationNature
		}
		if l.MaybeCompatible(r, NumberCheck, TimeCheck, DurationCheck) {
			return unknown
		}

	case "*":
		if l.IsNumber() && r.IsNumber() {
			return l.PromoteNumericNature(r)
		}
		if l.IsNumber() && r.IsDuration() {
			return durationNature
		}
		if l.IsDuration() && r.IsNumber() {
			return durationNature
		}
		if l.IsDuration() && r.IsDuration() {
			return durationNature
		}
		if l.MaybeCompatible(r, NumberCheck, DurationCheck) {
			return unknown
		}

	case "/":
		if l.IsNumber() && r.IsNumber() {
			return floatNature
		}
		if l.MaybeCompatible(r, NumberCheck) {
			return floatNature
		}

	case "**", "^":
		if l.IsNumber() && r.IsNumber() {
			return floatNature
		}
		if l.MaybeCompatible(r, NumberCheck) {
			return floatNature
		}

	case "%":
		if l.IsInteger() && r.IsInteger() {
			return integerNature
		}
		if l.MaybeCompatible(r, IntegerCheck) {
			return integerNature
		}

	case "+":
		if l.IsNumber() && r.IsNumber() {
			return l.PromoteNumericNature(r)
		}
		if l.IsString() && r.IsString() {
			return stringNature
		}
		if l.IsTime() && r.IsDuration() {
			return timeNature
		}
		if l.IsDuration() && r.IsTime() {
			return timeNature
		}
		if l.IsDuration() && r.IsDuration() {
			return durationNature
		}
		if l.MaybeCompatible(r, NumberCheck, StringCheck, TimeCheck, DurationCheck) {
			return unknown
		}

	case "in":
		if (l.IsString() || l.IsUnknown()) && r.IsStruct() {
			return boolNature
		}
		if r.IsMap() {
			rKey := r.Key()
			if !l.IsUnknown() && !l.AssignableTo(rKey) {
				return v.error(node, "cannot use %s as type %s in map key", l.String(), rKey.String())
			}
			return boolNature
		}
		if r.IsArray() {
			rElem := r.Elem()
			if !l.ComparableTo(rElem) {
				return v.error(node, "cannot use %s as type %s in array", l.String(), rElem.String())
			}
			return boolNature
		}
		if l.IsUnknown() && r.IsAnyOf(StringCheck, ArrayCheck, MapCheck) {
			return boolNature
		}
		if r.IsUnknown() {
			return boolNature
		}

	case "matches":
		if s, ok := node.Right.(*ast.StringNode); ok {
			_, err := regexp.Compile(s.Value)
			if err != nil {
				return v.error(node, err.Error())
			}
		}
		if l.IsString() && r.IsString() {
			return boolNature
		}
		if l.MaybeCompatible(r, StringCheck) {
			return boolNature
		}

	case "contains", "startsWith", "endsWith":
		if l.IsString() && r.IsString() {
			return boolNature
		}
		if l.MaybeCompatible(r, StringCheck) {
			return boolNature
		}

	case "..":
		if l.IsInteger() && r.IsInteger() {
			return integerNature.MakeArrayOf()
		}
		if l.MaybeCompatible(r, IntegerCheck) {
			return integerNature.MakeArrayOf()
		}

	case "??":
		if l.Nil && !r.Nil {
			return r
		}
		if !l.Nil && r.Nil {
			return l
		}
		if l.Nil && r.Nil {
			return nilNature
		}
		if r.AssignableTo(l) {
			return l
		}
		return unknown

	default:
		return v.error(node, "unknown operator (%s)", node.Operator)

	}

	return v.error(node, `invalid operation: %s (mismatched types %s and %s)`, node.Operator, l.String(), r.String())
}

func (v *Checker) chainNode(node *ast.ChainNode) Nature {
	return v.visit(node.Node)
}

func (v *Checker) memberNode(node *ast.MemberNode) Nature {
	// $env variable
	if an, ok := node.Node.(*ast.IdentifierNode); ok && an.Value == "$env" {
		if name, ok := node.Property.(*ast.StringNode); ok {
			strict := v.config.Strict
			if node.Optional {
				// If user explicitly set optional flag, then we should not
				// throw error if field is not found (as user trying to handle
				// this case). But if user did not set optional flag, then we
				// should throw error if field is not found & v.config.Strict.
				strict = false
			}
			return v.ident(node, name.Value, strict, false /* no builtins and no functions */)
		}
		return unknown
	}

	base := v.visit(node.Node)
	prop := v.visit(node.Property)

	if base.IsUnknown() {
		return unknown
	}

	if name, ok := node.Property.(*ast.StringNode); ok {
		if base.Nil {
			return v.error(node, "type nil has no field %s", name.Value)
		}

		// First, check methods defined on base type itself,
		// independent of which type it is. Without dereferencing.
		if m, ok := base.MethodByName(&v.config.NtCache, name.Value); ok {
			return m
		}
	}

	base = base.Deref()

	switch base.Kind() {
	case reflect.Map:
		if !prop.AssignableTo(base.Key()) && !prop.IsUnknown() {
			return v.error(node.Property, "cannot use %s to get an element from %s", prop.String(), base.String())
		}
		if prop, ok := node.Property.(*ast.StringNode); ok && base.MapData != nil {
			if field, ok := base.Fields[prop.Value]; ok {
				return field
			} else if base.Strict {
				return v.error(node.Property, "unknown field %s", prop.Value)
			}
		}
		return base.Elem()

	case reflect.Array, reflect.Slice:
		if !prop.IsInteger() && !prop.IsUnknown() {
			return v.error(node.Property, "array elements can only be selected using an integer (got %s)", prop.String())
		}
		return base.Elem()

	case reflect.Struct:
		if name, ok := node.Property.(*ast.StringNode); ok {
			propertyName := name.Value
			if field, ok := base.FieldByName(&v.config.NtCache, propertyName); ok {
				return Nature{Type: field.Type}
			}
			if node.Method {
				return v.error(node, "type %v has no method %v", base.String(), propertyName)
			}
			return v.error(node, "type %v has no field %v", base.String(), propertyName)
		}
	}

	// Not found.

	if name, ok := node.Property.(*ast.StringNode); ok {
		if node.Method {
			return v.error(node, "type %v has no method %v", base.String(), name.Value)
		}
		return v.error(node, "type %v has no field %v", base.String(), name.Value)
	}
	return v.error(node, "type %v[%v] is undefined", base.String(), prop.String())
}

func (v *Checker) sliceNode(node *ast.SliceNode) Nature {
	nt := v.visit(node.Node)

	if nt.IsUnknown() {
		return unknown
	}

	switch nt.Kind() {
	case reflect.String, reflect.Array, reflect.Slice:
		// ok
	default:
		return v.error(node, "cannot slice %s", nt.String())
	}

	if node.From != nil {
		from := v.visit(node.From)
		if !from.IsInteger() && !from.IsUnknown() {
			return v.error(node.From, "non-integer slice index %v", from.String())
		}
	}

	if node.To != nil {
		to := v.visit(node.To)
		if !to.IsInteger() && !to.IsUnknown() {
			return v.error(node.To, "non-integer slice index %v", to.String())
		}
	}

	return nt
}

func (v *Checker) callNode(node *ast.CallNode) Nature {
	// Check if type was set on node (for example, by patcher)
	// and use node type instead of function return type.
	//
	// If node type is anyType, then we should use function
	// return type. For example, on error we return anyType
	// for a call `errCall().Method()` and method will be
	// evaluated on `anyType.Method()`, so return type will
	// be anyType `anyType.Method(): anyType`. Patcher can
	// fix `errCall()` to return proper type, so on second
	// checker pass we should replace anyType on method node
	// with new correct function return type.
	if typ := node.Type(); typ != nil && typ != anyType {
		return node.Nature()
	}

	return v.functionReturnType(node)
}

func (v *Checker) functionReturnType(node *ast.CallNode) Nature {
	nt := v.visit(node.Callee)
	if nt.IsUnknown() {
		return unknown
	}

	if nt.Func != nil {
		return v.checkFunction(nt.Func, node, node.Arguments)
	}

	fnName := "function"
	if identifier, ok := node.Callee.(*ast.IdentifierNode); ok {
		fnName = identifier.Value
	}
	if member, ok := node.Callee.(*ast.MemberNode); ok {
		if name, ok := member.Property.(*ast.StringNode); ok {
			fnName = name.Value
		}
	}

	if nt.Nil {
		return v.error(node, "%v is nil; cannot call nil as function", fnName)
	}

	if nt.Kind() == reflect.Func {
		outType, err := v.checkArguments(fnName, nt, node.Arguments, node)
		if err != nil {
			if v.err == nil {
				v.err = err
			}
			return unknown
		}
		return outType
	}
	return v.error(node, "%s is not callable", nt.String())
}

func (v *Checker) builtinNode(node *ast.BuiltinNode) Nature {
	switch node.Name {
	case "all", "none", "any", "one":
		collection := v.visit(node.Arguments[0])
		collection = collection.Deref()
		if !collection.IsArray() && !collection.IsUnknown() {
			return v.error(node.Arguments[0], "builtin %v takes only array (got %v)", node.Name, collection.String())
		}

		v.begin(collection)
		predicate := v.visit(node.Arguments[1])
		v.end()

		if predicate.IsFunc() &&
			predicate.NumOut() == 1 &&
			predicate.NumIn() == 1 && predicate.In(0).IsUnknown() {

			if !predicate.Out(0).IsBool() && !predicate.Out(0).IsUnknown() {
				return v.error(node.Arguments[1], "predicate should return boolean (got %s)", predicate.Out(0))
			}
			return boolNature
		}
		return v.error(node.Arguments[1], "predicate should has one input and one output param")

	case "filter":
		collection := v.visit(node.Arguments[0])
		collection = collection.Deref()
		if !collection.IsArray() && !collection.IsUnknown() {
			return v.error(node.Arguments[0], "builtin %v takes only array (got %v)", node.Name, collection.String())
		}

		v.begin(collection)
		predicate := v.visit(node.Arguments[1])
		v.end()

		if predicate.IsFunc() &&
			predicate.NumOut() == 1 &&
			predicate.NumIn() == 1 && predicate.In(0).IsUnknown() {

			if !predicate.Out(0).IsBool() && !predicate.Out(0).IsUnknown() {
				return v.error(node.Arguments[1], "predicate should return boolean (got %s)", predicate.Out(0))
			}
			if collection.IsUnknown() {
				return arrayNature
			}
			collection = collection.Elem()
			return collection.MakeArrayOf()
		}
		return v.error(node.Arguments[1], "predicate should has one input and one output param")

	case "map":
		collection := v.visit(node.Arguments[0])
		collection = collection.Deref()
		if !collection.IsArray() && !collection.IsUnknown() {
			return v.error(node.Arguments[0], "builtin %v takes only array (got %v)", node.Name, collection.String())
		}

		v.begin(collection, varScope{"index", integerNature})
		predicate := v.visit(node.Arguments[1])
		v.end()

		if predicate.IsFunc() &&
			predicate.NumOut() == 1 &&
			predicate.NumIn() == 1 && predicate.In(0).IsUnknown() {

			return predicate.PredicateOut.MakeArrayOf()
		}
		return v.error(node.Arguments[1], "predicate should has one input and one output param")

	case "count":
		collection := v.visit(node.Arguments[0])
		collection = collection.Deref()
		if !collection.IsArray() && !collection.IsUnknown() {
			return v.error(node.Arguments[0], "builtin %v takes only array (got %v)", node.Name, collection.String())
		}

		if len(node.Arguments) == 1 {
			return integerNature
		}

		v.begin(collection)
		predicate := v.visit(node.Arguments[1])
		v.end()

		if predicate.IsFunc() &&
			predicate.NumOut() == 1 &&
			predicate.NumIn() == 1 && predicate.In(0).IsUnknown() {
			if !predicate.Out(0).IsBool() && !predicate.Out(0).IsUnknown() {
				return v.error(node.Arguments[1], "predicate should return boolean (got %s)", predicate.Out(0))
			}

			return integerNature
		}
		return v.error(node.Arguments[1], "predicate should has one input and one output param")

	case "sum":
		collection := v.visit(node.Arguments[0])
		collection = collection.Deref()
		if !collection.IsArray() && !collection.IsUnknown() {
			return v.error(node.Arguments[0], "builtin %v takes only array (got %v)", node.Name, collection.String())
		}

		if len(node.Arguments) == 2 {
			v.begin(collection)
			predicate := v.visit(node.Arguments[1])
			v.end()

			if predicate.IsFunc() &&
				predicate.NumOut() == 1 &&
				predicate.NumIn() == 1 && predicate.In(0).IsUnknown() {
				return *predicate.Out(0)
			}
		} else {
			if collection.IsUnknown() {
				return unknown
			}
			return collection.Elem()
		}

	case "find", "findLast":
		collection := v.visit(node.Arguments[0])
		collection = collection.Deref()
		if !collection.IsArray() && !collection.IsUnknown() {
			return v.error(node.Arguments[0], "builtin %v takes only array (got %v)", node.Name, collection.String())
		}

		v.begin(collection)
		predicate := v.visit(node.Arguments[1])
		v.end()

		if predicate.IsFunc() &&
			predicate.NumOut() == 1 &&
			predicate.NumIn() == 1 && predicate.In(0).IsUnknown() {

			if !predicate.Out(0).IsBool() && !predicate.Out(0).IsUnknown() {
				return v.error(node.Arguments[1], "predicate should return boolean (got %s)", predicate.Out(0))
			}
			if collection.IsUnknown() {
				return unknown
			}
			return collection.Elem()
		}
		return v.error(node.Arguments[1], "predicate should has one input and one output param")

	case "findIndex", "findLastIndex":
		collection := v.visit(node.Arguments[0])
		collection = collection.Deref()
		if !collection.IsArray() && !collection.IsUnknown() {
			return v.error(node.Arguments[0], "builtin %v takes only array (got %v)", node.Name, collection.String())
		}

		v.begin(collection)
		predicate := v.visit(node.Arguments[1])
		v.end()

		if predicate.IsFunc() &&
			predicate.NumOut() == 1 &&
			predicate.NumIn() == 1 && predicate.In(0).IsUnknown() {

			if !predicate.Out(0).IsBool() && !predicate.Out(0).IsUnknown() {
				return v.error(node.Arguments[1], "predicate should return boolean (got %s)", predicate.Out(0))
			}
			return integerNature
		}
		return v.error(node.Arguments[1], "predicate should has one input and one output param")

	case "groupBy":
		collection := v.visit(node.Arguments[0])
		collection = collection.Deref()
		if !collection.IsArray() && !collection.IsUnknown() {
			return v.error(node.Arguments[0], "builtin %v takes only array (got %v)", node.Name, collection.String())
		}

		v.begin(collection)
		predicate := v.visit(node.Arguments[1])
		v.end()

		if predicate.IsFunc() &&
			predicate.NumOut() == 1 &&
			predicate.NumIn() == 1 && predicate.In(0).IsUnknown() {

			collection = collection.Elem()
			collection = collection.MakeArrayOf()
			return Nature{Type: reflect.TypeOf(map[any][]any{}), ArrayOf: &collection}
		}
		return v.error(node.Arguments[1], "predicate should has one input and one output param")

	case "sortBy":
		collection := v.visit(node.Arguments[0])
		collection = collection.Deref()
		if !collection.IsArray() && !collection.IsUnknown() {
			return v.error(node.Arguments[0], "builtin %v takes only array (got %v)", node.Name, collection.String())
		}

		v.begin(collection)
		predicate := v.visit(node.Arguments[1])
		v.end()

		if len(node.Arguments) == 3 {
			_ = v.visit(node.Arguments[2])
		}

		if predicate.IsFunc() &&
			predicate.NumOut() == 1 &&
			predicate.NumIn() == 1 && predicate.In(0).IsUnknown() {

			return collection
		}
		return v.error(node.Arguments[1], "predicate should has one input and one output param")

	case "reduce":
		collection := v.visit(node.Arguments[0])
		collection = collection.Deref()
		if !collection.IsArray() && !collection.IsUnknown() {
			return v.error(node.Arguments[0], "builtin %v takes only array (got %v)", node.Name, collection.String())
		}

		v.begin(collection, varScope{"index", integerNature}, varScope{"acc", unknown})
		predicate := v.visit(node.Arguments[1])
		v.end()

		if len(node.Arguments) == 3 {
			_ = v.visit(node.Arguments[2])
		}

		if predicate.IsFunc() && predicate.NumOut() == 1 {
			return *predicate.PredicateOut
		}
		return v.error(node.Arguments[1], "predicate should has two input and one output param")

	}

	if id, ok := builtin.Index[node.Name]; ok {
		switch node.Name {
		case "get":
			return v.checkBuiltinGet(node)
		}
		return v.checkFunction(builtin.Builtins[id], node, node.Arguments)
	}

	return v.error(node, "unknown builtin %v", node.Name)
}

func (v *Checker) begin(collectionNature Nature, vars ...varScope) {
	v.predicateScopes = append(v.predicateScopes, predicateScope{
		collection: collectionNature,
		vars:       vars,
	})
}

func (v *Checker) end() {
	v.predicateScopes = v.predicateScopes[:len(v.predicateScopes)-1]
}

func (v *Checker) checkBuiltinGet(node *ast.BuiltinNode) Nature {
	if len(node.Arguments) != 2 {
		return v.error(node, "invalid number of arguments (expected 2, got %d)", len(node.Arguments))
	}

	base := v.visit(node.Arguments[0])
	prop := v.visit(node.Arguments[1])

	if id, ok := node.Arguments[0].(*ast.IdentifierNode); ok && id.Value == "$env" {
		if s, ok := node.Arguments[1].(*ast.StringNode); ok {
			if nt, ok := v.config.Env.Get(&v.config.NtCache, s.Value); ok {
				return nt
			}
		}
		return unknown
	}

	if base.IsUnknown() {
		return unknown
	}

	switch base.Kind() {
	case reflect.Slice, reflect.Array:
		if !prop.IsInteger() && !prop.IsUnknown() {
			return v.error(node.Arguments[1], "non-integer slice index %s", prop.String())
		}
		return base.Elem()
	case reflect.Map:
		if !prop.AssignableTo(base.Key()) && !prop.IsUnknown() {
			return v.error(node.Arguments[1], "cannot use %s to get an element from %s", prop.String(), base.String())
		}
		return base.Elem()
	}
	return v.error(node.Arguments[0], "type %v does not support indexing", base.String())
}

func (v *Checker) checkFunction(f *builtin.Function, node ast.Node, arguments []ast.Node) Nature {
	if f.Validate != nil {
		args := make([]reflect.Type, len(arguments))
		for i, arg := range arguments {
			argNature := v.visit(arg)
			if argNature.IsUnknown() {
				args[i] = anyType
			} else {
				args[i] = argNature.Type
			}
		}
		t, err := f.Validate(args)
		if err != nil {
			return v.error(node, "%v", err)
		}
		return Nature{Type: t}
	} else if len(f.Types) == 0 {
		nt, err := v.checkArguments(f.Name, Nature{Type: f.Type()}, arguments, node)
		if err != nil {
			if v.err == nil {
				v.err = err
			}
			return unknown
		}
		// No type was specified, so we assume the function returns any.
		return nt
	}
	var lastErr *file.Error
	for _, t := range f.Types {
		outNature, err := v.checkArguments(f.Name, Nature{Type: t}, arguments, node)
		if err != nil {
			lastErr = err
			continue
		}

		// As we found the correct function overload, we can stop the loop.
		// Also, we need to set the correct nature of the callee so compiler,
		// can correctly handle OpDeref opcode.
		if callNode, ok := node.(*ast.CallNode); ok {
			callNode.Callee.SetType(t)
		}

		return outNature
	}
	if lastErr != nil {
		if v.err == nil {
			v.err = lastErr
		}
		return unknown
	}

	return v.error(node, "no matching overload for %v", f.Name)
}

func (v *Checker) checkArguments(
	name string,
	fn Nature,
	arguments []ast.Node,
	node ast.Node,
) (Nature, *file.Error) {
	if fn.IsUnknown() {
		return unknown, nil
	}

	numOut := fn.NumOut()
	if numOut == 0 {
		return unknown, &file.Error{
			Location: node.Location(),
			Message:  fmt.Sprintf("func %v doesn't return value", name),
		}
	}
	if numOut > 2 {
		return unknown, &file.Error{
			Location: node.Location(),
			Message:  fmt.Sprintf("func %v returns more then two values", name),
		}
	}

	// If func is method on an env, first argument should be a receiver,
	// and actual arguments less than fnNumIn by one.
	fnNumIn := fn.NumIn()
	if fn.Method { // TODO: Move subtraction to the Nature.NumIn() and Nature.In() methods.
		fnNumIn--
	}
	// Skip first argument in case of the receiver.
	fnInOffset := 0
	if fn.Method {
		fnInOffset = 1
	}

	var err *file.Error
	isVariadic := fn.IsVariadic()
	if isVariadic {
		if len(arguments) < fnNumIn-1 {
			err = &file.Error{
				Location: node.Location(),
				Message:  fmt.Sprintf("not enough arguments to call %v", name),
			}
		}
	} else {
		if len(arguments) > fnNumIn {
			err = &file.Error{
				Location: node.Location(),
				Message:  fmt.Sprintf("too many arguments to call %v", name),
			}
		}
		if len(arguments) < fnNumIn {
			err = &file.Error{
				Location: node.Location(),
				Message:  fmt.Sprintf("not enough arguments to call %v", name),
			}
		}
	}

	if err != nil {
		// If we have an error, we should still visit all arguments to
		// type check them, as a patch can fix the error later.
		for _, arg := range arguments {
			_ = v.visit(arg)
		}
		return *fn.Out(0), err
	}

	for i, arg := range arguments {
		argNature := v.visit(arg)

		var in Nature
		if isVariadic && i >= fnNumIn-1 {
			// For variadic arguments fn(xs ...int), go replaces type of xs (int) with ([]int).
			// As we compare arguments one by one, we need underling type.
			in = fn.In(fnNumIn - 1).Elem()
		} else {
			in = *fn.In(i + fnInOffset)
		}

		if in.IsFloat() && argNature.IsInteger() {
			traverseAndReplaceIntegerNodesWithFloatNodes(&arguments[i], in)
			continue
		}

		inKind := in.Kind()
		if in.IsInteger() && argNature.IsInteger() && argNature.Kind() != inKind {
			traverseAndReplaceIntegerNodesWithIntegerNodes(&arguments[i], in)
			continue
		}

		if argNature.Nil {
			if inKind == reflect.Ptr || inKind == reflect.Interface {
				continue
			}
			return unknown, &file.Error{
				Location: arg.Location(),
				Message:  fmt.Sprintf("cannot use nil as argument (type %s) to call %v", in.String(), name),
			}
		}

		// Check if argument is assignable to the function input type.
		// We check original type (like *time.Time), not dereferenced type,
		// as function input type can be pointer to a struct.
		assignable := argNature.AssignableTo(in)

		// We also need to check if dereference arg type is assignable to the function input type.
		// For example, func(int) and argument *int. In this case we will add OpDeref to the argument,
		// so we can call the function with *int argument.
		if !assignable && argNature.IsPointer() {
			nt := argNature.Deref()
			assignable = nt.AssignableTo(in)
		}

		if !assignable && !argNature.IsUnknown() {
			return unknown, &file.Error{
				Location: arg.Location(),
				Message:  fmt.Sprintf("cannot use %s as argument (type %s) to call %v ", argNature.String(), in.String(), name),
			}
		}
	}

	return *fn.Out(0), nil
}

func traverseAndReplaceIntegerNodesWithFloatNodes(node *ast.Node, newNature Nature) {
	switch (*node).(type) {
	case *ast.IntegerNode:
		*node = &ast.FloatNode{Value: float64((*node).(*ast.IntegerNode).Value)}
		(*node).SetType(newNature.Type)
	case *ast.UnaryNode:
		unaryNode := (*node).(*ast.UnaryNode)
		traverseAndReplaceIntegerNodesWithFloatNodes(&unaryNode.Node, newNature)
	case *ast.BinaryNode:
		binaryNode := (*node).(*ast.BinaryNode)
		switch binaryNode.Operator {
		case "+", "-", "*":
			traverseAndReplaceIntegerNodesWithFloatNodes(&binaryNode.Left, newNature)
			traverseAndReplaceIntegerNodesWithFloatNodes(&binaryNode.Right, newNature)
		}
	}
}

func traverseAndReplaceIntegerNodesWithIntegerNodes(node *ast.Node, newNature Nature) {
	switch (*node).(type) {
	case *ast.IntegerNode:
		(*node).SetType(newNature.Type)
	case *ast.UnaryNode:
		(*node).SetType(newNature.Type)
		unaryNode := (*node).(*ast.UnaryNode)
		traverseAndReplaceIntegerNodesWithIntegerNodes(&unaryNode.Node, newNature)
	case *ast.BinaryNode:
		// TODO: Binary node return type is dependent on the type of the operands. We can't just change the type of the node.
		binaryNode := (*node).(*ast.BinaryNode)
		switch binaryNode.Operator {
		case "+", "-", "*":
			traverseAndReplaceIntegerNodesWithIntegerNodes(&binaryNode.Left, newNature)
			traverseAndReplaceIntegerNodesWithIntegerNodes(&binaryNode.Right, newNature)
		}
	}
}

func (v *Checker) predicateNode(node *ast.PredicateNode) Nature {
	nt := v.visit(node.Node)
	var out []reflect.Type
	if nt.IsUnknown() {
		out = append(out, anyType)
	} else if !nt.Nil {
		out = append(out, nt.Type)
	}
	return Nature{
		Type:         reflect.FuncOf([]reflect.Type{anyType}, out, false),
		PredicateOut: &nt,
	}
}

func (v *Checker) pointerNode(node *ast.PointerNode) Nature {
	if len(v.predicateScopes) == 0 {
		return v.error(node, "cannot use pointer accessor outside predicate")
	}
	scope := v.predicateScopes[len(v.predicateScopes)-1]
	if node.Name == "" {
		if scope.collection.IsUnknown() {
			return unknown
		}
		switch scope.collection.Kind() {
		case reflect.Array, reflect.Slice:
			return scope.collection.Elem()
		}
		return v.error(node, "cannot use %v as array", scope)
	}
	if scope.vars != nil {
		for i := range scope.vars {
			if node.Name == scope.vars[i].name {
				return scope.vars[i].nature
			}
		}
	}
	return v.error(node, "unknown pointer #%v", node.Name)
}

func (v *Checker) variableDeclaratorNode(node *ast.VariableDeclaratorNode) Nature {
	if _, ok := v.config.Env.Get(&v.config.NtCache, node.Name); ok {
		return v.error(node, "cannot redeclare %v", node.Name)
	}
	if _, ok := v.config.Functions[node.Name]; ok {
		return v.error(node, "cannot redeclare function %v", node.Name)
	}
	if _, ok := v.config.Builtins[node.Name]; ok {
		return v.error(node, "cannot redeclare builtin %v", node.Name)
	}
	if _, ok := v.lookupVariable(node.Name); ok {
		return v.error(node, "cannot redeclare variable %v", node.Name)
	}
	varNature := v.visit(node.Value)
	v.varScopes = append(v.varScopes, varScope{node.Name, varNature})
	exprNature := v.visit(node.Expr)
	v.varScopes = v.varScopes[:len(v.varScopes)-1]
	return exprNature
}

func (v *Checker) sequenceNode(node *ast.SequenceNode) Nature {
	if len(node.Nodes) == 0 {
		return v.error(node, "empty sequence expression")
	}
	var last Nature
	for _, node := range node.Nodes {
		last = v.visit(node)
	}
	return last
}

func (v *Checker) lookupVariable(name string) (varScope, bool) {
	for i := len(v.varScopes) - 1; i >= 0; i-- {
		if v.varScopes[i].name == name {
			return v.varScopes[i], true
		}
	}
	return varScope{}, false
}

func (v *Checker) conditionalNode(node *ast.ConditionalNode) Nature {
	c := v.visit(node.Cond)
	if !c.IsBool() && !c.IsUnknown() {
		return v.error(node.Cond, "non-bool expression (type %v) used as condition", c.String())
	}

	t1 := v.visit(node.Exp1)
	t2 := v.visit(node.Exp2)

	if t1.Nil && !t2.Nil {
		return t2
	}
	if !t1.Nil && t2.Nil {
		return t1
	}
	if t1.Nil && t2.Nil {
		return nilNature
	}
	if t1.AssignableTo(t2) {
		return t1
	}
	return unknown
}

func (v *Checker) arrayNode(node *ast.ArrayNode) Nature {
	var prev Nature
	allElementsAreSameType := true
	for i, node := range node.Nodes {
		curr := v.visit(node)
		if i > 0 {
			if curr.Kind() != prev.Kind() {
				allElementsAreSameType = false
			}
		}
		prev = curr
	}
	if allElementsAreSameType {
		return prev.MakeArrayOf()
	}
	return arrayNature
}

func (v *Checker) mapNode(node *ast.MapNode) Nature {
	for _, pair := range node.Pairs {
		v.visit(pair)
	}
	return mapNature
}

func (v *Checker) pairNode(node *ast.PairNode) Nature {
	v.visit(node.Key)
	v.visit(node.Value)
	return nilNature
}
