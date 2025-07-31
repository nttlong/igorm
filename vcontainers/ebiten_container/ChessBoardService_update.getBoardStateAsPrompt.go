package ebiten_container

import (
	"fmt"
	"strings"
)

// getBoardStateAsPrompt tạo một chuỗi mô tả trạng thái bàn cờ cho AI
func (g *ChessBoardService) getBoardStateAsPrompt() string {
	var sb strings.Builder
	sb.WriteString("Current Xiangqi board state:\n")
	sb.WriteString("Format: (PieceName)_[Color] (x,y)\n")
	sb.WriteString("Colors: Black (Đen), Red (Đỏ)\n")
	sb.WriteString("Pieces: 車(Xe), 馬(Mã), 象(Tượng), 士(Sĩ), 將(Tướng), 砲(Pháo), 卒(Tốt), 俥(Xe), 傌(Mã), 相(Tượng), 帥(Tướng), 炮(Pháo), 兵(Tốt)\n")
	sb.WriteString("Board coordinates are (column, row) from (0,0) top-left to (8,9) bottom-right.\n\n")

	pieceDescriptions := []string{}

	for y := 0; y < 10; y++ {
		for x := 0; x < 9; x++ {
			piece := g.board[y][x]
			if piece != "" {
				colorName := ""
				pieceColor := g.getPieceColor(x, y)
				if pieceColor == 1 {
					colorName = "Đen"
				} else if pieceColor == 2 {
					colorName = "Đỏ"
				}

				pieceDescriptions = append(pieceDescriptions, fmt.Sprintf("%s_%s (%d,%d)", piece, colorName, x, y))
			}
		}
	}
	sb.WriteString("Pieces on board:\n")
	sb.WriteString(strings.Join(pieceDescriptions, ", ") + "\n\n")

	return sb.String()
}
