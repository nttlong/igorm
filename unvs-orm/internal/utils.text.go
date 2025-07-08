package internal

import (
	"strconv"
	"strings"
)

func (u *utilsPackage) GetTableNameFromVirtualName(VirtualName string) string {
	//check form cache
	if v, ok := u.cacheGetTableNameFromVirtualName.Load(VirtualName); ok {
		return v.(string)
	}
	ret := strings.Split(VirtualName, "*")[0]
	u.cacheGetTableNameFromVirtualName.Store(VirtualName, ret)
	return ret
}

func (u *utilsPackage) Contains(list []string, item string) bool {
	item = strings.ToLower(item)
	for _, v := range list {
		if strings.ToLower(v) == item {
			return true
		}
	}
	return false
}

/*
this function will find any ? in the sql and replace it with placeholder+found index in the args
for example:
sql: "SELECT * FROM table WHERE id =? AND name =?"
placeholder: $
return: "SELECT * FROM table WHERE id = $1 AND name = $2"
*/
func (u *utilsPackage) replacePlaceHolder(placeholder string, sql string) string {
	key := placeholder + ":" + sql
	if v, ok := u.cacheReplacePlaceHolder.Load(key); ok {
		return v.(string)
	}
	var result strings.Builder
	paramIndex := 1
	for i := 0; i < len(sql); i++ {
		if sql[i] == '?' {
			result.WriteString(placeholder)
			result.WriteString(strconv.Itoa(paramIndex))
			paramIndex++
		} else {
			result.WriteByte(sql[i])
		}
	}
	ret := result.String()
	u.cacheReplacePlaceHolder.Store(key, ret)
	return ret

}
func (u *utilsPackage) Join(parts []string, sep string) string {
	return u.join(parts, sep)
}
func (u *utilsPackage) join(parts []string, sep string) string {
	out := ""
	for i, p := range parts {
		if i > 0 {
			out += sep
		}
		out += p
	}
	return out
}
