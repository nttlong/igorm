package expr

import (
	"sort"
	"strings"
)

func (e *expression) InsertMarks(input string, markList [][]int) string {
	// Sort theo vị trí start giảm dần để tránh làm lệch chỉ số khi chèn
	sort.Slice(markList, func(i, j int) bool {
		return markList[i][0] > markList[j][0]
	})

	builder := strings.Builder{}
	builder.WriteString(input)

	for _, mark := range markList {
		start, end := mark[0], mark[1]
		if start < 0 || end > builder.Len() || start >= end {
			continue // tránh lỗi khi input sai
		}

		// Chèn dấu '`' vào cuối đoạn trước
		builderStr := builder.String()
		builder.Reset()
		builder.WriteString(builderStr[:end])
		builder.WriteString("`")
		builder.WriteString(builderStr[end:])

		builderStr = builder.String()
		builder.Reset()
		builder.WriteString(builderStr[:start])
		builder.WriteString("`")
		builder.WriteString(builderStr[start:])
	}

	return builder.String()
}
