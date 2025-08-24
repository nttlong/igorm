package migrate

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type mUtils struct {
}

// if s is "true" or "false" retun true
func (m *mUtils) isBool(s string) bool {
	return strings.ToLower(s) == "true" || strings.ToLower(s) == "false"
}
func (m *mUtils) isNumber(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}
func (m *mUtils) isFloatNumber(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

func (m *mUtils) GetDefaultValue(defaultValue string, defaultValueByFromDbTag map[string]string) (string, error) {
	if strings.Contains(defaultValue, "'") {
		return defaultValue, nil
	}
	if typeUtils.isFloatNumber(defaultValue) {
		return defaultValue, nil

	} else if typeUtils.isNumber(defaultValue) {
		return defaultValue, nil

	} else if typeUtils.isBool(defaultValue) {
		return defaultValue, nil

	} else if val, ok := defaultValueByFromDbTag[defaultValue]; ok {
		return val, nil
	} else {
		return "", fmt.Errorf("not support default value from %s, review GetGetDefaultValueByFromDbTag() function in %s", defaultValue, reflect.TypeOf(m).Elem())
	}
}

var typeUtils = &mUtils{}
