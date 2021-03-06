package main

import (
    "log"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	screenWidth = 640
	screenHeight = 480
)

type Mode int

const (
	ModeTitle Mode = iota
	ModeGame
	ModeGameOver
)

type Game struct{
	mode Mode
}

func NewGame() *Game {
	g := &Game{}
	return g
}

func clickMouseButton() bool {
	return inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft)
}

func (g *Game) Update() error {
	switch g.mode{
	case ModeTitle:
		if clickMouseButton(){ g.mode = ModeGame }
	case ModeGame:
		if clickMouseButton(){ g.mode = ModeGameOver }
	case ModeGameOver:
		if clickMouseButton(){ g.mode = ModeTitle }	
	}

    return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	switch g.mode{
	case ModeTitle:
		ebitenutil.DebugPrint(screen, "ModeTitle")
	case ModeGame:
		ebitenutil.DebugPrint(screen, "ModeGame")
	case ModeGameOver:
		ebitenutil.DebugPrint(screen, "ModeGameOver")
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
    return 300, 300
}

func main() {
    ebiten.SetWindowSize(screenWidth, screenHeight)
    ebiten.SetWindowTitle("FlappyGorira")
    if err := ebiten.RunGame(NewGame()); err != nil {
        log.Fatal(err)
    }
}