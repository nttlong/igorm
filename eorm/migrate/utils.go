package migrate

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

/*
this struct is used when DEV wants to make their struct as Model
example:

	type BaseModel struct {
		CreatedAt time.Time
	}
	type User struct {
		Entity `db:"table:users"` // <- if db tag not defined, table name is converted to SnakeCase of struct name and pluralized
		BaseModel
		ID int `db:"pk"`
	}
*/
type Entity struct {
	TableName string
	Cols      []ColumnDef          //<-- list of all columns
	MapCols   map[string]ColumnDef //<-- used for faster access to column by name
}

type ColumnDef struct {
	Name string

	Nullable bool
	/*
		if tag looks like this:
			`db:"pk"` or `db:primary` or `db:primaryKey`  => PKName = Name
			`db:"pk(<name of pk constraint>)"` or `db:primary(<name of pk constraint>)` or `db:primaryKey(<name of pk constraint>)` => PKName = <name of pk constraint>
		else:
			PKName = ""
	*/
	PKName  string
	IsAuto  bool
	Default string
	/*
		if tag looks like this:
			`db:"uk"` or `db:unique`   => UniqueName = Name
			`db:"uk(<name of the unique index>)"` or `db:unique(<name of the unique index>)` => UniqueName = <name of the unique index>
		else:
			UniqueName = ""
	*/
	UniqueName string
	/*
		if tag looks like this:
			`db:"idx"` or `db:index`   => IndexName = Name
			`db:"idx(<name of the index>)"` or `db:index(<name of the index>)` => IndexName = <name of the index>
		else:
			IndexName = ""
	*/
	IndexName string
	Field     reflect.StructField
}

type utils struct{}

var dbTagPattern = regexp.MustCompile(`([a-zA-Z]+)(\((.*?)\))?`)

/*
this function will parse the db tag of a field
if tag not found, it will return ColumnDef with default values
look at the example below:

	ColumnDef {
		Name:     name of the field, // the other info is default
	}
*/
func (u *utils) ParseTagFromStruct(field reflect.StructField) ColumnDef {
	tag := field.Tag.Get("db")
	name := field.Name
	col := ColumnDef{
		Name:     u.SnakeCase(name),
		Field:    field,
		Nullable: field.Type.Kind() == reflect.Ptr, // auto derive from Go type
	}

	if tag == "" {
		return col
	}

	parts := strings.Split(tag, ",")
	for _, part := range parts {
		m := dbTagPattern.FindStringSubmatch(part)
		if len(m) < 2 {
			continue
		}
		key := strings.ToLower(m[1])
		val := ""
		if len(m) >= 4 {
			val = m[3]
		}
		switch key {
		case "pk", "primary", "primarykey":
			col.PKName = val
			if col.PKName == "" {
				col.PKName = col.Name
			}
		case "uk", "unique":
			col.UniqueName = val
			if col.UniqueName == "" {
				col.UniqueName = col.Name
			}
		case "idx", "index":
			col.IndexName = val
			if col.IndexName == "" {
				col.IndexName = col.Name
			}
		case "auto":
			col.IsAuto = true
		case "default":
			col.Default = val
		case "column":
			col.Name = val
		}
	}

	return col
}

/*
This function will parse the struct tag and return a slice of ColumnDef
Example:

	type BaseModel struct {
		CreatedAt time.Time
	}
	type User struct {
		Entity `db:"table:users"` // <- if db tag not defined, table name is converted to SnakeCase of struct name and pluralized
		BaseModel
		ID int `db:"pk"`
	}
*/
func (u *utils) ParseStruct(obj interface{}) ([]ColumnDef, error) {
	t := reflect.TypeOf(obj)
	if t.Kind() != reflect.Struct {
		return nil, fmt.Errorf("only struct is supported")
	}

	var cols []ColumnDef
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if f.Anonymous {
			if f.Type == reflect.TypeOf(Entity{}) {
				continue // skip embedded Entity
			}
			subCols, _ := u.ParseStruct(reflect.New(f.Type).Elem().Interface())
			cols = append(cols, subCols...)
			continue
		}
		cols = append(cols, u.ParseTagFromStruct(f))
	}
	return cols, nil
}

/*
this function will return a map of primary key columns with their names
@return map[string][]ColumnDef //<-- map[Primary key constraint name][]ColumnDef
*/
func (u *utils) GetPrimaryKey(e *Entity) map[string][]ColumnDef {
	m := make(map[string][]ColumnDef)
	for _, col := range e.Cols {
		if col.PKName != "" {
			m[col.PKName] = append(m[col.PKName], col)
		}
	}
	return m
}

/*
this function will return a map of unique key columns with their names
@return map[string][]ColumnDef //<-- map[Unique constraint name][]ColumnDef
*/
func (u *utils) GetUnique(e *Entity) map[string][]ColumnDef {
	m := make(map[string][]ColumnDef)
	for _, col := range e.Cols {
		if col.UniqueName != "" {
			m[col.UniqueName] = append(m[col.UniqueName], col)
		}
	}
	return m
}

/*
this function will return a map of index columns with their names
@return map[string][]ColumnDef //<-- map[Index constraint name][]ColumnDef
*/
func (u *utils) GetIndex(e *Entity) map[string][]ColumnDef {
	m := make(map[string][]ColumnDef)
	for _, col := range e.Cols {
		if col.IndexName != "" {
			m[col.IndexName] = append(m[col.IndexName], col)
		}
	}
	return m
}

/*
this function will return a string of unsynced columns
unsynced columns are columns in Entity that do not exist in the db table
@dbColumnName is a list of column names in the db table
*/
func (u *utils) GetUnSyncColumns(e *Entity, dbColumnName []string) string {
	dbCols := make(map[string]bool)
	for _, c := range dbColumnName {
		dbCols[c] = true
	}

	var unsync []string
	for _, col := range e.Cols {
		if _, found := dbCols[col.Name]; !found {
			unsync = append(unsync, col.Name)
		}
	}
	return strings.Join(unsync, ", ")
}

func (u *utils) SnakeCase(s string) string {
	var out []rune
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			out = append(out, '_')
		}
		out = append(out, r)
	}
	return strings.ToLower(string(out))
}

func (u *utils) Pluralize(word string) string {
	if strings.HasSuffix(word, "s") {
		return word
	}
	return word + "s"
}
