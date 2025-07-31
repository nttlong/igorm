package ebiten_container

import (
	"vdi"
)

type EbitenContainer struct {
	Game    vdi.Singleton[EbitenContainer, *ChessBoardService]
	Console vdi.Singleton[EbitenContainer, *ConsoleService]
}

func NewEbitenContainer() *EbitenContainer {

	ret := vdi.NewContainer(func(owner *EbitenContainer) error {
		/*
			576, 640
		*/
		owner.Game.Init = func(owner *EbitenContainer) *ChessBoardService {
			ret := &ChessBoardService{
				Width:       1200,
				Height:      640,
				BoardWitdh:  576,
				BoardHeight: 640,
				Console:     owner.Console.Get(),
			}
			ret.initFont()
			ret.initBoard()
			return ret
		}
		owner.Console.Init = func(owner *EbitenContainer) *ConsoleService {
			ret := &ConsoleService{
				MaxLines: 1000,
			}
			ret.LoadFontFace(12)
			return ret
		}
		return nil

	})
	return ret
}
func (c *EbitenContainer) Run() {
	c.Game.Get().Run()
}
