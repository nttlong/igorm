package expr

import "bytes"

func (e *expression) IsSpecialChar(input byte) bool {
	if e.specialChar == nil {
		e.specialChar = []byte{
			'(', ')', ',', '.', ' ', '/', '+', '-', '*', '%', '=', '<', '>', '!', '&', '|', '^', '~', '?', ':', ';', '[', ']', '{', '}', '@', '#', '$',
		}
	}
	return bytes.Contains(e.specialChar, []byte{input})
}
