package db

import (
	"fmt"
	"regexp"
	"strings"
)

type SearchStmt struct {
	Select []string
	From   string
	Where  []string
}

func (s SearchStmt) String() string {
	var sb strings.Builder

	var selectStmt string
	if len(s.Select) == 0 {
		selectStmt = "*"
	} else {
		selectStmt = strings.Join(s.Select, ", ")
	}

	sb.WriteString(fmt.Sprintf("select %v from %v", selectStmt, s.From))

	if len(s.Where) > 0 {
		sb.WriteString("\nwhere")
	}
	newLine := false
	for _, where := range s.Where {
		sb.WriteString("\n")
		if newLine {
			sb.WriteString("and " + where)
		} else {
			sb.WriteString(where)
			newLine = true
		}
	}

	return sb.String()
}

// Roughly convert regex to sql-like pattern
// This does not support escape characters in regex
func ConvertRegexToLikePattern(inputRegex string) string {
	escaped := regexp.MustCompile(`([\\%_])`).ReplaceAllString(inputRegex, `\$1`)

	escaped = strings.ReplaceAll(escaped, "_", "\\_")
	escaped = strings.ReplaceAll(escaped, "%", "\\%")
	escaped = strings.ReplaceAll(escaped, ".", "_")
	escaped = strings.ReplaceAll(escaped, "*", "%")
	escaped = strings.ReplaceAll(escaped, "+", "%")
	escaped = strings.ReplaceAll(escaped, "?", "_")

	return escaped
}
