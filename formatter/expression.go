package formatter

import (
	"strings"

	"github.com/DataDog/go-sqllexer"
)

type Expression interface {
	expressionNode()
	String() string
	Type() string
}

type NumericExpression struct {
	Token sqllexer.Token
}

func (ne NumericExpression) expressionNode() {}
func (ne NumericExpression) Type() string    { return "numeric" }
func (ne NumericExpression) String() string  { return ne.Token.Value }

type StringExpression struct {
	Token sqllexer.Token
}

func (se StringExpression) expressionNode() {}
func (se StringExpression) Type() string    { return "string" }
func (se StringExpression) String() string  { return se.Token.Value }

type BooleanExpression struct {
	Token sqllexer.Token
}

func (be BooleanExpression) expressionNode() {}
func (be BooleanExpression) Type() string    { return "boolean" }
func (be BooleanExpression) String() string  { return strings.ToUpper(be.Token.Value) }

type IdentExpression struct {
	Token sqllexer.Token
}

func (ie IdentExpression) expressionNode() {}
func (ie IdentExpression) Type() string    { return "identifier" }
func (ie IdentExpression) String() string  { return strings.ToLower(ie.Token.Value) }

type OperationKeywordExpression struct {
	Token sqllexer.Token
}

func (oe OperationKeywordExpression) expressionNode() {}
func (oe OperationKeywordExpression) Type() string    { return "operation-keyword" }
func (oe OperationKeywordExpression) String() string  { return strings.ToUpper(oe.Token.Value) }

type QuotedIdentExpression struct {
	Token sqllexer.Token
}

func (qie QuotedIdentExpression) expressionNode() {}
func (qie QuotedIdentExpression) Type() string    { return "quoted-identifier" }
func (qie QuotedIdentExpression) String() string  { return qie.Token.Value }

type OperatorExpression struct {
	Token sqllexer.Token
}

func (oe OperatorExpression) expressionNode() {}
func (oe OperatorExpression) Type() string    { return "operator" }
func (oe OperatorExpression) String() string {
	switch oe.Token.Value {
	case ",":
		return oe.Token.Value + " "
	case ".", "(", ")":
		return oe.Token.Value
	default:
		return " " + oe.Token.Value + " "
	}
}

type TypecastExpression struct {
	Datatype DatatypeExpression
}

func (te TypecastExpression) expressionNode() {}
func (te TypecastExpression) Type() string    { return "typecast" }
func (te TypecastExpression) String() string {
	str := "::"
	if te.Datatype.Datatype.Type != sqllexer.ERROR {
		str += " " + te.Datatype.String()
	}

	return str
}

type DatatypeExpression struct {
	Datatype sqllexer.Token
	Args     GroupedExpression
	hasArgs  bool
}

func (de DatatypeExpression) expressionNode() {}
func (de DatatypeExpression) Type() string    { return "datatype" }
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

type CallExpression struct {
	Function sqllexer.Token
	Args     GroupedExpression
}

func (ce CallExpression) expressionNode() {}
func (ce CallExpression) Type() string    { return "call" }
func (ce CallExpression) String() string {
	str := ""
	if len(strings.Split(ce.Function.Value, ".")) >= 2 {
		str += strings.ToLower(ce.Function.Value) // custom function
	} else {
		str += strings.ToUpper(ce.Function.Value) // inbuilt function
	}

	str += ce.Args.String()

	return str
}

type GroupedExpression struct {
	Exps     []Expression
	Comment  CommentExpression
	HasParen bool
}

func (ge GroupedExpression) expressionNode() {}
func (ge GroupedExpression) Type() string    { return "group" }
func (ge GroupedExpression) String() string {
	str := ""
	if ge.HasParen {
		str += "("
	}

	prevType := "operator"
	for i := 0; i < len(ge.Exps); i++ {
		if !(prevType == "operator") && ge.Exps[i].Type() != "operator" { // this is so fucking bad holy shit
			str += " "
		}

		str += ge.Exps[i].String()
		prevType = ge.Exps[i].Type()
	}

	if ge.HasParen {
		str += ")"
	}
	return str
}

type WindowExpression struct {
	Call CallExpression
	Args GroupedExpression
}

func (we WindowExpression) expressionNode() {}
func (we WindowExpression) Type() string    { return "window-function" }
func (we WindowExpression) String() string {
	str := we.Call.String()
	str += " OVER "
	str += we.Args.String()

	return str
}

type AsExpression struct {
	Token sqllexer.Token // not an ident expression because it can't handle qualified identifiers
}

func (ae AsExpression) expressionNode() {}
func (ae AsExpression) Type() string    { return "as" }
func (ae AsExpression) String() string {
	str := "AS"

	if ae.Token.Type != sqllexer.ERROR {
		if ae.Token.Type == sqllexer.IDENT {
			str += " " + strings.ToLower(ae.Token.Value)
		} else if ae.Token.Type == sqllexer.QUOTED_IDENT {
			str += " " + ae.Token.Value
		}
	}

	return str
}

type CommentExpression struct {
	Token sqllexer.Token
}

func (ce CommentExpression) expressionNode() {}
func (ce CommentExpression) Type() string    { return "comment" }
func (ce CommentExpression) String() string  { return ce.Token.Value }
