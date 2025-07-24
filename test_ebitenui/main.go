package main

import (
	"image/color"
	"log"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/themes"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	ui *ebitenui.UI
}

func (g *Game) Update() error {
	g.ui.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.ui.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 800, 600 // Kích thước cửa sổ
}

func main() {
	// Thử khởi tạo Dark Theme
	theme := themes.NewDarkTheme()
	if theme == nil {
		log.Fatal("themes.NewDarkTheme() returned nil!")
	}
	log.Println("themes.NewDarkTheme() initialized successfully.")

	// Thử khởi tạo Color Resource
	colorRes := themes.NewColorResource(color.RGBA{255, 0, 0, 255})
	if colorRes == nil {
		log.Fatal("themes.NewColorResource() returned nil!")
	}
	log.Println("themes.NewColorResource() initialized successfully.")

	// Thử khởi tạo GridLayoutSpecFixed
	fixedSpec := widget.GridLayoutSpecFixed(100)
	log.Printf("widget.GridLayoutSpecFixed(100) created: %+v\n", fixedSpec)

	// Tạo một TextAreaTheme (chỉ là struct, không phải hàm khởi tạo)
	textAreaTheme := &widget.TextAreaTheme{
		Idle:  colorRes,
		Face:  theme.Face, // Sử dụng font từ theme
		Color: color.White,
	}
	log.Printf("widget.TextAreaTheme created: %+v\n", textAreaTheme)

	// Tạo UI cơ bản
	rootContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			widget.GridLayoutOpts.Columns(1),
			widget.GridLayoutOpts.ColumnSpecs(widget.GridLayoutSpecFixed(800)),
		)),
	)

	textInput := widget.NewTextInput(
		widget.TextInputOpts.ContainerOpts(widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.MinSize(100, 30),
		)),
		widget.TextInputOpts.Theme(theme),
		widget.TextInputOpts.MaxLen(50),
		widget.TextInputOpts.Placeholder("Type something..."),
	)
	rootContainer.AddChild(textInput)

	game := &Game{
		ui: &ebitenui.UI{
			Container: rootContainer,
		},
	}

	ebiten.SetWindowSize(800, 600)
	ebiten.SetWindowTitle("EbitenUI Test")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
