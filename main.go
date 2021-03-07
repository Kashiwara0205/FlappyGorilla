package main

import (
	"image/color"
	"log"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
)

const (
	screenWidth = 640
	screenHeight = 480
	fontSize = 32
)

var (
	arcadeFont	font.Face
)

func init() {
	tt, err := opentype.Parse(fonts.PressStart2P_ttf)
	if err != nil {
		log.Fatal(err)
	}
	const dpi = 72
	arcadeFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    fontSize,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
}

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

func drawText(screen *ebiten.Image, texts []string){
	for i, l := range texts {
		x := (screenWidth - len(l)*fontSize) / 2
		text.Draw(screen, l, arcadeFont, x, (i+4)*fontSize, color.White)
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0x80, 0xa0, 0xc0, 0xff})
	
	var texts []string

	switch g.mode{
	case ModeTitle:
		ebitenutil.DebugPrint(screen, "ModeTitle")
		texts = []string{"FLAPPY GORIRA", "", "", "", "CLICK MOUSE BUTTON"}
		drawText(screen, texts)
	case ModeGame:
		ebitenutil.DebugPrint(screen, "ModeGame")
	case ModeGameOver:
		ebitenutil.DebugPrint(screen, "ModeGameOver")
		texts = []string{"", "", "", "GAME OVER"}
		drawText(screen, texts)

	}
 
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
    return screenWidth, screenHeight
}

func main() {
    ebiten.SetWindowSize(screenWidth, screenHeight)
    ebiten.SetWindowTitle("FlappyGorira")
    if err := ebiten.RunGame(NewGame()); err != nil {
        log.Fatal(err)
    }
}