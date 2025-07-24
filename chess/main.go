package main

import (
	"embed"
	"fmt"
	"image/color"
	"log"
	"strings"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/themes" // Đảm bảo import này đúng
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	// Import này cần thiết cho embed.FS
)

//go:embed NotoSansCJK-Regular.ttf
var fontData embed.FS

// Point struct để dễ dàng làm việc với tọa độ (x, y)
type Point struct {
	X int
	Y int
}

type Game struct {
	board         [10][9]string
	selectedPiece *Point  // Thay đổi để dùng struct Point
	possibleMoves []Point // Lưu trữ các nước đi khả dĩ
	fontFace      font.Face

	// EbitenUI fields
	ui                  *ebitenui.UI
	promptTextAreaLeft  *widget.TextArea
	promptTextAreaRight *widget.TextArea
}

// Khởi tạo font chữ Hán
func (g *Game) initFont() {
	data, err := fontData.ReadFile("NotoSansCJK-Regular.ttf")
	if err != nil {
		log.Fatal("Error reading font file:", err)
	}

	tt, err := opentype.Parse(data)
	if err != nil {
		log.Fatal("Error parsing font:", err)
	}

	g.fontFace, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    28, // Kích thước chữ vừa với quân cờ
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal("Error creating font face:", err)
	}
}

func (g *Game) initBoard() {
	g.board = [10][9]string{
		{"車", "馬", "象", "士", "將", "士", "象", "馬", "車"}, // Hàng 0: Quân Đen
		{"", "", "", "", "", "", "", "", ""},
		{"", "砲", "", "", "", "", "", "砲", ""},    // Pháo Đen
		{"卒", "", "卒", "", "卒", "", "卒", "", "卒"}, // Tốt Đen
		{"", "", "", "", "", "", "", "", ""},
		{"", "", "", "", "", "", "", "", ""},      // Sông
		{"兵", "", "兵", "", "兵", "", "兵", "", "兵"}, // Tốt Đỏ
		{"", "炮", "", "", "", "", "", "炮", ""},    // Pháo Đỏ
		{"", "", "", "", "", "", "", "", ""},
		{"俥", "傌", "相", "士", "帥", "士", "相", "傌", "俥"}, // Hàng 9: Quân Đỏ
	}
}

// isInsideBoard kiểm tra xem tọa độ (x, y) có nằm trong bàn cờ không
func (g *Game) isInsideBoard(x, y int) bool {
	return x >= 0 && x < 9 && y >= 0 && y < 10
}

// getPieceColor xác định màu của quân cờ tại (x, y)
// 0: Không có quân, 1: Đen, 2: Đỏ
func (g *Game) getPieceColor(x, y int) int {
	if !g.isInsideBoard(x, y) || g.board[y][x] == "" {
		return 0 // Không có quân hoặc ngoài bàn cờ
	}

	piece := g.board[y][x]
	// Các quân đen thường dùng ký tự đơn giản, quân đỏ dùng ký tự phức tạp hơn hoặc khác biệt (vd: 帥 vs 將)
	switch piece {
	case "車", "馬", "象", "將", "砲", "卒":
		return 1 // Quân đen
	case "士": // Sĩ là quân đặc biệt, cần xem vị trí
		if y <= 4 { // Quân sĩ đen ở hàng trên của quân đen
			return 1
		}
		return 2 // Quân sĩ đỏ ở hàng dưới của quân đỏ
	case "俥", "傌", "相", "帥", "炮", "兵":
		return 2 // Quân đỏ
	}
	return 0 // Mặc định không xác định
}

// calculatePossibleMoves tính toán các nước đi khả dĩ cho quân cờ tại (pieceX, pieceY)
func (g *Game) calculatePossibleMoves(pieceX, pieceY int) []Point {
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

// containsPoint kiểm tra xem một Point có trong slice Point không
func containsPoint(list []Point, p Point) bool {
	for _, item := range list {
		if item.X == p.X && item.Y == p.Y {
			return true
		}
	}
	return false
}

// getBoardStateAsPrompt tạo một chuỗi mô tả trạng thái bàn cờ cho AI
func (g *Game) getBoardStateAsPrompt() string {
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

func (g *Game) Update() error {
	// EbitenUI xử lý input trước, sau đó chúng ta xử lý input của game
	g.ui.Update()

	// Logic xử lý chuột của bạn cho game chỉ hoạt động trên vùng bàn cờ.
	// Bàn cờ được đặt ở giữa với kích thước 576x640.
	// TextArea trái rộng 200px, TextArea phải rộng 200px.
	// Tổng chiều rộng: 200 (TextArea trái) + 576 (bàn cờ) + 200 (TextArea phải) = 976px
	// Vị trí X của bàn cờ sẽ từ 200 đến 200 + 576 = 776
	cursorX, cursorY := ebiten.CursorPosition()
	boardMouseX := cursorX - 200 // Trừ đi chiều rộng của TextArea bên trái để lấy tọa độ tương đối trên bàn cờ
	boardMouseY := cursorY       // Chiều Y không đổi

	// Chuyển đổi tọa độ pixel trên bàn cờ thành tọa độ ô cờ
	boardGridX, boardGridY := boardMouseX/64, boardMouseY/64

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		// Chỉ xử lý click nếu nằm trong vùng bàn cờ (x từ 200 đến 776)
		if cursorX >= 200 && cursorX < (200+576) && g.isInsideBoard(boardGridX, boardGridY) {
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
					g.promptTextAreaLeft.SetText("Black AI Prompt:\n\n" + boardPrompt)
					g.promptTextAreaRight.SetText("Red AI Prompt:\n\n" + boardPrompt)

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
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Tạo một Image riêng cho bàn cờ và vẽ nó
	boardImage := ebiten.NewImage(576, 640)
	ebitenutil.DrawRect(boardImage, 0, 0, 576, 640, color.RGBA{210, 180, 140, 255}) // Nền

	// Vẽ lưới bàn cờ
	for i := 0; i < 10; i++ {
		for j := 0; j < 9; j++ {
			ebitenutil.DrawRect(boardImage, float64(j*64+28), float64(i*64+28), 8, 8, color.RGBA{50, 50, 50, 255})
		}
	}

	// Vẽ các đường ngang và dọc
	for i := 0; i < 10; i++ {
		ebitenutil.DrawLine(boardImage, 0, float64(i*64+32), 8*64, float64(i*64+32), color.Black)
	}
	for j := 0; j < 9; j++ {
		if j == 0 || j == 8 {
			ebitenutil.DrawLine(boardImage, float64(j*64+32), 0, float64(j*64+32), 9*64+32, color.Black)
		} else {
			ebitenutil.DrawLine(boardImage, float64(j*64+32), 0, float64(j*64+32), 4*64+32, color.Black)
			ebitenutil.DrawLine(boardImage, float64(j*64+32), 5*64+32, float64(j*64+32), 9*64+32, color.Black)
		}
	}

	// Vẽ sông (màu xanh lam)
	ebitenutil.DrawRect(boardImage, 0, 5*64, 9*64, 64, color.RGBA{0, 128, 255, 100})
	text.Draw(boardImage, "楚河漢界", g.fontFace, 4*64-30, 5*64+40, color.Black)

	// Vẽ cung Tướng (đường chéo màu vàng)
	for _, y := range []int{0, 7} {
		ebitenutil.DrawLine(boardImage, 3*64+32, float64(y*64+32), 5*64+32, float64((y+2)*64+32), color.RGBA{255, 215, 0, 255})
		ebitenutil.DrawLine(boardImage, 5*64+32, float64(y*64+32), 3*64+32, float64((y+2)*64+32), color.RGBA{255, 215, 0, 255})
	}

	// Vẽ quân cờ
	for i := 0; i < 10; i++ {
		for j := 0; j < 9; j++ {
			if g.board[i][j] != "" {
				pieceColor := color.RGBA{255, 0, 0, 255}    // Mặc định là đỏ
				textColor := color.RGBA{255, 255, 255, 255} // Mặc định là chữ trắng

				pieceSide := g.getPieceColor(j, i)

				if pieceSide == 1 { // Quân đen
					pieceColor = color.RGBA{0, 0, 0, 255}
					textColor = color.RGBA{255, 255, 255, 255}
				} else if pieceSide == 2 { // Quân đỏ
					// Các quân đỏ thông thường có chữ đen, riêng pháo (炮) có chữ trắng trên nền đỏ
					if g.board[i][j] == "炮" {
						textColor = color.RGBA{255, 255, 255, 255}
					} else {
						textColor = color.RGBA{0, 0, 0, 255}
					}
				}

				ebitenutil.DrawCircle(boardImage, float64(j*64+32), float64(i*64+32), 28, pieceColor)
				text.Draw(boardImage, g.board[i][j], g.fontFace, j*64+20, i*64+40, textColor)
			}
		}
	}

	// Đánh dấu quân cờ được chọn
	if g.selectedPiece != nil {
		ebitenutil.DrawRect(boardImage, float64(g.selectedPiece.X*64), float64(g.selectedPiece.Y*64), 64, 64, color.RGBA{255, 255, 0, 100})

		// Vẽ các nước đi khả dĩ
		for _, move := range g.possibleMoves {
			ebitenutil.DrawCircle(boardImage, float64(move.X*64+32), float64(move.Y*64+32), 10, color.RGBA{0, 255, 0, 150})
		}
	}

	// Vẽ boardImage lên màn hình chính tại vị trí của bàn cờ (sau TextArea trái)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(200, 0) // Di chuyển bàn cờ sang phải 200 pixel (chiều rộng của TextArea trái)
	screen.DrawImage(boardImage, op)

	// Vẽ UI của ebitenui lên màn hình (nó sẽ tự đặt các TextAreas vào đúng vị trí)
	g.ui.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	// Tổng chiều rộng: 200 (left) + 576 (board) + 200 (right) = 976
	return 200 + 576 + 200, 640
}

func main() {
	game := &Game{}
	game.initBoard()
	game.initFont()

	// --- Khởi tạo EbitenUI ---
	data, err := fontData.ReadFile("NotoSansCJK-Regular.ttf")
	if err != nil {
		log.Fatal("Error reading font file for theme:", err)
	}

	// SỬA: opentype.Parse trả về *ttf.Font. NewFace giờ nhận *ttf.Font
	ttFont, err := opentype.Parse(data)
	if err != nil {
		log.Fatal("Error parsing font for theme:", err)
	}

	faceSource, err := opentype.NewFace(ttFont, &opentype.FaceOptions{
		Size:    16, // Kích thước font cho UI (nhỏ hơn font cờ để hiển thị được nhiều text)
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal("Error creating font face for theme:", err)
	}

	// SỬA LỖI: customTheme := themes.NewTheme()
	// themes.NewTheme() KHÔNG CÓ TRONG ebitenui MỚI.
	// BẠN CẦN SỬ DỤNG theme CÓ SẴN HOẶC TẠO MỘT CÁCH KHÁC.
	// Tôi sẽ dùng themes.NewDarkTheme() để đơn giản hóa.
	customTheme := themes.NewDarkTheme() // SỬA: Dùng theme có sẵn (NewDarkTheme)

	customTheme.Face = faceSource // Áp dụng fontface của bạn

	// SỬA LỖI: customTheme.TextArea
	// API của TextAreaTheme đã thay đổi. idle và disabled không phải là ButtonGraphic
	// Mà là ImageRenderer (Resource).
	// Và NewImageGraphic KHÔNG TỒN TẠI. Thay vào đó dùng themes.NewColorResource.
	customTheme.TextArea = &widget.TextAreaTheme{
		Idle:     themes.NewColorResource(color.RGBA{100, 100, 100, 200}), // SỬA: Dùng NewColorResource
		Disabled: themes.NewColorResource(color.RGBA{80, 80, 80, 150}),    // SỬA: Dùng NewColorResource
		Face:     faceSource,                                              // Font cho nội dung TextArea
		Color:    color.White,                                             // Màu chữ
		// CaretColor:    color.White, // Màu con trỏ
		// SelectionColor: color.RGBA{0, 100, 200, 100}, // Màu highlight khi chọn text
	}

	// Tạo một Container làm root (container gốc) để chứa tất cả các widget khác
	rootContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			widget.GridLayoutOpts.Columns(3), // 3 cột: Left TextArea, Board, Right TextArea
			// SỬA LỖI: ColumnSpecs và GridLayoutSpec
			widget.GridLayoutOpts.ColumnSpecs(
				// Trong ebitenui v1.4.0, GridLayoutSpecFixed và GridLayoutSpecStretch nằm trực tiếp trong gói widget
				widget.GridLayoutSpecFixed(200), // Left TextArea (200px)
				widget.GridLayoutSpecFixed(576), // Board (576px)
				widget.GridLayoutSpecFixed(200), // Right TextArea (200px)
			),
		)),
		// Bỏ WidgetOpts.MouseButtonAny() như đã thảo luận trước đó
	)

	// TextArea cho quân đen (bên trái)
	game.promptTextAreaLeft = widget.NewTextArea(
		widget.TextAreaOpts.ContainerOpts(widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.MinSize(200, 640))),
		widget.TextAreaOpts.Theme(customTheme), // SỬA: .Theme() là một tùy chọn hợp lệ trong v1.4.0
		widget.TextAreaOpts.Text("Black AI Prompt:\n\nTrạng thái bàn cờ sẽ xuất hiện ở đây sau mỗi nước đi."),
	)
	rootContainer.AddChild(game.promptTextAreaLeft)

	// Panel làm placeholder cho bàn cờ. Bàn cờ sẽ được vẽ riêng và sau đó đưa vào vị trí này.
	boardPlaceholder := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.MinSize(576, 640)), // Kích thước placeholder cho bàn cờ
	)
	rootContainer.AddChild(boardPlaceholder)

	// TextArea cho quân đỏ (bên phải)
	game.promptTextAreaRight = widget.NewTextArea(
		widget.TextAreaOpts.ContainerOpts(widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.MinSize(200, 640))),
		widget.TextAreaOpts.Theme(customTheme), // SỬA: .Theme() là một tùy chọn hợp lệ trong v1.4.0
		widget.TextAreaOpts.Text("Red AI Prompt:\n\nTrạng thái bàn cờ sẽ xuất hiện ở đây sau mỗi nước đi."),
	)
	rootContainer.AddChild(game.promptTextAreaRight)

	// Khởi tạo ebitenui UI với root container
	game.ui = &ebitenui.UI{
		Container: rootContainer,
	}

	// Thay đổi kích thước cửa sổ game để phù hợp với 3 cột
	ebiten.SetWindowSize(976, 640) // 200 (left) + 576 (board) + 200 (right) = 976
	ebiten.SetWindowTitle("Xiangqi - Chinese Chess")

	// Lần đầu khởi tạo, cập nhật trạng thái bàn cờ vào TextAreas
	initialBoardPrompt := game.getBoardStateAsPrompt()
	game.promptTextAreaLeft.SetText("Black AI Prompt:\n\n" + initialBoardPrompt)
	game.promptTextAreaRight.SetText("Red AI Prompt:\n\n" + initialBoardPrompt)
	log.Println("--- Initial Board State for OpenAI ---")
	log.Println(initialBoardPrompt)
	log.Println("------------------------------------")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
