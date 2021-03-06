package main

import (
    "log"

    "github.com/hajimehoshi/ebiten"
    "github.com/hajimehoshi/ebiten/ebitenutil"
)

type Game struct{}

func (g *Game) Update(screen *ebiten.Image) error {
    return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
    ebitenutil.DebugPrint(screen, "Ebiten")
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
    return 120, 140
}

func main() {
    ebiten.SetWindowSize(500, 500)
    ebiten.SetWindowTitle("Ebiten")
    if err := ebiten.RunGame(&Game{}); err != nil {
        log.Fatal(err)
    }
}