package migrate

import "strconv"

type mUtils struct {
}

func (m *mUtils) isNumber(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}
func (m *mUtils) isFloatNumber(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

var typeUtils = &mUtils{}
