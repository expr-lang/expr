package ast

type BaseVisitor struct{}

func (BaseVisitor) NilNode(node *NilNode) {}

func (BaseVisitor) IdentifierNode(node *IdentifierNode) {}

func (BaseVisitor) IntegerNode(node *IntegerNode) {}

func (BaseVisitor) FloatNode(node *FloatNode) {}

func (BaseVisitor) BoolNode(node *BoolNode) {}

func (BaseVisitor) StringNode(node *StringNode) {}

func (BaseVisitor) UnaryNode(node *UnaryNode) {}

func (BaseVisitor) BinaryNode(node *BinaryNode) {}

func (BaseVisitor) MatchesNode(node *MatchesNode) {}

func (BaseVisitor) PropertyNode(node *PropertyNode) {}

func (BaseVisitor) IndexNode(node *IndexNode) {}

func (BaseVisitor) MethodNode(node *MethodNode) {}

func (BaseVisitor) FunctionNode(node *FunctionNode) {}

func (BaseVisitor) BuiltinNode(node *BuiltinNode) {}

func (BaseVisitor) ClosureNode(node *ClosureNode) {}

func (BaseVisitor) PointerNode(node *PointerNode) {}

func (BaseVisitor) ConditionalNode(node *ConditionalNode) {}

func (BaseVisitor) ArrayNode(node *ArrayNode) {}

func (BaseVisitor) MapNode(node *MapNode) {}

func (BaseVisitor) PairNode(node *PairNode) {}

func (BaseVisitor) Node(node *Node) {}
