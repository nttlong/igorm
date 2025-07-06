package expr

import (
	"strings"
)

func (e *expression) GetMarkList(input string, keyword string) ([][]int, error) {
	keywordLower := strings.ToLower(keyword)
	if !strings.Contains(strings.ToLower(input), keywordLower) {
		return nil, nil
	}

	if len(keywordLower) == 0 || len(input) < len(keywordLower) {
		return nil, nil
	}

	markList := make([][]int, 0, 4) // preallocate
	inputBytes := []byte(input)
	keywordBytes := []byte(keywordLower)
	keywordLen := len(keywordBytes)

	i := 0
	for i <= len(inputBytes)-keywordLen {
		c := inputBytes[i]
		if e.IsSpecialChar(c) {
			i++
			continue
		}

		// Kiểm tra từ vị trí i có khớp keyword không (case-insensitive)
		match := true
		for j := 0; j < keywordLen; j++ {
			if i+j >= len(inputBytes) || toLowerAscii(inputBytes[i+j]) != keywordBytes[j] {
				match = false
				break
			}
		}

		if match {
			start := i
			end := i + keywordLen

			// Đảm bảo tiếp sau là special char hoặc kết thúc
			if end < len(inputBytes) && !e.IsSpecialChar(inputBytes[end]) {
				i++
				continue
			}

			markList = append(markList, []int{start, end})
			i = end
		} else {
			i++
		}
	}
	return markList, nil
}

// Hỗ trợ: chỉ lowercase ASCII → nhanh và không tạo string mới
func toLowerAscii(b byte) byte {
	if b >= 'A' && b <= 'Z' {
		return b + 32
	}
	return b
}
