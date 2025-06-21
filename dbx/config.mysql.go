package dbx

import "fmt"

func (c *Cfg) makeDnsMySql(dbname string) string {
	ret := ""
	if dbname == "" {
		ret = fmt.Sprintf("%s:%s@tcp(%s:%d)/?multiStatements=true&parseTime=true&loc=Local", c.User, c.Password, c.Host, c.Port)
	} else {
		ret = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?multiStatements=true&parseTime=true&loc=Local", c.User, c.Password, c.Host, c.Port, dbname)
	}
	return ret
}
