package formatter

import (
	"strings"

	"github.com/DataDog/go-sqllexer"
)

type Expression interface {
	expressionNode()
	String() string
}

type IntExpression struct {
	Token sqllexer.Token
}

func (ie IntExpression) expressionNode() {}
func (ie IntExpression) String() string  { return ie.Token.Value }

type DecimalExpression struct {
	Left  sqllexer.Token
	Right sqllexer.Token
}

func (de DecimalExpression) expressionNode() {}
func (de DecimalExpression) String() string {
	right := ""
	if de.Right.Type != sqllexer.ERROR {
		right = de.Right.Value
	}

	left := ""
	if de.Left.Type != sqllexer.ERROR {
		left = de.Left.Value
	}

	return left + "." + right
}

type IdentExpression struct {
	Token sqllexer.Token
}

func (ie IdentExpression) expressionNode() {}
func (ie IdentExpression) String() string  { return strings.ToLower(ie.Token.Value) }

type QualifiedIdentExpression struct {
	Left  sqllexer.Token
	Right sqllexer.Token
}

func (qie QualifiedIdentExpression) expressionNode() {}
func (qie QualifiedIdentExpression) String() string {
	right := ""
	if qie.Right.Type != sqllexer.ERROR {
		right = qie.Right.Value
	}

	left := ""
	if qie.Left.Type != sqllexer.ERROR {
		left = qie.Left.Value
	}

	return left + "." + right
}

type TypecastExpression struct {
	Datatype DatatypeExpression
}

func (te TypecastExpression) expressionNode() {}
func (te TypecastExpression) String() string {
	return ":: " + te.Datatype.String()
}

type DatatypeExpression struct {
	Datatype sqllexer.Token
	Args     GroupExpression
	hasArgs  bool
}

func (de DatatypeExpression) expressionNode() {}
func (de DatatypeExpression) String() string {
	str := ""

	if de.Datatype.Type != sqllexer.ERROR {
		str += strings.ToUpper(de.Datatype.Value)
	}

	if de.hasArgs {
		str += de.Args.String()
	}

	return str
}

// should be renamed to args expression
type GroupExpression struct {
	Exps []Expression
}

func (ge GroupExpression) expressionNode() {}
func (ge GroupExpression) String() string {
	str := "("

	for i := 0; i < len(ge.Exps); i++ {
		str += ge.Exps[i].String()

		if i != len(ge.Exps)-1 {
			str += ", "
		}
	}

	str += ")"
	return str
}
