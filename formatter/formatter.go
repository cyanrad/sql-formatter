package formatter

import (
	"strings"

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

	for !(isStatementKeyword(f.currToken) ||
		f.currTypeIs(sqllexer.EOF) ||
		f.currTokenIs(sqllexer.PUNCTUATION, ";")) {

		sc := f.formatSelectedColumnStatement()
		ss.Columns = append(ss.Columns, sc)
	}

	if f.currTokenIs(sqllexer.PUNCTUATION, ";") {
		f.nextToken()
	}

	return ss.Format(0)
}

func (f *Formatter) formatSelectedColumnStatement() SelectedColumn {
	sc := SelectedColumn{Exps: []Expression{}}

	var exp Expression
	var prefixExpOperator OperatorExpression
	for i := 0; !(f.currTokenIs(sqllexer.PUNCTUATION, ",") ||
		f.currTokenIs(sqllexer.IDENT, "AS") ||
		f.currTypeIs(sqllexer.EOF) ||
		f.currTokenIs(sqllexer.PUNCTUATION, ";")); i++ { // [todo] I should prolly turn this into a map
		exp = f.parseExpression()
		if exp == nil {
			break
		}

		if exp.Type() == "comment" {
			sc.Comment = exp.(CommentExpression)
			if i == 0 {
				break
			}
			f.nextToken()
			continue
		} else if isPrefixOperator(strings.Trim(exp.String(), " ")) {
			if f.isPotentialPrefixExpression(sc) {
				prefixExpOperator = exp.(OperatorExpression)
				f.nextToken()
				continue
			}
		} else if exp.Type() == "group" {
			group := exp.(GroupedExpression)
			if group.Comment.String() != "" {
				sc.Comment = group.Comment
			}
		}

		// this means that the previous expression was a valid prefix operator
		// now we need to check if the current expression can take the prefix
		if prefixExpOperator.String() != "" {
			if isValidPrefixOperatorExpression(exp) {
				exp = PrefixOperatorExpression{Operator: prefixExpOperator.Token, Exp: exp}
				prefixExpOperator = OperatorExpression{} // resetting
			} else {
				// in this case we treat current and next exp as seperate expressions
				sc.Exps = append(sc.Exps, prefixExpOperator)
				sc.Exps = append(sc.Exps, exp)
				prefixExpOperator = OperatorExpression{} // resetting
				f.nextToken()
				continue
			}
		}

		// handling the final typecast
		if exp.Type() == "typecast" && f.peekTokenIs(sqllexer.IDENT, "AS") {
			sc.Cast = exp.(TypecastExpression)
			sc.hasCast = true
			f.nextToken()
			break
		}

		sc.Exps = append(sc.Exps, exp)
		f.nextToken()
	}

	if prefixExpOperator.String() != "" {
		// in this case we treat current and next exp as seperate expressions
		sc.Exps = append(sc.Exps, prefixExpOperator)
		prefixExpOperator = OperatorExpression{} // resetting
		f.nextToken()
	}

	sc.HandlePotentialComment(f)
	if f.currTokenIs(sqllexer.IDENT, "AS") {
		sc.hasAlias = true
		f.nextToken()
		sc.HandlePotentialComment(f)

		if (f.currTypeIs(sqllexer.IDENT) && !isStatementKeyword(f.currToken)) || f.currTypeIs(sqllexer.QUOTED_IDENT) {
			sc.Alias = AsExpression{Token: f.currToken}
			f.nextToken()
		}
	}

	sc.HandlePotentialComment(f)
	if f.currTokenIs(sqllexer.PUNCTUATION, ",") {
		f.nextToken()
	}

	sc.HandlePotentialComment(f)
	return sc
}

// expects current token to NOT be a keyword
// expect that when it's done executing the nextToken is the start of the next section
func (f *Formatter) parseExpression() Expression {
	var exp Expression

	if f.currTypeIs(sqllexer.IDENT) {
		exp = f.parseIdentifier()
	} else if f.currTypeIs(sqllexer.QUOTED_IDENT) {
		exp = f.parseQuotedIdentExpression()
	} else if f.currTypeIs(sqllexer.NUMBER) {
		exp = f.parseNumericExpression()
	} else if f.currTypeIs(sqllexer.STRING) {
		exp = f.parseStringExpression()
	} else if f.currTokenIs(sqllexer.PUNCTUATION, "(") {
		exp = f.parseGroupedExpression()
	} else if f.currTypeIs(sqllexer.FUNCTION) { // [todo] - this should be replaced with the OperationExpression in the future
		callExp := f.parseCallExpression()
		exp = callExp

		if f.peekTokenIs(sqllexer.IDENT, "OVER") {
			f.nextToken()
			exp = f.parseWindowExpression(callExp)
		}
	} else if f.currTokenIs(sqllexer.OPERATOR, "::") {
		exp = f.parseTypecastExpression()
	} else if f.currTypeIs(sqllexer.PUNCTUATION) || f.currTypeIs(sqllexer.OPERATOR) || f.currTypeIs(sqllexer.WILDCARD) {
		exp = f.parseOperatorExpression()
	} else if f.currTypeIs(sqllexer.COMMENT) {
		exp = f.parseCommentExpression()
	}

	return exp
}

func (f *Formatter) parseIdentifier() Expression {
	if isBoolean(f.currToken) {
		return f.parseBooleanExpression()
	} else if isOperationKeyword(f.currToken) {
		return f.parseOperationKeywordExpression()
	} else if !isStatementKeyword(f.currToken) {
		return f.parseIdentExpression()
	}

	return nil
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

	var prefixExpOperator OperatorExpression
	for !(f.currTokenIs(sqllexer.PUNCTUATION, ")") ||
		f.currTokenIs(sqllexer.IDENT, "AS") ||
		f.currTypeIs(sqllexer.EOF) ||
		f.currTokenIs(sqllexer.PUNCTUATION, ";")) {
		e := f.parseExpression()

		if e == nil {
			break
		} else if e.Type() == "comment" {
			group.Comment = e.(CommentExpression)
			f.nextToken()
			continue
		} else if isPrefixOperator(strings.Trim(e.String(), " ")) {
			if f.isPotentialPrefixExpression(group) {
				prefixExpOperator = e.(OperatorExpression)
				f.nextToken()
				continue
			}
		}

		// this means that the previous expression was a valid prefix operator
		// now we need to check if the current expression can take the prefix
		if prefixExpOperator.String() != "" {
			if isValidPrefixOperatorExpression(e) {
				e = PrefixOperatorExpression{Operator: prefixExpOperator.Token, Exp: e}
				prefixExpOperator = OperatorExpression{} // resetting
			} else {
				// in this case we treat current and next exp as seperate expressions
				group.Exps = append(group.Exps, prefixExpOperator)
				group.Exps = append(group.Exps, e)
				prefixExpOperator = OperatorExpression{} // resetting
				f.nextToken()
				continue
			}
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

func (f *Formatter) parseCommentExpression() CommentExpression {
	return CommentExpression{Token: f.currToken}
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

// I'm going to hell for this code
// Assumes it's being called when current token is a valid prefix operator
// this is happening because I fucking suck at programming
// returns either an operator expression or prefix expression
func (f *Formatter) isPotentialPrefixExpression(me MultiExpression) bool {
	Exps := me.Expressions()
	if len(Exps) == 0 {
		return true
	}

	prevExp := Exps[len(Exps)-1]
	leftParen := sqllexer.Token{Type: sqllexer.PUNCTUATION, Value: "("}
	if !isValidPrefixOperatorExpression(prevExp) ||
		f.index == 0 ||
		f.tokens[f.index-1] == leftParen {
		return true
	}
	return false
}

func isValidPrefixOperatorExpression(e Expression) bool {
	return e.Type() == "identifier" ||
		e.Type() == "numeric" ||
		e.Type() == "quoted-identifier" ||
		e.Type() == "prefix-operator" ||
		e.Type() == "group" ||
		e.Type() == "call"
}
