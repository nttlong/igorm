/*This package contains the configuration for the database connection.*/
package dbx

import "fmt"

func (c *Cfg) makeDnsPostgres(dbname string) string {
	ret := ""
	if c.Password != "" {
		if c.SSL {

			if dbname == "" {
				ret = fmt.Sprintf("postgres://%s:%s@%s:%d", c.User, c.Password, c.Host, c.Port)
			} else {
				ret = fmt.Sprintf("postgres://%s:%s@%s:%d/%s", c.User, c.Password, c.Host, c.Port, dbname)
			}
		} else {
			if dbname == "" {
				ret = fmt.Sprintf("postgres://%s:%s@%s:%d?sslmode=disable", c.User, c.Password, c.Host, c.Port)
			} else {
				ret = fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", c.User, c.Password, c.Host, c.Port, dbname)
			}
		}
	} else {
		if c.SSL {
			if dbname == "" {
				ret = fmt.Sprintf("postgres://%s@%s:%d", c.User, c.Host, c.Port)
			} else {
				ret = fmt.Sprintf("postgres://%s@%s:%d/%s", c.User, c.Host, c.Port, dbname)
			}
		} else {
			if dbname == "" {
				ret = fmt.Sprintf("postgres://%s@%s:%d?sslmode=disable", c.User, c.Host, c.Port)
			} else {
				ret = fmt.Sprintf("postgres://%s@%s:%d/%s?sslmode=disable", c.User, c.Host, c.Port, dbname)
			}
		}
	}
	return ret
}
