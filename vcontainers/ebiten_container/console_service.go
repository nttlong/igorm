package ebiten_container

import (
	"embed"
	"image"
	"image/color"
	"log"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

var fontBytes []byte

//go:embed Roboto-Regular.ttf
var ConsolefontData embed.FS

type ConsoleService struct {
	Message   string
	MaxLines  int // VD: 100 dòng tối đa lưu trữ
	FontFace  font.Face
	Width     int // chiều rộng console tính bằng pixel (để tính khi wrap)
	Height    int
	fontBytes []byte
}

func (c *ConsoleService) LoadFontFace(size float64) *ConsoleService {
	data, err := ConsolefontData.ReadFile("Roboto-Regular.ttf")
	if err != nil {
		log.Fatal("Error reading font file:", err)
	}
	c.fontBytes = data
	ft, err := opentype.Parse(c.fontBytes)
	if err != nil {
		log.Fatal(err)
	}

	face, err := opentype.NewFace(ft, &opentype.FaceOptions{
		Size:    size,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
	c.FontFace = face

	return c
}
func (c *ConsoleService) Log(text string) {
	c.Message = text
	// wrapped := wrapText(text, c.FontFace, c.Width-20) // 20px padding
	// c.Lines = append(c.Lines, wrapped...)

	// if len(c.Lines) > c.MaxLines {
	// 	c.Lines = c.Lines[len(c.Lines)-c.MaxLines:]
	// }
}
func wrapText1(s string, face font.Face, maxWidth int) []string {
	var lines []string
	words := strings.Fields(s)
	if len(words) == 0 {
		return []string{""}
	}

	current := words[0]
	for _, word := range words[1:] {
		test := current + " " + word
		width := text.BoundString(face, test).Dx()
		if width <= maxWidth {
			current = test
		} else {
			lines = append(lines, current)
			current = word
		}
	}
	lines = append(lines, current)
	return lines
}
func (c *ConsoleService) Draw(screen *ebiten.Image, x, y int) {
	if c.FontFace == nil {
		return // Không có font thì không vẽ gì
	}
	//ebitenutil.DrawRect(screen, float64(x), float64(y), 9*64, 64, color.RGBA{0, 128, 255, 100})
	// lineHeight := c.FontFace.Metrics().Height.Ceil() + 4
	// maxVisibleLines := 600 / lineHeight // Có thể cấu hình

	// start := 0
	// if len(c.Lines) > maxVisibleLines {
	// 	start = len(c.Lines) - maxVisibleLines
	// }
	// linesToDraw := c.Lines[start:]
	// height := len(linesToDraw) * lineHeight

	// Fill background màu trắng (hoặc sửa thành màu bạn muốn)
	bgColor := color.White
	c.Height = 600
	rect := image.Rect(x, y, x+c.Width, y+c.Height)
	screen.SubImage(rect).(*ebiten.Image).Fill(bgColor)

	// Vẽ từng dòng chữ
	textColor := color.White // Đổi lại trắng nếu nền là xám
	text.Draw(screen, c.Message, c.FontFace, x+10, 10, textColor)
	// cursorY := y + lineHeight
	// for _, line := range linesToDraw {
	// 	text.Draw(screen, line, c.FontFace, x+10, cursorY, textColor)
	// 	cursorY += lineHeight
	// }
}
