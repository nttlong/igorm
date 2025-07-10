package expr

// func (e *expression) PrepareOld(input string) (string, error) {
// 	if e.keywords == nil {
// 		e.keywords = []string{
// 			"select",
// 			"from",
// 			"where",
// 			"group",
// 			"order",
// 			"limit",
// 			"offset",
// 		}
// 	}
// 	for _, keyword := range e.keywords {
// 		markList, err := e.GetMarkList(input, keyword)
// 		if err != nil {
// 			return "", err
// 		}
// 		input = e.InsertMarks(input, markList)
// 	}
// 	return input, nil

// }
// func (e *expression) Prepare(sql string) string {
// 	reField := regexp.MustCompile(`\b([A-Za-z_][A-Za-z0-9_]*)\.([A-Za-z_][A-Za-z0-9_]*)\b`)
// 	sql = reField.ReplaceAllString(sql, "`$1`.`$2`")

// 	// 2. Match FROM, JOIN ... table names (employees, orders, etc.)
// 	reTable := regexp.MustCompile(`(?i)(FROM|JOIN)\s+([A-Za-z_][A-Za-z0-9_]*)\b`)
// 	sql = reTable.ReplaceAllString(sql, "$1 `$2`")

// 	// 3. Match AS aliases: AS T1 → AS `T1`
// 	reAlias := regexp.MustCompile(`(?i)\bAS\s+([A-Za-z_][A-Za-z0-9_]*)\b`)
// 	sql = reAlias.ReplaceAllString(sql, "AS `$1`")

// 	return sql
// }
