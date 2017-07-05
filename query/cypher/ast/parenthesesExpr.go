package ast

type Parentheses int

const (
	RPAREN Parentheses = iota // )
	LPAREN                    // (
)

type ParenthesesExpr struct {
	Parentheses
}

func (ParenthesesExpr) exprNode() {}

func ParenthesesPrecedence(item ParenthesesExpr) int {
	if item.Parentheses == LPAREN {
		return 11
	} else if item.Parentheses == RPAREN {
		return 12
	} else {
		return 20
	}
}