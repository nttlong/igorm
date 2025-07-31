package ebiten_container

func (g *ChessBoardService) getPossibleMoves1(x, y int) []Point {
	var moves []Point
	piece := g.board[y][x]

	switch piece {
	case "卒", "兵":
		// Tốt đen: chỉ đi xuống (tăng y), không lùi
		if piece == "卒" {
			if g.isInsideBoard(x, y+1) && g.getPieceColor(x, y+1) != g.getPieceColor(x, y) {
				moves = append(moves, Point{X: x, Y: y + 1})
			}
		}
		// Tốt đỏ: chỉ đi lên (giảm y), không lùi
		if piece == "兵" {
			if g.isInsideBoard(x, y-1) && g.getPieceColor(x, y-1) != g.getPieceColor(x, y) {
				moves = append(moves, Point{X: x, Y: y - 1})
			}
		}

	case "車", "俥":
		// Xe: đi thẳng các hướng đến khi gặp quân
		dirs := []Point{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}
		for _, dir := range dirs {
			nx, ny := x+dir.X, y+dir.Y
			for g.isInsideBoard(nx, ny) {
				if g.board[ny][nx] == "" {
					moves = append(moves, Point{X: nx, Y: ny})
				} else {
					// Nếu khác màu thì có thể ăn
					if g.getPieceColor(nx, ny) != g.getPieceColor(x, y) {
						moves = append(moves, Point{X: nx, Y: ny})
					}
					break
				}
				nx += dir.X
				ny += dir.Y
			}
		}

		// TODO: các quân khác như Mã, Tượng, Sĩ, Tướng,...
	}

	return moves
}

// calculatePossibleMoves tính toán các nước đi khả dĩ cho quân cờ tại (pieceX, pieceY)
func (g *ChessBoardService) calculatePossibleMoves(pieceX, pieceY int) []Point {
	moves := []Point{}
	pieceName := g.board[pieceY][pieceX]
	pieceColor := g.getPieceColor(pieceX, pieceY)

	if pieceColor == 0 {
		return moves
	}

	// Tướng/Tướng (General/Marshal) - "將" / "帥"
	if pieceName == "將" || pieceName == "帥" {
		palaceMinX, palaceMaxX := 3, 5
		palaceMinY, palaceMaxY := 0, 2
		if pieceColor == 2 { // Nếu là Tướng đỏ, cung ở dưới
			palaceMinY, palaceMaxY = 7, 9
		}

		for dy := -1; dy <= 1; dy++ {
			for dx := -1; dx <= 1; dx++ {
				if (dx == 0 && dy == 0) || (dx != 0 && dy != 0) { // Chỉ đi ngang hoặc dọc 1 ô
					continue
				}

				newX, newY := pieceX+dx, pieceY+dy
				if g.isInsideBoard(newX, newY) &&
					newX >= palaceMinX && newX <= palaceMaxX &&
					newY >= palaceMinY && newY <= palaceMaxY {

					targetColor := g.getPieceColor(newX, newY)
					if targetColor == 0 || targetColor != pieceColor {
						moves = append(moves, Point{newX, newY})
					}
				}
			}
		}
	}

	// Sĩ (Advisor) - "士"
	if pieceName == "士" {
		palaceMinX, palaceMaxX := 3, 5
		palaceMinY, palaceMaxY := 0, 2
		if pieceColor == 2 { // Nếu là Sĩ đỏ, cung ở dưới
			palaceMinY, palaceMaxY = 7, 9
		}

		for dy := -1; dy <= 1; dy += 2 { // Đi chéo 1 ô
			for dx := -1; dx <= 1; dx += 2 {
				newX, newY := pieceX+dx, pieceY+dy
				if g.isInsideBoard(newX, newY) &&
					newX >= palaceMinX && newX <= palaceMaxX &&
					newY >= palaceMinY && newY <= palaceMaxY {

					targetColor := g.getPieceColor(newX, newY)
					if targetColor == 0 || targetColor != pieceColor {
						moves = append(moves, Point{newX, newY})
					}
				}
			}
		}
	}

	// Tượng/Tượng (Elephant/Minister) - "象" / "相"
	if pieceName == "象" || pieceName == "相" {
		deltas := []Point{{2, 2}, {2, -2}, {-2, 2}, {-2, -2}} // Đi chéo 2 ô
		for _, d := range deltas {
			newX, newY := pieceX+d.X, pieceY+d.Y
			midX, midY := pieceX+d.X/2, pieceY+d.Y/2 // Vị trí chân tượng

			if g.isInsideBoard(newX, newY) && g.board[midY][midX] == "" { // Không bị cản chân
				// Tượng không được vượt sông
				if (pieceColor == 1 && newY >= 5) || (pieceColor == 2 && newY < 5) {
					continue
				}

				targetColor := g.getPieceColor(newX, newY)
				if targetColor == 0 || targetColor != pieceColor {
					moves = append(moves, Point{newX, newY})
				}
			}
		}
	}

	// Mã (Horse/Knight) - "馬" / "傌"
	if pieceName == "馬" || pieceName == "傌" {
		horseMoves := []struct {
			dx, dy           int
			blockDx, blockDy int // Vị trí chân mã
		}{
			{1, 2, 0, 1}, {1, -2, 0, -1}, {-1, 2, 0, 1}, {-1, -2, 0, -1},
			{2, 1, 1, 0}, {2, -1, 1, 0}, {-2, 1, -1, 0}, {-2, -1, -1, 0},
		}

		for _, move := range horseMoves {
			newX, newY := pieceX+move.dx, pieceY+move.dy
			blockX, blockY := pieceX+move.blockDx, pieceY+move.blockDy

			if g.isInsideBoard(blockX, blockY) && g.board[blockY][blockX] != "" { // Bị cản chân
				continue
			}

			if g.isInsideBoard(newX, newY) {
				targetColor := g.getPieceColor(newX, newY)
				if targetColor == 0 || targetColor != pieceColor {
					moves = append(moves, Point{newX, newY})
				}
			}
		}
	}

	// Xe (Chariot/Rook) - "車" / "俥"
	if pieceName == "車" || pieceName == "俥" {
		directions := []Point{{0, 1}, {0, -1}, {1, 0}, {-1, 0}} // 4 hướng: xuống, lên, phải, trái

		for _, dir := range directions {
			for step := 1; ; step++ {
				newX, newY := pieceX+dir.X*step, pieceY+dir.Y*step
				if !g.isInsideBoard(newX, newY) {
					break // Ra ngoài bàn cờ
				}

				targetColor := g.getPieceColor(newX, newY)
				if targetColor == 0 {
					moves = append(moves, Point{newX, newY})
				} else if targetColor != pieceColor {
					moves = append(moves, Point{newX, newY}) // Ăn quân
					break
				} else {
					break // Gặp quân cùng màu
				}
			}
		}
	}

	// Pháo (Cannon) - "砲" / "炮"
	if pieceName == "砲" || pieceName == "炮" {
		directions := []Point{{0, 1}, {0, -1}, {1, 0}, {-1, 0}} // 4 hướng

		for _, dir := range directions {
			seenPivot := false // Biến để kiểm tra đã gặp "ngòi" (quân cản đầu tiên) chưa
			for step := 1; ; step++ {
				newX, newY := pieceX+dir.X*step, pieceY+dir.Y*step
				if !g.isInsideBoard(newX, newY) {
					break // Ra ngoài bàn cờ
				}

				targetColor := g.getPieceColor(newX, newY)
				if !seenPivot {
					if targetColor == 0 {
						moves = append(moves, Point{newX, newY}) // Đi như Xe nếu chưa gặp ngòi
					} else {
						seenPivot = true // Đã gặp ngòi
					}
				} else { // Đã gặp ngòi
					if targetColor != 0 {
						if targetColor != pieceColor {
							moves = append(moves, Point{newX, newY}) // Ăn quân khác màu sau ngòi
						}
						break // Gặp quân sau ngòi (dù cùng màu hay khác màu) thì dừng lại
					}
				}
			}
		}
	}

	// Tốt (Pawn/Soldier) - "卒" / "兵"
	if pieceName == "卒" || pieceName == "兵" {
		if pieceColor == 1 { // Tốt đen (卒)
			newY := pieceY + 1 // Luôn tiến thẳng
			if g.isInsideBoard(pieceX, newY) {
				targetColor := g.getPieceColor(pieceX, newY)
				if targetColor == 0 || targetColor != pieceColor {
					moves = append(moves, Point{pieceX, newY})
				}
			}
			if newY >= 5 { // Đã qua sông, có thể đi ngang
				if g.isInsideBoard(pieceX+1, pieceY) { // Sang phải
					targetColor := g.getPieceColor(pieceX+1, pieceY)
					if targetColor == 0 || targetColor != pieceColor {
						moves = append(moves, Point{pieceX + 1, pieceY})
					}
				}
				if g.isInsideBoard(pieceX-1, pieceY) { // Sang trái
					targetColor := g.getPieceColor(pieceX-1, pieceY)
					if targetColor == 0 || targetColor != pieceColor {
						moves = append(moves, Point{pieceX - 1, pieceY})
					}
				}
			}
		}

		if pieceColor == 2 { // Tốt đỏ (兵)
			newY := pieceY - 1 // Luôn tiến thẳng
			if g.isInsideBoard(pieceX, newY) {
				targetColor := g.getPieceColor(pieceX, newY)
				if targetColor == 0 || targetColor != pieceColor {
					moves = append(moves, Point{pieceX, newY})
				}
			}
			if newY < 5 { // Đã qua sông, có thể đi ngang
				if g.isInsideBoard(pieceX+1, pieceY) { // Sang phải
					targetColor := g.getPieceColor(pieceX+1, pieceY)
					if targetColor == 0 || targetColor != pieceColor {
						moves = append(moves, Point{pieceX + 1, pieceY})
					}
				}
				if g.isInsideBoard(pieceX-1, pieceY) { // Sang trái
					targetColor := g.getPieceColor(pieceX-1, pieceY)
					if targetColor == 0 || targetColor != pieceColor {
						moves = append(moves, Point{pieceX - 1, pieceY})
					}
				}
			}
		}
	}

	return moves
}
