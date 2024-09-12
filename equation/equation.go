package equation

import (
	"math"
	"math/rand"
	"strconv"
)

type BaseNode interface {
	Eval(x, y float32) float32
	String() string
	GetParent() BaseNode
	SetParent(parent BaseNode)
	GetChildren() []BaseNode
	SetChildren(children []BaseNode)
	AddRandomNode(node BaseNode)
	AddLeafNode(leaf BaseNode) bool // true if successfully added
	NodeCount() int
}

// ANCHOR
type Node struct {
	Parent   BaseNode
	Children []BaseNode
}

func (node *Node) Eval(x, y float32) float32 {
	panic("call eval on basenode")
}

func (node *Node) String() string {
	panic("call string on basenode")
}

func (node *Node) GetParent() BaseNode {
	return node.Parent
}

func (node *Node) SetParent(parent BaseNode) {
	node.Parent = parent
}

func (node *Node) GetChildren() []BaseNode {
	return node.Children
}

func (node *Node) SetChildren(children []BaseNode) {
	node.Children = children
}

func (node *Node) AddRandomNode(newNode BaseNode) {
	index := rand.Intn(len(node.Children))
	if node.Children[index] == nil {
		node.Children[index] = newNode
		newNode.SetParent(node)
	} else {
		node.Children[index].AddRandomNode(newNode)
	}
}

func (node *Node) AddLeafNode(leaf BaseNode) bool {
	for i, child := range node.Children {
		if child == nil {
			node.Children[i] = leaf
			leaf.SetParent(node)
			return true
		} else if node.Children[i].AddLeafNode(leaf) {
			return true
		}
	}
	return false
}

func (node *Node) NodeCount() int {
	count := 1
	for _, child := range node.Children {
		count += child.NodeCount()
	}
	return count
}

func NewNode(size int) Node {
	return Node{nil, make([]BaseNode, size)}
}

func GetNthNode(tree BaseNode, n int) BaseNode {
	count := 0
	var result BaseNode

	var dfs func(node BaseNode)
	dfs = func(node BaseNode) {
		count++
		if count == n {
			result = node
			return
		}

		for _, child := range node.GetChildren() {
			dfs(child)
			if result != nil {
				return
			}
		}
	}

	dfs(tree)

	if result == nil {
		panic("node does not exist")
	}

	return result
}

func Mutate(node BaseNode) BaseNode {
	var newNode BaseNode

	opTypeCount := 8
	leafTypeCount := 3
	n := rand.Intn(opTypeCount + leafTypeCount)
	if n < 8 {
		newNode = RandomOpNode()
	} else {
		newNode = RandomLeafNode()
	}

	// Point ParentNode to NewNode
	if parentNode := node.GetParent(); parentNode != nil {
		for i, child := range parentNode.GetChildren() {
			if child == node {
				parentNode.GetChildren()[i] = newNode
			}
		}
		newNode.SetParent(parentNode)
	}

	// Add children from OldNode to NewNode
	for i, child := range node.GetChildren() {
		if i >= len(newNode.GetChildren()) {
			break
		}
		newNode.GetChildren()[i] = child
		child.SetParent(newNode)
	}

	// Add leaf to children if they are empty
	for i, child := range newNode.GetChildren() {
		if child == nil {
			leaf := RandomLeafNode()
			newNode.GetChildren()[i] = leaf
			leaf.SetParent(newNode)
		}
	}

	return newNode
}

// ANCHOR
type OpX struct {
	Node
}

func NewOpX() *OpX {
	return &OpX{NewNode(0)}
}

func (op *OpX) Eval(x, y float32) float32 {
	return x
}

func (op *OpX) String() string {
	return "X"
}

// ANCHOR
type OpY struct {
	Node
}

func NewOpY() *OpY {
	return &OpY{NewNode(0)}
}

func (op *OpY) Eval(x, y float32) float32 {
	return y
}

func (op *OpY) String() string {
	return "Y"
}

// ANCHOR
type OpConstant struct {
	Node
	value float32
}

func NewOpConstant() *OpConstant {
	return &OpConstant{NewNode(0), rand.Float32()*2 - 1}
}

func (op *OpConstant) Eval(x, y float32) float32 {
	return op.value
}

func (op *OpConstant) String() string {
	return strconv.FormatFloat(float64(op.value), 'f', 9, 32)
}

// ANCHOR
type OpPlus struct {
	Node
}

func NewOpPlus() *OpPlus {
	return &OpPlus{NewNode(2)}
}

func (op *OpPlus) Eval(x, y float32) float32 {
	return op.Children[0].Eval(x, y) + op.Children[1].Eval(x, y)
}

func (op *OpPlus) String() string {
	return "Plus(" + op.Children[0].String() + ", " + op.Children[1].String() + ")"
}

// ANCHOR
type OpMinus struct {
	Node
}

func NewOpMinus() *OpMinus {
	return &OpMinus{NewNode(2)}
}

func (op *OpMinus) Eval(x, y float32) float32 {
	return op.Children[0].Eval(x, y) - op.Children[1].Eval(x, y)
}

func (op *OpMinus) String() string {
	return "Minus(" + op.Children[0].String() + ", " + op.Children[1].String() + ")"
}

// ANCHOR
type OpMult struct {
	Node
}

func NewOpMult() *OpMult {
	return &OpMult{NewNode(2)}
}

func (op *OpMult) Eval(x, y float32) float32 {
	return op.Children[0].Eval(x, y) * op.Children[1].Eval(x, y)
}

func (op *OpMult) String() string {
	return "Mult(" + op.Children[0].String() + ", " + op.Children[1].String() + ")"
}

// ANCHOR
type OpDiv struct {
	Node
}

func NewOpDiv() *OpDiv {
	return &OpDiv{NewNode(2)}
}

func (op *OpDiv) Eval(x, y float32) float32 {
	return op.Children[0].Eval(x, y) / op.Children[1].Eval(x, y)
}

func (op *OpDiv) String() string {
	return "Div(" + op.Children[0].String() + ", " + op.Children[1].String() + ")"
}

// ANCHOR
type OpSin struct {
	Node
}

func NewOpSin() *OpSin {
	return &OpSin{NewNode(1)}
}

func (op *OpSin) Eval(x, y float32) float32 {
	return float32(math.Sin(float64(op.Children[0].Eval(x, y))))
}

func (op *OpSin) String() string {
	return "Sin(" + op.Children[0].String() + ")"
}

// ANCHOR
type OpCos struct {
	Node
}

func NewOpCos() *OpCos {
	return &OpCos{NewNode(1)}
}

func (op *OpCos) Eval(x, y float32) float32 {
	return float32(math.Cos(float64(op.Children[0].Eval(x, y))))
}

func (op *OpCos) String() string {
	return "Cos(" + op.Children[0].String() + ")"
}

// ANCHOR
type OpAtan struct {
	Node
}

func NewOpAtan() *OpAtan {
	return &OpAtan{NewNode(1)}
}

func (op *OpAtan) Eval(x, y float32) float32 {
	return float32(math.Atan(float64(op.Children[0].Eval(x, y))))
}

func (op *OpAtan) String() string {
	return "Atan(" + op.Children[0].String() + ")"
}

// ANCHOR
type OpAtan2 struct {
	Node
}

func NewOpAtan2() *OpAtan2 {
	return &OpAtan2{NewNode(2)}
}

func (op *OpAtan2) Eval(x, y float32) float32 {
	return float32(math.Atan2(float64(op.Children[0].Eval(x, y)), float64(op.Children[1].Eval(x, y))))
}

func (op *OpAtan2) String() string {
	return "Atan2(" + op.Children[0].String() + ", " + op.Children[1].String() + ")"
}

func RandomOpNode() BaseNode {
	n := rand.Intn(8)
	switch n {
	case 0:
		return NewOpPlus()
	case 1:
		return NewOpMinus()
	case 2:
		return NewOpMult()
	case 3:
		return NewOpDiv()
	case 4:
		return NewOpSin()
	case 5:
		return NewOpCos()
	case 6:
		return NewOpAtan()
	case 7:
		return NewOpAtan2()
	}
	panic("get random node failed")
}

func RandomLeafNode() BaseNode {
	n := rand.Intn(3)
	switch n {
	case 0:
		return NewOpX()
	case 1:
		return NewOpY()
	case 2:
		return NewOpConstant()
	}
	panic("get random node failed")
}
