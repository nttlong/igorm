package dbx

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"

	"github.com/go-sql-driver/mysql"
)

type mysqlErrorParser struct {
}

func (p *mysqlErrorParser) ParseError(ctx context.Context, db *sql.DB, err error) error {
	if err == nil {
		return nil
	}
	if mysqlErr, ok := err.(*mysql.MySQLError); ok {
		if mysqlErr.Number == 1062 {
			return p.ParseError1062(ctx, db, *mysqlErr)
		}

		strMySQLErrNumber := strconv.FormatUint(uint64(mysqlErr.Number), 10)
		panic(fmt.Sprintf("not implemented  func (p *mysqlErrorParser) ParseError in dbxerrors.mysql,error code %s", strMySQLErrNumber))
	}

	panic("not implemented  func (p *mysqlErrorParser) ParseError in dbxerrors.mysql")
}

var mySqlErrorParser *mysqlErrorParser

func init() {
	mySqlErrorParser = &mysqlErrorParser{}
}
