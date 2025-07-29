package vdb

import "strings"

func (d *mySqlDialect) SqlFunction(delegator *DialectDelegateFunction) (string, error) {
	fnName := strings.ToLower(delegator.FuncName)
	switch fnName {
	case "now":
		delegator.HandledByDialect = true
		return "NOW()", nil
	case "len":
		delegator.FuncName = "LENGTH"
		return delegator.FuncName, nil
	default:
		return "", nil
	}

}
