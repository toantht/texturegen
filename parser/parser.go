package parser

import (
	"strconv"

	. "github.com/toantht/texturegen/equation"
)

func Parse(tokens []Token) BaseNode {
	index := 0

	var buildTree func(parent BaseNode) BaseNode
	buildTree = func(parent BaseNode) BaseNode {
		token := tokens[index]
		index++

		switch token.typ {
		case EOF:
			return nil
		case OPEN_PAREN:
			return buildTree(parent)
		case CLOSE_PAREN:
			return buildTree(parent)
		case CONSTANT:
			value, err := strconv.ParseFloat(token.value, 32)
			if err != nil {
				panic(err)
			}
			node := NewOpConstant(float32(value))
			node.SetParent(parent)
			return node
		case OPERATION:
			node := tokenToNode(token)
			node.SetParent(parent)
			for i := range node.GetChildren() {
				node.GetChildren()[i] = buildTree(node)
			}
			return node
		}
		return nil
	}

	return buildTree(nil)

}

func tokenToNode(token Token) BaseNode {
	switch token.value {
	case "X":
		return NewOpX()
	case "Y":
		return NewOpY()
	case "Plus":
		return NewOpPlus()
	case "Minus":
		return NewOpMinus()
	case "Mult":
		return NewOpMult()
	case "Div":
		return NewOpDiv()
	case "Sin":
		return NewOpSin()
	case "Cos":
		return NewOpCos()
	case "Atan":
		return NewOpAtan()
	case "Atan2":
		return NewOpAtan2()
	case "EquationImage":
		return NewOpImage()
	}
	panic("unknown token")
}
