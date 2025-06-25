package unvsef

import (
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"unicode"

	"github.com/jinzhu/inflection"
)

// --------------------- Tag Metadata ---------------------
// FieldTag holds parsed metadata from struct field tags.
type FieldTag struct {
	PrimaryKey    bool
	AutoIncrement bool
	Unique        bool
	UniqueName    string
	Index         bool
	IndexName     string
	Length        *int
	FTSName       string
	DBType        string
	TableName     string
	Check         string
	Nullable      bool
	Field         reflect.StructField
	Default       string
}
type utilsPackage struct {
	cacheGetMetaInfo         sync.Map
	CacheTableNameFromStruct sync.Map
	cacheGetPkFromMeta       sync.Map
	// future: add cache or shared state here
}

var utils = &utilsPackage{}

// ParseDBTag parses the `db` struct tag into a FieldTag struct.
func (u *utilsPackage) ParseDBTag(field reflect.StructField) FieldTag {

	tag := strings.TrimSpace(field.Tag.Get("db"))
	t := FieldTag{
		Field: field,
	}

	t.Nullable = strings.Contains(field.Type.String(), ".DbField[*")
	tag = strings.TrimSpace(tag)
	if tag == "" {
		return t
	}
	parts := strings.Split(tag, ";")
	for _, p := range parts {
		p = strings.TrimSpace(p)
		switch {

		case p == "primaryKey":
			t.PrimaryKey = true
		case p == "autoIncrement":
			t.AutoIncrement = true
		case p == "unique":
			t.Unique = true
		case strings.HasPrefix(p, "unique("):
			t.Unique = true
			t.UniqueName = u.extractName(p)
		case p == "index":
			t.Index = true
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

func (u *utilsPackage) extractName(s string) string {
	re := regexp.MustCompile(`\((.*?)\)`)
	matches := re.FindStringSubmatch(s)
	if len(matches) == 2 {
		return matches[1]
	}
	return ""
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
func (u *utilsPackage) ToSnakeCase(str string) string {
	var result []rune
	for i, r := range str {
		if i > 0 && unicode.IsUpper(r) &&
			(unicode.IsLower(rune(str[i-1])) || (i+1 < len(str) && unicode.IsLower(rune(str[i+1])))) {
			result = append(result, '_')
		}
		result = append(result, unicode.ToLower(r))
	}
	return string(result)
}
func (u *utilsPackage) TableNameFromStruct(typ reflect.Type) string {
	// Check for override via table(...) tag
	if v, ok := u.CacheTableNameFromStruct.Load(typ); ok {
		return v.(string)
	}
	for i := 0; i < typ.NumField(); i++ {
		if typ.Field(i).Name == "_" {
			parsed := u.ParseDBTag(typ.Field(i))
			if parsed.TableName != "" {
				return parsed.TableName
			}
		}
	}
	base := u.ToSnakeCase(typ.Name())
	ret := inflection.Plural(base)
	u.CacheTableNameFromStruct.Store(typ, ret)
	return ret
}
func (u *utilsPackage) ResolveFieldKind(field reflect.StructField) string {
	strFt := field.Type.String()

	if strings.Contains(strFt, ".DbField[") {
		typeParam := strings.Split(strFt, ".DbField[")[1]
		typeParam = strings.Split(typeParam, "]")[0]
		typeParam = strings.TrimLeft(typeParam, "*")
		return typeParam
	}
	return field.Type.Kind().String()
}

// GetMetaInfo extracts and caches parsed FieldTag metadata for a given struct type.
// Lấy và cache metadata của các field trong struct theo dạng FieldTag.
func (u *utilsPackage) GetMetaInfo(typ reflect.Type) map[string]map[string]FieldTag {
	// 1. Kiểm tra cache trước (check cache first)
	if metaInfo, ok := u.cacheGetMetaInfo.Load(typ); ok {
		return metaInfo.(map[string]map[string]FieldTag)
	}

	// 2. Tạo mới metadata
	metaInfo := make(map[string]map[string]FieldTag)
	tableName := utils.TableNameFromStruct(typ)
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)

		// Bỏ qua field đặc biệt "_" (used for table(...) override)
		if field.Name == "_" {

			continue
		}

		// 3. Nếu là embedded struct (anonymous), đệ quy lấy metadata của nó
		if field.Anonymous {
			embeddedMeta := u.GetMetaInfo(field.Type)
			for tableName, fields := range embeddedMeta {
				if _, ok := metaInfo[tableName]; !ok {
					metaInfo[tableName] = make(map[string]FieldTag)
				}
				for fieldName, fieldTag := range fields {

					metaInfo[tableName][u.ToSnakeCase(fieldName)] = fieldTag
				}
			}
			continue
		}

		if _, ok := metaInfo[tableName]; !ok {
			metaInfo[tableName] = make(map[string]FieldTag)
		}

		// 5. Gán tag metadata cho field
		metaInfo[tableName][u.ToSnakeCase(field.Name)] = u.ParseDBTag(field)
	}

	// 6. Cache lại và trả về
	u.cacheGetMetaInfo.Store(typ, metaInfo)
	return metaInfo
}
func (u *utilsPackage) Quote(strQuote string, str ...string) string {
	left := strQuote[0:1]
	right := strQuote[1:2]
	ret := left + strings.Join(str, left+"."+right) + right
	return ret
}
func (u *utilsPackage) GetPkFromMetaByType(typ reflect.Type) map[string]map[string]FieldTag {
	// 1. Kiểm tra cache trước (check cache first)
	if pk, ok := u.cacheGetPkFromMeta.Load(typ); ok {
		return pk.(map[string]map[string]FieldTag)
	}
	metaInfo := u.GetMetaInfo(typ)
	ret := make(map[string]map[string]FieldTag)

	for tableName, fields := range metaInfo {
		ret[tableName+"_pk"] = make(map[string]FieldTag)
		for fieldName, fieldTag := range fields {
			if fieldTag.PrimaryKey {
				ret[tableName+"_pk"][fieldName] = fieldTag
			}
		}
	}
	return ret
}
