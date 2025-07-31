package ebiten_container

import (
	"log"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

func (g *ChessBoardService) initFont() {
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
