package orm

import "bytes"

func (e *expression) isSpecialChar(input byte) bool {
	if e.specialChar == nil {
		e.specialChar = []byte{
			'(', ')', ',', '.', ' ', '/', '+', '-', '*', '%', '=', '<', '>', '!', '&', '|', '^', '~', '?', ':', ';', '[', ']', '{', '}', '@', '#', '$',
		}
	}
	return bytes.Contains(e.specialChar, []byte{input})
}
