package formatter

import (
	"strings"

	"github.com/DataDog/go-sqllexer"
)

type SelectStatement struct {
	Columns []SelectedColumn
}

func (ss SelectStatement) Format(indent int) string {
	formatted := padding(indent, 0) + "SELECT" + padding(FIRST_COLUMN_WIDTH-6, 1)
	indentTracker := len(formatted)

	for i := 0; i < len(ss.Columns); i++ {
		formatted += ss.Columns[i].Format(indentTracker)

		if i != len(ss.Columns)-1 {
			padding := padding(indent+FIRST_COLUMN_WIDTH, 0)
			formatted += ",\n" + padding
			indentTracker = len(padding)
		}
	}

	return formatted
}

type SelectedColumn struct {
	Exp   Expression
	Cast  TypecastExpression
	Alias sqllexer.Token
}

func (sc SelectedColumn) Format(indent int) string {
	formatted := ""
	if sc.Exp != nil {
		formatted += sc.Exp.String()
	}

	indentTracker := indent + len(formatted)

	if sc.Cast.Datatype.Datatype.Type != sqllexer.ERROR {
		castStr := padding(CAST_COLUMN_WIDTH-indentTracker, 1) + sc.Cast.String()
		indentTracker += len(castStr)
		formatted += castStr
	}

	// we should handle cases where the current length is larger than 120
	if sc.Alias.Type != sqllexer.ERROR {
		formatted += padding(SECOND_COLUMN_WIDTH-indentTracker, 1) + "AS " + strings.ToLower(sc.Alias.Value)
	}

	return formatted
}
