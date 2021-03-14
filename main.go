package main

import (
	"image"
	_ "image/png"
	"image/color"
	"log"
	"os"
	"math"
	"fmt"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
)

const (
	screenWidth      = 640
	screenHeight     = 480
	fontSize         = 32
	tileSize         = 32
	pipeWidth        = tileSize * 2
	pipeStartOffsetX = 8
	pipeIntervalX    = 8
	pipeGapY         = 5
)

var (
	arcadeFont     font.Face
	gorillaImage   *ebiten.Image
	tilesImage     *ebiten.Image
)

func init(){
	file, _ := os.Open("image/gorilla.png")
	img, _, err := image.Decode(file)
	if err != nil {
		log.Fatal(err)
	}
	gorillaImage = ebiten.NewImageFromImage(img)

	file, _ = os.Open("image/tiles.png")
	img, _, err = image.Decode(file)
	if err != nil {
		log.Fatal(err)
	}
	tilesImage = ebiten.NewImageFromImage(img)
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

func floorDiv(x, y int) int {
	d := x / y
	if d*y == x || x >= 0 {
		return d
	}
	return d - 1
}

func floorMod(x, y int) int {
	return x - floorDiv(x, y)*y
}
type Mode int

const (
	ModeTitle Mode = iota
	ModeGame
	ModeGameOver
)

type Game struct{
	mode Mode

	gorillaX int
	gorillaY int
	gorillaVy int

	cameraX int
	cameraY int

	pipeTileYs []int

	updateCount int

	ga *GA
}

const POPULATION = 1
const NUMBER_GENES = 100

type CpuPlayer struct {
	gene []int
	score int
	death bool
	idx int
}

func (player *CpuPlayer) shouldJump() bool {
	return true
}

func (player *CpuPlayer) nextStep() {
	if !player.death{
		player.idx++

		if 100 == player.idx{
			player.idx = 0
		}
	}
}

type GA struct{
	cpuPlayers [] CpuPlayer
	population int
}

func getRotateValue(values []int, i int) int{
	length := len(values)
	x := (length + i) / length
	idx := length + i - length * x 

	return values[idx]
}

func NewGame() *Game {
	g := &Game{}
	g.init()
	return g
}

func (g *Game) init() {
	g.gorillaX = 0
	g.gorillaY = 100 * 16
	g.cameraX = -240
	g.cameraY = 0
	g.pipeTileYs = make([]int, 256)

	// 土管の位置
	values := []int{2, 3, 4, 3, 5, 7, 2, 3, 4, 5}
	for i := range g.pipeTileYs {
		g.pipeTileYs[i] = getRotateValue(values, i)
	}

	// 遺伝子の初期化
	g.ga = NewGA()

	// 描画回数を記録する(評価タイミングに使用)
	g.updateCount = 0
}

func NewGA() *GA{
	ga := &GA{}
	ga.init()

	return ga
}

func (g *GA) init() {
	cnt := 0
	cpuPlayers := [] CpuPlayer{}

	gene := [] int{1}
	for cnt < POPULATION {
		player := CpuPlayer{ gene: gene, score: 0, death: false, idx: 0 }
		cpuPlayers = append(cpuPlayers, player)

		cnt++
	}

	g.population = POPULATION
	g.cpuPlayers = cpuPlayers
}

func clickMouseButton() bool {
	return inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft)
}

func (g *Game) Update() error {
	switch g.mode{
	case ModeTitle:
		if clickMouseButton(){ g.mode = ModeGame }
	case ModeGame:
		g.gorillaX += 32
		g.cameraX += 2

		// 40回目のUpdateでAIが行動する
		g.updateCount += 1
		if 40 == g.updateCount {
			g.updateCount = 0

			for _, player := range g.ga.cpuPlayers{
				if player.shouldJump() {
					g.gorillaVy = -96
				}

				player.nextStep()
			}
		}

		g.gorillaY += g.gorillaVy

		g.gorillaVy += 4
		if g.gorillaVy > 96 {
			g.gorillaVy = 96
		}

		if g.hit(){
			g.mode = ModeGameOver
		}

	case ModeGameOver:
		if clickMouseButton(){ 
			g.init()
			g.mode = ModeTitle 
		}
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
	g.drawTiles(screen)

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

	scoreStr := fmt.Sprintf("%04d", g.score())
	text.Draw(screen, scoreStr, arcadeFont, screenWidth-len(scoreStr)*fontSize, fontSize, color.White)
}

func (g *Game) pipeAt(tileX int) (tileY int, ok bool) {
	if (tileX - pipeStartOffsetX) <= 0 {
		return 0, false
	}
	if floorMod(tileX-pipeStartOffsetX, pipeIntervalX) != 0 {
		return 0, false
	}
	idx := floorDiv(tileX-pipeStartOffsetX, pipeIntervalX)
	return g.pipeTileYs[idx%len(g.pipeTileYs)], true
}

func (g *Game) score() int {
	x := floorDiv(g.gorillaX, 16) / tileSize
	if (x - pipeStartOffsetX) <= 0 {
		return 0
	}
	return floorDiv(x-pipeStartOffsetX, pipeIntervalX)
}

func (g *Game) hit() bool{	
	const (
		gorillaWidth  = 30
		gorillaHeight = 65
	)
	
	w, h := gorillaImage.Size()

	y0 := floorDiv(g.gorillaY, 16) + (h - gorillaHeight) / 2
	y1 := y0 + gorillaHeight

	if y0 < -tileSize * 3{
		return true
	}

	if y1 >= screenHeight-tileSize {
		return true
	}

	x0 := floorDiv(g.gorillaX, 16) + (w-gorillaWidth)/2
	x1 := x0 + gorillaWidth

	xMin := floorDiv(x0-pipeWidth, tileSize)
	xMax := floorDiv(x0+gorillaWidth, tileSize)
	for x := xMin; x <= xMax; x++ {
		y, ok := g.pipeAt(x)
		if !ok {
			continue
		}
		if x0 >= x*tileSize+pipeWidth {
			continue
		}
		if x1 < x*tileSize {
			continue
		}
		if y0 < y*tileSize {
			return true
		}
		if y1 >= (y+pipeGapY)*tileSize {
			return true
		}
	}
	
	return false
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
    return screenWidth, screenHeight
}

func (g *Game) drawGorilla(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	w, h := gorillaImage.Size()
	op.GeoM.Translate(-float64(w)/2.0, -float64(h)/2.0)
	op.GeoM.Rotate(float64(g.gorillaVy) / 96.0 * math.Pi / 6)
	op.GeoM.Translate(float64(w)/2.0, float64(h)/2.0)
	op.GeoM.Translate(float64(g.gorillaX/16.0)-float64(g.cameraX), float64(g.gorillaY/16.0)-float64(g.cameraY))
	op.Filter = ebiten.FilterLinear
	screen.DrawImage(gorillaImage, op)
}

func (g *Game) drawTiles(screen *ebiten.Image) {
	const (
		nx           = screenWidth / tileSize
		ny           = screenHeight / tileSize
		pipeTileSrcX = 128
		pipeTileSrcY = 192
	)

	op := &ebiten.DrawImageOptions{}
	for i := -2; i < nx+1; i++ {
		// ground
		op.GeoM.Reset()
		op.GeoM.Translate(float64(i*tileSize-floorMod(g.cameraX, tileSize)),
			float64((ny-1)*tileSize-floorMod(g.cameraY, tileSize)))
		screen.DrawImage(tilesImage.SubImage(image.Rect(0, 0, tileSize, tileSize)).(*ebiten.Image), op)

		// pipe
		if tileY, ok := g.pipeAt(floorDiv(g.cameraX, tileSize) + i); ok {
			for j := 0; j < tileY; j++ {
				op.GeoM.Reset()
				op.GeoM.Scale(1, -1)
				op.GeoM.Translate(float64(i*tileSize-floorMod(g.cameraX, tileSize)),
					float64(j*tileSize-floorMod(g.cameraY, tileSize)))
				op.GeoM.Translate(0, tileSize)
				var r image.Rectangle
				if j == tileY-1 {
					r = image.Rect(pipeTileSrcX, pipeTileSrcY, pipeTileSrcX+tileSize*2, pipeTileSrcY+tileSize)
				} else {
					r = image.Rect(pipeTileSrcX, pipeTileSrcY+tileSize, pipeTileSrcX+tileSize*2, pipeTileSrcY+tileSize*2)
				}
				screen.DrawImage(tilesImage.SubImage(r).(*ebiten.Image), op)
			}
			for j := tileY + pipeGapY; j < screenHeight/tileSize-1; j++ {
				op.GeoM.Reset()
				op.GeoM.Translate(float64(i*tileSize-floorMod(g.cameraX, tileSize)),
					float64(j*tileSize-floorMod(g.cameraY, tileSize)))
				var r image.Rectangle
				if j == tileY+pipeGapY {
					r = image.Rect(pipeTileSrcX, pipeTileSrcY, pipeTileSrcX+pipeWidth, pipeTileSrcY+tileSize)
				} else {
					r = image.Rect(pipeTileSrcX, pipeTileSrcY+tileSize, pipeTileSrcX+pipeWidth, pipeTileSrcY+tileSize+tileSize)
				}
				screen.DrawImage(tilesImage.SubImage(r).(*ebiten.Image), op)
			}
		}
	}
}

func main() {
    ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("FlappyGORILLA")
    if err := ebiten.RunGame(NewGame()); err != nil {
        log.Fatal(err)
    }
}