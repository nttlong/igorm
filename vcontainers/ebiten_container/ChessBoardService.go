package ebiten_container

import (
	"embed"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
)

//go:embed NotoSansCJK-Regular.ttf
var fontData embed.FS

// Point struct để dễ dàng làm việc với tọa độ (x, y)
type Point struct {
	X int
	Y int
}
type ChessBoardService struct {
	Left          int
	Width         int
	Height        int
	BoardWitdh    int
	BoardHeight   int
	board         [10][9]string
	selectedPiece *Point  // Thay đổi để dùng struct Point
	possibleMoves []Point // Lưu trữ các nước đi khả dĩ
	fontFace      font.Face
	Console       *ConsoleService
}

// Khởi tạo font chữ Hán
func (g *ChessBoardService) Draw1(screen *ebiten.Image) {
	g.Console.Log("Hello, Console!")
	op := &ebiten.DrawImageOptions{}
	boardImage := ebiten.NewImage(1200, 640)
	g.Console.Draw(boardImage, g.BoardWitdh, 0)
	screen.DrawImage(boardImage, op)
	ebitenutil.DrawRect(boardImage, float64(g.Left), 0, float64(g.BoardWitdh), float64(g.BoardHeight), color.RGBA{210, 180, 140, 255})

}
func (g *ChessBoardService) Draw(screen *ebiten.Image) {
	// Tạo một Image riêng cho bàn cờ và vẽ nó
	boardImage := ebiten.NewImage(1200, 640)

	ebitenutil.DrawRect(boardImage, float64(g.Left), 0, float64(g.BoardWitdh), float64(g.BoardHeight), color.RGBA{210, 180, 140, 255}) // Nền

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
	//op.GeoM.Translate(200, 0) // Di chuyển bàn cờ sang phải 200 pixel (chiều rộng của TextArea trái)

	boardPrompt := g.getBoardStateAsPrompt()
	g.Console.Log(boardPrompt)
	g.Console.Draw(boardImage, g.BoardWitdh, 0)
	screen.DrawImage(boardImage, op)
	// Vẽ UI của ebitenui lên màn hình (nó sẽ tự đặt các TextAreas vào đúng vị trí)

}

// isInsideBoard kiểm tra xem tọa độ (x, y) có nằm trong bàn cờ không
func (g *ChessBoardService) isInsideBoard(x, y int) bool {
	return x >= 0 && x < 9 && y >= 0 && y < 10
}

// getPieceColor xác định màu của quân cờ tại (x, y)
// 0: Không có quân, 1: Đen, 2: Đỏ
func (g *ChessBoardService) getPieceColor(x, y int) int {
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
func (g *ChessBoardService) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.Width, g.Height
}
func (g *ChessBoardService) Run() {
	ebiten.SetWindowSize(g.Width, g.Height)
	ebiten.SetWindowTitle("My Ebiten Game")

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
