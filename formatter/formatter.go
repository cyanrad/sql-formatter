package formatter

import (
	"github.com/DataDog/go-sqllexer"
)

type Formatter struct {
	tokens []sqllexer.Token
	index  int

	currToken sqllexer.Token
	peekToken sqllexer.Token
}

func Create(sql string) Formatter {
	lexer := sqllexer.New(sql)
	tokens := lexer.ScanAll()
	f := Formatter{tokens: clearWhiteSpace(tokens), index: -2}

	// loading tokens into curr and peek
	f.nextToken()
	f.nextToken()

	return f
}

func (f *Formatter) Format() string {
	formatted := ""

	for !f.currTokenIs(sqllexer.EOF, "") {
		s := ""
		if f.currTokenIs(sqllexer.IDENT, "SELECT") {
			s = f.formatSelectStatement() + "\n;\n"
		}

		formatted += s
		if s == "" {
			f.nextToken()
		}
	}

	return formatted[:len(formatted)-1]
}

// when it's done currToken will be the first token of the next statement
func (f *Formatter) formatSelectStatement() string {
	ss := SelectStatement{}
	f.nextToken()

	for !(isStatementKeyword(f.currToken) || f.currTypeIs(sqllexer.EOF) || f.currTokenIs(sqllexer.PUNCTUATION, ";")) {
		// [todo] - this should be broken down into a fromatSelectColumnStatement
		sc := SelectedColumn{}

		exp, _ := f.parseExpression()
		sc.Exp = exp
		f.nextToken()

		if f.currTokenIs(sqllexer.OPERATOR, "::") {
			sc.Cast = f.parseTypecastExpression()
			f.nextToken()
		}

		if f.currTokenIs(sqllexer.IDENT, "AS") {
			if (f.peekTypeIs(sqllexer.IDENT) || f.peekTypeIs(sqllexer.QUOTED_IDENT)) && !isStatementKeyword(f.peekToken) {
				f.nextToken()
				sc.Alias = f.currToken
			}
			f.nextToken()
		}

		if f.currTokenIs(sqllexer.PUNCTUATION, ",") {
			f.nextToken()
		}

		ss.Columns = append(ss.Columns, sc)
	}

	if f.currTokenIs(sqllexer.PUNCTUATION, ";") {
		f.nextToken()
	}

	return ss.Format(0)
}

// expects current token to NOT be a keyword
// expect that when it's done executing the nextToken is the start of the next section
func (f *Formatter) parseExpression() (Expression, bool) {
	var exp Expression

	if f.currTypeIs(sqllexer.IDENT) {
		exp = f.parseIdentifier()
	} else if f.currTypeIs(sqllexer.QUOTED_IDENT) {
		exp = f.parseQuotedIdentExpression()
	} else if f.currTypeIs(sqllexer.NUMBER) {
		exp = f.parseNumericExpression()
	} else if f.currTypeIs(sqllexer.STRING) {
		exp = f.parseStringExpression()
	} else if f.currTokenIs(sqllexer.PUNCTUATION, ".") { // [edgecase] - stupid piece of shit istg
		if f.peekTypeIs(sqllexer.IDENT) {
			exp = f.parseIdentifier()
		} else if f.peekTypeIs(sqllexer.QUOTED_IDENT) {
			exp = f.parseIdentExpression()
		} else if f.peekTypeIs(sqllexer.NUMBER) {
			exp = f.parseNumericExpression()
		}
	} else if f.currTokenIs(sqllexer.PUNCTUATION, "(") || f.peekTypeIs(sqllexer.OPERATOR) {
		exp = f.parseGroupedExpression()
	} else if f.currTypeIs(sqllexer.FUNCTION) { // [todo] - this should be replaced with the OperationExpression in the future
		callExp := f.parseCallExpression()

		if f.peekTokenIs(sqllexer.IDENT, "OVER") {
			f.nextToken()
			exp = f.parseWindowExpression(callExp)
		} else {
			exp = callExp
		}
	} else if f.currTokenIs(sqllexer.OPERATOR, "::") {
		exp = f.parseTypecastExpression()
	} else if f.currTypeIs(sqllexer.PUNCTUATION) || f.currTypeIs(sqllexer.OPERATOR) {
		exp = f.parseOperatorExpression()
	} else {
		return nil, false
	}

	return exp, true
}

func (f *Formatter) parseIdentifier() Expression {
	// [edgecase] - in the case a dumbass does ".asdf"
	if f.currTokenIs(sqllexer.PUNCTUATION, ".") {
		f.nextToken()
		f.currToken.Value = "." + f.currToken.Value
	}

	if isBoolean(f.currToken) {
		return f.parseBooleanExpression()
	} else if isOperationKeyword(f.currToken) {
		return f.parseOperationKeywordExpression()
	} else {
		return f.parseIdentExpression()
	}
}

func (f *Formatter) parseBooleanExpression() BooleanExpression {
	return BooleanExpression{Token: f.currToken}
}

func (f *Formatter) parseIdentExpression() IdentExpression {
	return IdentExpression{Token: f.currToken}
}

func (f *Formatter) parseOperationKeywordExpression() OperationKeywordExpression {
	return OperationKeywordExpression{Token: f.currToken}
}

func (f *Formatter) parseQuotedIdentExpression() Expression { // this is should actually be considered a legal crime
	// [edgecase] - in the case a dumbass does "."asdf""
	if f.currTokenIs(sqllexer.PUNCTUATION, ".") {
		f.nextToken()
		f.currToken.Value = "." + f.currToken.Value
	}

	// [edgecase] - because sqllexer doesn't treat quoted strings the same as other things
	if f.peekTokenIs(sqllexer.PUNCTUATION, ".") {
		return f.parseCallExpression()
	} else {
		return QuotedIdentExpression{Token: f.currToken}
	}
}

func (f *Formatter) parseNumericExpression() NumericExpression {
	// [edgecase] - in the case a dumbass does ".123"
	if f.currTokenIs(sqllexer.PUNCTUATION, ".") {
		f.nextToken()
		f.currToken.Value = "." + f.currToken.Value
	}

	return NumericExpression{Token: f.currToken}
}

func (f *Formatter) parseStringExpression() StringExpression {
	return StringExpression{Token: f.currToken}
}

func (f *Formatter) parseOperatorExpression() OperatorExpression {
	return OperatorExpression{Token: f.currToken}
}

func (f *Formatter) parseCallExpression() CallExpression {
	call := CallExpression{Function: f.currToken}
	if f.currTypeIs(sqllexer.QUOTED_IDENT) {
		f.nextToken() // skipping '.'
		f.nextToken() // at the function call name
		call.Function.Value += "." + f.currToken.Value
	}
	f.nextToken()

	args := f.parseGroupedExpression()
	call.Args = args

	return call
}

func (f *Formatter) parseWindowExpression(call CallExpression) WindowExpression {
	window := WindowExpression{Call: call}
	f.nextToken()
	window.Args = f.parseGroupedExpression()

	return window
}

func (f *Formatter) parseGroupedExpression() GroupedExpression {
	group := GroupedExpression{Exps: []Expression{}}
	if f.currTokenIs(sqllexer.PUNCTUATION, "(") {
		group.HasParen = true
	}
	f.nextToken()

	for !f.currTokenIs(sqllexer.PUNCTUATION, ")") {
		e, ok := f.parseExpression()
		if !ok {
			break
		}
		f.nextToken()

		group.Exps = append(group.Exps, e)
	}

	return group
}

func (f *Formatter) parseTypecastExpression() TypecastExpression {
	exp := TypecastExpression{}

	if f.peekTypeIs(sqllexer.IDENT) || f.peekTypeIs(sqllexer.FUNCTION) {
		f.nextToken()
		exp.Datatype = f.parseDatatypeExpression()
	}

	return exp
}

func (f *Formatter) parseDatatypeExpression() DatatypeExpression {
	exp := DatatypeExpression{Datatype: f.currToken, Args: GroupedExpression{}}

	if f.peekTokenIs(sqllexer.PUNCTUATION, "(") {
		f.nextToken()
		ge := f.parseGroupedExpression()
		exp.Args = ge
		exp.hasArgs = true
	}

	return exp
}

func (f *Formatter) nextToken() {
	f.currToken = f.peekToken

	f.index++
	if f.index+1 < len(f.tokens) {
		f.peekToken = f.tokens[f.index+1]
	} else {
		f.peekToken = sqllexer.Token{Type: sqllexer.EOF, Value: ""}
	}
}
