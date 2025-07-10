package orm

func (s *SqlStmt) String() string {
	ret := "SELECT " + s.Select + " FROM " + s.From
	if s.Where != "" {
		ret += " WHERE " + s.Where
	}
	if s.GroupBy != "" {
		ret += " GROUP BY " + s.GroupBy
	}
	if s.Having != "" {
		ret += " HAVING " + s.Having
	}
	if s.OrderBy != "" {
		ret += " ORDER BY " + s.OrderBy
	}
	if s.Limit != "" {
		ret += " LIMIT " + s.Limit
	}
	if s.Offset != "" {
		ret += " OFFSET " + s.Offset
	}
	return ret
}

func (s *SqlStmt) Err() error {
	return s.err
}
