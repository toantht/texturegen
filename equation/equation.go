package equation

import (
	"math"
	"math/rand"
	"strconv"
)

type BaseNode interface {
	Eval(x, y float32) float32
	String() string
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

func NewNode(size int) Node {
	return Node{nil, make([]BaseNode, size)}
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
	return &OpConstant{NewNode(0), rand.Float32()}
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
