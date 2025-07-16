package dbv

import (
	"fmt"
	"regexp"
	"strings"
)

type JoinClause struct {
	JoinType string
	Table    string
	Alias    string
	On       string
}

type DslJoin struct {
	FromAlias string
	FromTable string
	Joins     []JoinClause
}

func ParseDslJoin(dsl string) (*DslJoin, error) {
	re := regexp.MustCompile(`\(([^)]+)\)\s*(<-|->|<->|-)\s*(.+)`)
	matches := re.FindStringSubmatch(dsl)
	if len(matches) != 4 {
		return nil, fmt.Errorf("invalid DSL format")
	}

	aliasDefs := matches[1]
	joinOp := matches[2]
	on := matches[3]

	joinTypeMap := map[string]string{
		"<->": "FULL JOIN",
		"->":  "LEFT JOIN",
		"<-":  "RIGHT JOIN",
		"-":   "INNER JOIN",
	}
	joinType := joinTypeMap[joinOp]

	// Parse alias:table pairs, keep order
	parts := strings.Split(aliasDefs, ",")
	if len(parts) < 2 {
		return nil, fmt.Errorf("at least 2 tables required")
	}

	aliases := []struct {
		Alias string
		Table string
	}{}
	for _, part := range parts {
		p := strings.Split(strings.TrimSpace(part), ":")
		if len(p) != 2 {
			return nil, fmt.Errorf("invalid alias format: %s", part)
		}
		aliases = append(aliases, struct {
			Alias string
			Table string
		}{
			Alias: strings.TrimSpace(p[0]),
			Table: strings.TrimSpace(p[1]),
		})
	}

	return &DslJoin{
		FromAlias: aliases[0].Alias,
		FromTable: aliases[0].Table,
		Joins: []JoinClause{
			{
				JoinType: joinType,
				Table:    aliases[1].Table,
				Alias:    aliases[1].Alias,
				On:       on,
			},
		},
	}, nil
}

func (j *DslJoin) JoinString() string {
	var sb strings.Builder
	for _, join := range j.Joins {
		sb.WriteString(fmt.Sprintf(" %s %s AS %s ON %s",
			join.JoinType, join.Table, join.Alias, join.On))
	}
	return sb.String()
}

//	func (j *DslJoin) FullSQL() string {
//		return fmt.Sprintf("FROM %s AS %s%s", j.FromTable, j.FromAlias, j.JoinString())
//	}
func (j *DslJoin) FullSQL() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("FROM %s AS %s", j.FromTable, j.FromAlias))
	for _, join := range j.Joins {
		sb.WriteString(fmt.Sprintf(" %s %s AS %s ON %s",
			join.JoinType, join.Table, join.Alias, join.On))
	}
	return sb.String()
}
