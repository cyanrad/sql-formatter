package formatter

import (
	"strings"

	"github.com/DataDog/go-sqllexer"
)

const (
	FIRST_COLUMN_WIDTH  = 8
	SECOND_COLUMN_WIDTH = 120
	CAST_COLUMN_WIDTH   = 102
	LEVEL_INDENT        = 4
)

func padding(indent int, least int) string {
	if indent <= 0 {
		if least > 0 {
			return strings.Repeat(" ", least)
		}

		return ""
	}

	return strings.Repeat(" ", indent)
}

func (f *Formatter) currTokenIs(t sqllexer.TokenType, v string) bool {
	return f.currToken.Type == t && strings.ToUpper(f.currToken.Value) == v
}
func (f *Formatter) peekTokenIs(t sqllexer.TokenType, v string) bool {
	return f.peekToken.Type == t && strings.ToUpper(f.peekToken.Value) == v
}

func (f *Formatter) currTypeIs(t sqllexer.TokenType) bool {
	return f.currToken.Type == t
}
func (f *Formatter) peekTypeIs(t sqllexer.TokenType) bool {
	return f.peekToken.Type == t
}

func clearWhiteSpace(tokens []sqllexer.Token) []sqllexer.Token {
	newTokens := []sqllexer.Token{}

	for _, t := range tokens {
		if !(t.Type == sqllexer.WS) {
			newTokens = append(newTokens, t)
		}
	}

	return newTokens
}

var keywords = map[string]struct{}{
	"SELECT":    {},
	"FROM":      {},
	"WHERE":     {},
	"GROUP":     {},
	"ORDER":     {},
	"LEFT":      {},
	"RIGHT":     {},
	"INNER":     {},
	"OUTER":     {},
	"JOIN":      {},
	"ON":        {},
	"CASE":      {},
	"WHEN":      {},
	"THEN":      {},
	"ELSE":      {},
	"END":       {},
	"HAVING":    {},
	"INSERT":    {},
	"UPDATE":    {},
	"DELETE":    {},
	"CREATE":    {},
	"ALTER":     {},
	"DROP":      {},
	"TRUNCATE":  {},
	"MERGE":     {},
	"UNION":     {},
	"ALL":       {},
	"DISTINCT":  {},
	"EXCEPT":    {},
	"INTERSECT": {},
	"NULL":      {},
	"NOT":       {},
	"LIKE":      {},
	"IN":        {},
	"BETWEEN":   {},
	"AND":       {},
	"OR":        {},
	"IS":        {},
	"EXISTS":    {},
	"AS":        {},
	"INTO":      {},
	"VALUES":    {},
	"LIMIT":     {},
	"OFFSET":    {},
	"FETCH":     {},
	"FOR":       {},
	"BY":        {},
	"ASC":       {},
	"DESC":      {},
	"WITH":      {},
	"QUALIFY":   {},
}

func isKeyword(t sqllexer.Token) bool {
	_, ok := keywords[strings.ToUpper(t.Value)]
	return ok && t.Type == sqllexer.IDENT
}
