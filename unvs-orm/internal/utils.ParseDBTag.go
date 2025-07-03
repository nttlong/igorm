package internal

import (
	"reflect"
	"strconv"
	"strings"
)

type FieldTag struct {
	PrimaryKey    bool
	AutoIncrement bool
	Unique        bool
	/*
		can be field name if no unique index name in tag else name of unique index in tag
	*/
	UniqueName string
	Index      bool
	/*
		can be field name if no  index name in tag else name of  index in tag
	*/
	IndexName string
	Length    *int
	FTSName   string
	DBType    string
	TableName string
	Check     string
	Nullable  bool
	Field     reflect.StructField
	Default   string
}

// ParseDBTag parses the `db` struct tag into a FieldTag struct.
func (u *utilsPackage) ParseDBTag(field reflect.StructField) FieldTag {

	tag := strings.TrimSpace(field.Tag.Get("db"))
	t := FieldTag{
		Field: field,
	}

	tag = strings.TrimSpace(tag)
	if tag == "" {
		return t
	}
	parts := strings.Split(tag, ";")
	for _, p := range parts {
		p = strings.TrimSpace(p)
		switch {
		case p == "nullable":
			t.Nullable = true
		case p == "null":
			t.Nullable = true
		case p == "notnull":
			t.Nullable = false
		case p == "primaryKey":
			t.PrimaryKey = true
		case p == "autoIncrement":
			t.AutoIncrement = true
		case p == "unique":
			t.Unique = true
			t.UniqueName = u.ToSnakeCase(field.Name)
		case strings.HasPrefix(p, "unique("):
			t.Unique = true
			t.UniqueName = u.extractName(p)
		case p == "index":
			t.Index = true
			t.IndexName = u.ToSnakeCase(field.Name)
		case strings.HasPrefix(p, "index("):
			t.Index = true
			t.IndexName = u.extractName(p)
		case strings.HasPrefix(p, "table("):
			t.TableName = u.extractName(p)
		case strings.HasPrefix(p, "length("):
			if s := u.extractName(p); s != "" {
				if n, err := strconv.Atoi(s); err == nil {
					t.Length = &n
				}
			}
		case strings.HasPrefix(p, "check("):
			t.Check = u.extractName(p)

			if s := u.extractName(p); s != "" {
				if n, err := strconv.Atoi(s); err == nil {
					t.Length = &n
				}
			}
		case strings.HasPrefix(p, "FTS("):
			t.FTSName = u.extractName(p)
		case strings.HasPrefix(p, "type:"):
			t.DBType = strings.TrimPrefix(p, "type:")
		case strings.HasPrefix(p, "default:"):
			t.Default = strings.TrimPrefix(p, "default:")
		}

	}
	return t
}
