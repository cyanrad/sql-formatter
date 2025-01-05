package formatter

import "github.com/DataDog/go-sqllexer"

type SelectStatement struct {
	Columns []SelectedColumn
}

func (ss SelectStatement) Format(indent int) string {
	formatted := padding(indent, 0) + "SELECT" + padding(FIRST_COLUMN_WIDTH-6, 1)
	indentTracker := len(formatted)

	for i, c := range ss.Columns {
		formatted += c.Format(indentTracker)

		if i != len(ss.Columns)-1 && len(c.Exps) > 0 { // second condition is in case of comment line
			formatted += ","
		}

		if c.Comment.String() != "" {
			if len(c.Exps) > 0 {
				formatted += padding(1, 0)
			}
			formatted += c.Comment.String()
		}

		if i != len(ss.Columns)-1 {
			padding := padding(indent+FIRST_COLUMN_WIDTH, 0)
			formatted += "\n" + padding
			indentTracker = len(padding)
		}
	}

	return formatted
}

type SelectedColumn struct {
	Exps     []Expression
	Cast     TypecastExpression
	Alias    AsExpression
	Comment  CommentExpression
	hasAlias bool
	hasCast  bool
}

func (sc SelectedColumn) Format(indent int) string {
	formatted := ""

	prevType := "operator"
	for _, e := range sc.Exps {
		if !(prevType == "operator") && e.Type() != "operator" { // this is so fucking bad holy shit - edit: i genuinly have no clue what I was thinking when I was writing this but holy shit it's so simple and solve a more complicated problem
			formatted += " "
		}

		// formatted += padExpression(e, i, len(sc.Exps))
		formatted += e.String()
		prevType = e.Type()
	}

	indentTracker := indent + len(formatted)

	if sc.hasCast {
		castStr := padding(CAST_COLUMN_WIDTH-indentTracker, 1) + sc.Cast.String()
		indentTracker += len(castStr)
		formatted += castStr
	}

	// we should handle cases where the current length is larger than 120
	if sc.hasAlias {
		formatted += padding(SECOND_COLUMN_WIDTH-indentTracker, 1) + sc.Alias.String()
	}

	return formatted
}

func (sc *SelectedColumn) HandlePotentialComment(f *Formatter) {
	if f.currTypeIs(sqllexer.COMMENT) {
		sc.Comment = CommentExpression{Token: f.currToken}
		f.nextToken()
	}
}

// this is very fucking stupid holy shit
// like actually this might be some of the worst code I've every wrote
// a toddler can come up with better logic
// func padExpression(e Expression, i int, expsCount int) string {
// 	fmt.Println(e.Type() + e.String())
// 	if e.Type() == "operator" {
// 		switch e.String() {
// 		case ",":
// 			return e.String() + " "
// 		case ".", "(", ")", "{", "}", "[", "]":
// 			return e.String()
// 		default:
// 			return " " + e.String() + " "
// 		}
// 	} else {
// 		if i != expsCount-1 {
// 			return e.String() + " "
// 		}
// 		return e.String()
// 	}

// 	return "FUCK"
// }
