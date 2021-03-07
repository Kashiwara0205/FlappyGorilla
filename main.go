package main

import (
	"image"
	_ "image/png"
	"image/color"
	"log"
	"os"
	"math"

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
	arcadeFont     font.Face
	gorillaImage   *ebiten.Image
)

func init(){
	file, _ := os.Open("image/gorilla.png")
	img, _, err := image.Decode(file)
	if err != nil {
		log.Fatal(err)
	}
	gorillaImage = ebiten.NewImageFromImage(img)
}

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

	gorilla_x int
	gorilla_y int
	gorilla_vy int

	cameraX int
	cameraY int


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
		texts = []string{"FLAPPY GORILLA", "", "", "", "CLICK MOUSE BUTTON"}
		drawText(screen, texts)
	case ModeGame:
		g.drawGorilla(screen)
		ebitenutil.DebugPrint(screen, "ModeGame")
	case ModeGameOver:
		g.drawGorilla(screen)
		ebitenutil.DebugPrint(screen, "ModeGameOver")
		texts = []string{"", "", "", "GAME OVER"}
		drawText(screen, texts)

	}
 
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
    return screenWidth, screenHeight
}


func (g *Game) drawGorilla(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	w, h := gorillaImage.Size()
	op.GeoM.Translate(-float64(w)/2.0, -float64(h)/2.0)
	op.GeoM.Rotate(float64(g.gorilla_vy) / 96.0 * math.Pi / 6)
	op.GeoM.Translate(float64(w)/2.0, float64(h)/2.0)
	op.GeoM.Translate(float64(g.gorilla_x/16.0)-float64(g.cameraX), float64(g.gorilla_y/16.0)-float64(g.cameraY))
	op.Filter = ebiten.FilterLinear
	screen.DrawImage(gorillaImage, op)
}

func main() {
    ebiten.SetWindowSize(screenWidth, screenHeight)
    ebiten.SetWindowTitle("FlappyGORILLA")
    if err := ebiten.RunGame(NewGame()); err != nil {
        log.Fatal(err)
    }
}