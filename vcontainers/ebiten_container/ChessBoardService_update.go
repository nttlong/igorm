package ebiten_container

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// containsPoint kiểm tra xem một Point có trong slice Point không
func containsPoint(list []Point, p Point) bool {
	for _, item := range list {
		if item.X == p.X && item.Y == p.Y {
			return true
		}
	}
	return false
}
func (g *ChessBoardService) Update() error {
	// EbitenUI xử lý input trước, sau đó chúng ta xử lý input của game

	// Logic xử lý chuột của bạn cho game chỉ hoạt động trên vùng bàn cờ.
	// Bàn cờ được đặt ở giữa với kích thước 576x640.
	// TextArea trái rộng 200px, TextArea phải rộng 200px.
	// Tổng chiều rộng: 200 (TextArea trái) + 576 (bàn cờ) + 200 (TextArea phải) = 976px
	// Vị trí X của bàn cờ sẽ từ 200 đến 200 + 576 = 776
	cursorX, cursorY := ebiten.CursorPosition()
	boardMouseX := cursorX - g.Left // Trừ đi chiều rộng của TextArea bên trái để lấy tọa độ tương đối trên bàn cờ
	boardMouseY := cursorY          // Chiều Y không đổi

	// Chuyển đổi tọa độ pixel trên bàn cờ thành tọa độ ô cờ
	boardGridX, boardGridY := boardMouseX/64, boardMouseY/64

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		// Chỉ xử lý click nếu nằm trong vùng bàn cờ (x từ 200 đến 776)
		if cursorX >= g.Left && cursorX < (g.Left+576) && g.isInsideBoard(boardGridX, boardGridY) {
			// Trường hợp 1: Đã có quân cờ được chọn
			if g.selectedPiece != nil {
				selectedPiece := *g.selectedPiece // Tạo bản sao của Point đã chọn

				targetPoint := Point{boardGridX, boardGridY}
				if containsPoint(g.possibleMoves, targetPoint) {
					// Thực hiện di chuyển quân cờ
					g.board[boardGridY][boardGridX] = g.board[selectedPiece.Y][selectedPiece.X]
					g.board[selectedPiece.Y][selectedPiece.X] = ""
					g.selectedPiece = nil // Hủy chọn quân
					g.possibleMoves = nil // Xóa các nước đi khả dĩ

					// Cập nhật nội dung của cả hai TextAreas sau mỗi nước đi hợp lệ
					boardPrompt := g.getBoardStateAsPrompt()
					// g.promptTextAreaLeft.SetText("Black AI Prompt:\n\n" + boardPrompt)
					// g.promptTextAreaRight.SetText("Red AI Prompt:\n\n" + boardPrompt)

					log.Println("--- Board State for OpenAI (after move) ---")
					log.Println(boardPrompt)
					log.Println("-------------------------------------------")

				} else {
					// Nếu nhấp vào một ô không hợp lệ (không phải nước đi khả dĩ) hoặc nhấp lại vào chính quân cờ
					if selectedPiece.X == boardGridX && selectedPiece.Y == boardGridY {
						g.selectedPiece = nil // Hủy chọn quân
						g.possibleMoves = nil // Xóa các nước đi khả dĩ
					} else {
						// Nếu chọn một quân khác khi đã có quân được chọn nhưng nước đi không hợp lệ
						g.selectedPiece = nil
						g.possibleMoves = nil
						if g.board[boardGridY][boardGridX] != "" {
							g.selectedPiece = &Point{boardGridX, boardGridY}
							g.possibleMoves = g.calculatePossibleMoves(boardGridX, boardGridY)
						}
					}
				}
			} else {
				// Trường hợp 2: Chưa có quân cờ nào được chọn, và bạn nhấp vào một ô có quân
				if g.board[boardGridY][boardGridX] != "" {
					g.selectedPiece = &Point{boardGridX, boardGridY}
					g.possibleMoves = g.calculatePossibleMoves(boardGridX, boardGridY)
				}
			}
		} else {
			// Nếu click chuột ngoài bàn cờ, hủy chọn quân
			g.selectedPiece = nil
			g.possibleMoves = nil
		}
	}
	boardPrompt := g.getBoardStateAsPrompt()
	g.Console.Log(boardPrompt)
	return nil
}
