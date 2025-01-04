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
	var s string

	for !f.currTokenIs(sqllexer.EOF, "") {
		if f.currTokenIs(sqllexer.IDENT, "SELECT") {
			s = f.formatSelectStatement()
		}

		f.nextToken()
	}

	return s
}

func (f *Formatter) formatSelectStatement() string {
	ss := SelectStatement{}
	f.nextToken()

	for !(isKeyword(f.currToken) || f.currTypeIs(sqllexer.EOF)) {
		sc := SelectedColumn{}

		exp, _ := f.parseExpression()
		sc.Exp = exp
		f.nextToken()

		if f.currTokenIs(sqllexer.OPERATOR, "::") {
			sc.Cast = f.parseTypecastExpression()
			f.nextToken()
		}

		if f.currTokenIs(sqllexer.IDENT, "AS") {
			if f.peekTypeIs(sqllexer.IDENT) && !isKeyword(f.peekToken) {
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

	return ss.Format(0)
}

// expects current token to NOT be a keyword
// expect that when it's done executing the nextToken is the start of the next section
func (f *Formatter) parseExpression() (Expression, bool) {
	var exp Expression

	if f.currTypeIs(sqllexer.IDENT) {
		exp = f.parseIdentExpression()
	} else if f.currTypeIs(sqllexer.NUMBER) {
		exp = f.parseNumericExpression()
	} else if f.currTokenIs(sqllexer.PUNCTUATION, ".") {
		if f.peekTypeIs(sqllexer.IDENT) {
			exp = f.parseIdentExpression()
		} else if f.peekTypeIs(sqllexer.NUMBER) {
			exp = f.parseNumericExpression()
		}
	} else if f.currTokenIs(sqllexer.PUNCTUATION, "(") {
		exp = f.parseArgsExpression()
	} else if f.currTokenIs(sqllexer.OPERATOR, "::") {
		exp = f.parseTypecastExpression()
	} else {
		return nil, false
	}

	return exp, true
}

func (f *Formatter) parseIdentExpression() Expression { // this is bad practice
	// [edgecase] - in the case a dumbass does ".asdf"
	if f.currTokenIs(sqllexer.PUNCTUATION, ".") {
		f.nextToken()
		return QualifiedIdentExpression{Left: sqllexer.Token{Type: sqllexer.ERROR}, Right: f.currToken}
	}

	if f.peekTokenIs(sqllexer.PUNCTUATION, ".") {
		qie := QualifiedIdentExpression{Left: f.currToken}
		f.nextToken()

		if f.peekTypeIs(sqllexer.IDENT) && !isKeyword(f.currToken) {
			f.nextToken()
			qie.Right = f.currToken
		} else {
			qie.Right = sqllexer.Token{Type: sqllexer.ERROR}
		}

		return qie
	}

	return IdentExpression{Token: f.currToken}
}

func (f *Formatter) parseNumericExpression() Expression { // This is bad practice
	// [edgecase] - in the case a dumbass does ".123"
	if f.currTokenIs(sqllexer.PUNCTUATION, ".") {
		f.nextToken()
		return DecimalExpression{Right: f.currToken}
	}

	if f.peekTokenIs(sqllexer.PUNCTUATION, ".") {
		de := DecimalExpression{Left: f.currToken}
		f.nextToken()

		if f.peekTypeIs(sqllexer.NUMBER) {
			f.nextToken()
			de.Right = f.currToken
		} else {
			de.Right = sqllexer.Token{}
		}

		return de
	}

	return IntExpression{Token: f.currToken}
}

func (f *Formatter) parseArgsExpression() ArgsExpression {
	args := ArgsExpression{Exps: []Expression{}}
	f.nextToken()

	for !f.currTokenIs(sqllexer.PUNCTUATION, ")") {
		e, ok := f.parseExpression()
		if !ok {
			break
		}
		f.nextToken()

		if f.currTokenIs(sqllexer.PUNCTUATION, ",") {
			f.nextToken()
		}

		args.Exps = append(args.Exps, e)
	}

	return args
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
	exp := DatatypeExpression{Datatype: f.currToken, Args: ArgsExpression{}}

	if f.peekTokenIs(sqllexer.PUNCTUATION, "(") {
		f.nextToken()
		ge := f.parseArgsExpression()
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
