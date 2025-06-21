package dbx

import "fmt"

func (c *Cfg) makeDnsMssql(dbname string) string {
	ret := ""
	if dbname == "" {
		if c.Port > 0 {
			ret = fmt.Sprintf("sqlserver://%s:%s@%s:%d", c.User, c.Password, c.Host, c.Port)
		} else {
			ret = fmt.Sprintf("sqlserver://%s:%s@%s", c.User, c.Password, c.Host)
		}
	} else {
		if c.Port > 0 {
			ret = fmt.Sprintf("sqlserver://%s:%s@%s:%d?database=%s", c.User, c.Password, c.Host, c.Port, dbname)
		} else {
			ret = fmt.Sprintf("sqlserver://%s:%s@%s?database=%s", c.User, c.Password, c.Host, dbname)
		}

	}
	return ret
}
