package internal

import (
	"fmt"
	"reflect"
	"strings"
)

/*
Extract all info of reflect.Type
After fetch all info in reflect.Type the meta info will be cache
Next call will get from cache instead of fetch again

# Example 1:

	type User struct {
			_ DbField[any] `db:"table(MyUser)"` //<-- Optional
			Id   DbField[uint64] `db:"primaryKey;autoIncrement"`
			Code DbField[string] `db:"unique;length(50)"`
			Name DbField[string] `db:"index;length(50)"`
		}

		return map {
				users: {
					id:{
						PrimaryKey    : true
					},
					code: {
							Unique:true,
							UniqueName: code_uk
					},
					name: {
							Index: true,
							IndexName; name_uk
					}
				}
		}
		# Exmaple 2
		type struct UserRole {
			Id   DbField[uint64]     `db:"primaryKey;autoIncrement"`
			UserId  DbField[uint64]  `db:"unique(user_role)"` //<-- unique constraint 2 columns
			RoleId  DbField[uint64]  `db:"unique(user_role)"` //<-- unique constraint 2 columns
		}
		return
		{
			users: {
				id: {
					PrimaryKey: true
				},
				user_id: {
					Unique: true,
					UniqueName: user_role_uk
				},
				role_id: {
					Unique: true,
					UniqueName: user_role_uk
				}
			}
		}
*/
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
		ftType := field.Type
		if ftType.Kind() == reflect.Ptr {
			ftType = ftType.Elem()
		}

		// Bỏ qua field đặc biệt "_" (used for table(...) override)
		if strings.HasPrefix(ftType.String(), u.entityTypeName) {

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
func (u *utilsPackage) getAutoPkKey(typ reflect.Type) *autoNumberKey {
	//check from cache
	if v, ok := u.cacheGetAutoPkKey.Load(typ.String()); ok {
		return v.(*autoNumberKey)
	}

	tableMap := u.GetMetaInfo(typ)
	for _, fields := range tableMap {
		for fieldName, fieldTag := range fields {
			if fieldTag.AutoIncrement {
				ret := &autoNumberKey{
					FieldName: fieldName,
					KeyType:   fieldTag.Field.Type,
					fieldTag:  &fieldTag,
				}
				u.cacheGetAutoPkKey.Store(typ.String(), ret)
				return ret
			}
		}
	}
	ret := &autoNumberKey{
		FieldName: "",
		KeyType:   nil,
	}
	u.cacheGetAutoPkKey.Store(typ.String(), ret)
	return ret
}
func (u *utilsPackage) extractValue(entityType reflect.Type, data interface{}) (*autoNumberKey, reflect.Type, map[string]interface{}, error) {
	ret := make(map[string]interface{})
	typ := reflect.TypeOf(data)
	valData := reflect.ValueOf(data)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		valData = valData.Elem()
	}
	tableMap := u.getRequireFields(entityType)
	for fieldName, fieldTag := range tableMap {
		fieldValue := valData.FieldByName(fieldName)
		if !fieldValue.IsValid() {
			dataFields := []string{}
			for i := 0; i < typ.NumField(); i++ {
				dataFields = append(dataFields, typ.Field(i).Name)
			}

			return nil, nil, nil, fmt.Errorf("%s is require but not found in %s/nFields: %s", fieldName, typ.String(), ToJsonString(dataFields))
		} else {
			val := fieldValue.Interface()
			if fieldTag.Length != nil {
				if fieldValue.Kind() == reflect.String {
					if len(val.(string)) == 0 {
						return nil, nil, nil, fmt.Errorf("%s is require but value is empty in %s/nFields: %d", fieldName, typ.String(), fieldTag.Length)
					}
					if len(val.(string)) > *fieldTag.Length {
						return nil, nil, nil, fmt.Errorf("size of '%s' in value of '%s' is exceed of %d", fieldName, typ.String(), *fieldTag.Length)
					}
				}

			}
			ret[fieldTag.Field.Name] = fieldValue.Interface()
		}

	}
	autoKey := u.getAutoPkKey(entityType)
	if autoKey.FieldName == "" {
		return nil, nil, ret, nil
	}
	if keyValueField, ok := typ.FieldByName(autoKey.fieldTag.Field.Name); ok {
		return autoKey, keyValueField.Type, ret, nil
	}
	return nil, nil, ret, nil

}
func (u *utilsPackage) getRequireFields(typ reflect.Type) map[string]FieldTag {
	//check cache
	if v, ok := u.cacheGetRequireFields.Load(typ.String()); ok {
		return v.(map[string]FieldTag)
	}
	tableMap := u.GetMetaInfo(typ)

	ret := make(map[string]FieldTag)
	for _, fields := range tableMap {
		for _, fieldTag := range fields {
			if (!fieldTag.AutoIncrement) && (!fieldTag.Nullable) && (fieldTag.Default == "") {
				ret[fieldTag.Field.Name] = fieldTag
			}
		}
	}
	u.cacheGetRequireFields.Store(typ.String(), ret)
	return ret

}
