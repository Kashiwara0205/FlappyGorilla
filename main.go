package main

import (
	"flappyGorilla/utils"
	"flappyGorilla/ga"

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
	gorillaImages   [] *ebiten.Image
	tilesImage     *ebiten.Image
)

func init(){
	file, _ := os.Open("image/gorilla.png")
	img, _, err := image.Decode(file)
	if err != nil {
		log.Fatal(err)
	}

	gorillaImages = make([]*ebiten.Image, ga.POPULATION_SIZE)
	for i:= 0; i < ga.POPULATION_SIZE; i++{
		gorillaImages[i] = ebiten.NewImageFromImage(img)
	}

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

type Mode int

const (
	ModeTitle Mode = iota
	ModeGame
	ModeGameOver
)

type Game struct{
	mode Mode

	gorillaX []int
	gorillaY []int
	gorillaVy []int

	cameraX int
	cameraY int

	pipeTileYs []int

	updateCount int

	GA *ga.GA
}

func NewGame() *Game {
	g := &Game{}
	g.init()
	return g
}

func (g *Game) init() {
	g.gorillaX = make([]int, ga.POPULATION_SIZE)
	g.gorillaY = make([]int, ga.POPULATION_SIZE)
	g.gorillaVy = make([]int, ga.POPULATION_SIZE)

	for i:= 0; i < ga.POPULATION_SIZE; i++{
		g.gorillaX[i] = 0
		g.gorillaY[i] = 1600
	}

	g.cameraX = -240
	g.cameraY = 0
	g.pipeTileYs = make([]int, 256)

	// 土管の位置
	values := []int{2, 3, 4, 3, 5, 7, 2, 3, 4, 5}
	for i := range g.pipeTileYs {
		g.pipeTileYs[i] = utils.GetRotateValue(values, i)
	}

	// 遺伝子の初期化
	if g.GA == nil{
		g.GA = ga.NewGA()
	}else{
		fmt.Printf("=> UPDATE GA\n")
	}

	// 描画回数を記録する(評価タイミングに使用)
	g.updateCount = 0
}

func clickMouseButton() bool {
	return inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft)
}

func (g *Game) Update() error {
	switch g.mode{
	case ModeTitle:
		g.mode = ModeGame
	case ModeGame:
		g.cameraX += 2

		for i:= 0; i < ga.POPULATION_SIZE; i++{
			g.gorillaX[i] += 32
		}

		g.updateCount += 1
		if ga.ACTION_SPAN == g.updateCount {
			g.updateCount = 0

			for i, individual := range g.GA.Individuals{
				if individual.ShouldJump() {
					g.gorillaVy[i] = -96
				}

				individual.NextStep()
			}
		}

		// 各、ゴリラの重力計算
		for i:= 0; i < ga.POPULATION_SIZE; i++{
			g.gorillaY[i] += g.gorillaVy[i]

			g.gorillaVy[i] += 4
			if g.gorillaVy[i] > 96 {
				g.gorillaVy[i] = 96
			}
		}

		// 各、ゴリラの当たり判定
		for i, individual := range g.GA.Individuals{
			if !individual.CheckDead() && g.hit(i) {
				individual.Dead()
			}
		}

		// 全部のゴリラが死んだか判定
		if g.GA.CheckAllDead(){
			// 遺伝子の更新
			g.GA.Update()
			g.mode = ModeGameOver 
			
		}

	case ModeGameOver:
		g.init()
		g.mode = ModeGame
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
		ebitenutil.DebugPrint(screen, "ModeGameOver")
		texts = []string{"", "", "", "GAME OVER"}
		drawText(screen, texts)
	}
}

func (g *Game) pipeAt(tileX int) (tileY int, ok bool) {
	if (tileX - pipeStartOffsetX) <= 0 {
		return 0, false
	}

	if utils.FloorMod(tileX-pipeStartOffsetX, pipeIntervalX) != 0 {
		return 0, false
	}

	idx := utils.FloorDiv(tileX-pipeStartOffsetX, pipeIntervalX)
	return g.pipeTileYs[idx%len(g.pipeTileYs)], true
}

func (g *Game) hit(idx int) bool{	
	const (
		gorillaWidth  = 30
		gorillaHeight = 65
	)
	
	// gorilla Image size
	w, h := 75, 75
	y0 := utils.FloorDiv(g.gorillaY[idx], 16) + (h - gorillaHeight) / 2
	y1 := y0 + gorillaHeight

	if y0 >= tileSize * 2{
		g.GA.Individuals[idx].Score += 1
	}

	if y0 < -tileSize * 3{
		g.GA.Individuals[idx].Score -= 300
		return true
	}

	if y1 < screenHeight - tileSize * 2 {
		g.GA.Individuals[idx].Score += 1
	}

	if y1 >= screenHeight-tileSize {
		g.GA.Individuals[idx].Score -= 300
		return true
	}

	x0 := utils.FloorDiv(g.gorillaX[idx], 16) + (w - gorillaWidth) / 2
	x1 := x0 + gorillaWidth

	xMin := utils.FloorDiv(x0-pipeWidth, tileSize)
	xMax := utils.FloorDiv(x0+gorillaWidth, tileSize)

	passedPoint := 0

	for x := xMin; x <= xMax; x++ {
		y, ok := g.pipeAt(x)

		if !ok {
			continue	
		}

		topBoaderLine := y * tileSize
		underBoaderLine := (y+pipeGapY)*tileSize

		if x0 >= x*tileSize+pipeWidth {
			continue
		}
		if x1 < x*tileSize {	
			continue	
		}

		if y0 < topBoaderLine {
			diff := topBoaderLine - y0
			if 5 > diff{
				g.GA.Individuals[idx].Score -= 50
			}else if 10 > diff{
				g.GA.Individuals[idx].Score -= 100
			}else if 15 > diff{
				g.GA.Individuals[idx].Score -= 150
			}else if 20 > diff{
				g.GA.Individuals[idx].Score -= 200
			}else{
				g.GA.Individuals[idx].Score -= 300
			}
			return true
		}

		if y1 >= underBoaderLine {
			diff := y1 - underBoaderLine
			if 5 > diff{
				g.GA.Individuals[idx].Score -= 50
			}else if 10 > diff{
				g.GA.Individuals[idx].Score -= 100
			}else if 15 > diff{
				g.GA.Individuals[idx].Score -= 150
			}else if 20 > diff{
				g.GA.Individuals[idx].Score -= 200
			}else{
				g.GA.Individuals[idx].Score -= 300
			}
			return true
		}

		passedPoint = 10
	}

	g.GA.Individuals[idx].Score += passedPoint
	g.GA.Individuals[idx].Score += 1

	return false
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
    return screenWidth, screenHeight
}

func (g *Game) drawGorilla(screen *ebiten.Image) {
	// gorilla Image size
	w, h := 75, 75

	for i:= 0; i < ga.POPULATION_SIZE; i++{
		// 死んだゴリラは描画しない
		if g.GA.Individuals[i].CheckDead(){
			continue
		}
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(-float64(w)/2.0, -float64(h)/2.0)
		op.GeoM.Rotate(float64(g.gorillaVy[i]) / 96.0 * math.Pi / 6)
		op.GeoM.Translate(float64(w)/2.0, float64(h)/2.0)
		op.GeoM.Translate(float64(g.gorillaX[i]/16.0)-float64(g.cameraX), float64(g.gorillaY[i]/16.0)-float64(g.cameraY))
		op.Filter = ebiten.FilterLinear
		screen.DrawImage(gorillaImages[i], op)
	}
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

		op.GeoM.Reset()
		op.GeoM.Translate(float64(i*tileSize-utils.FloorMod(g.cameraX, tileSize)),
			float64((ny-1)*tileSize-utils.FloorMod(g.cameraY, tileSize)))
		screen.DrawImage(tilesImage.SubImage(image.Rect(0, 0, tileSize, tileSize)).(*ebiten.Image), op)

		if tileY, ok := g.pipeAt(utils.FloorDiv(g.cameraX, tileSize) + i); ok {
			for j := 0; j < tileY; j++ {
				op.GeoM.Reset()
				op.GeoM.Scale(1, -1)
				op.GeoM.Translate(float64(i*tileSize-utils.FloorMod(g.cameraX, tileSize)),
					float64(j*tileSize-utils.FloorMod(g.cameraY, tileSize)))
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
				op.GeoM.Translate(float64(i*tileSize-utils.FloorMod(g.cameraX, tileSize)),
					float64(j*tileSize-utils.FloorMod(g.cameraY, tileSize)))
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
	ebiten.SetWindowTitle("Flappy Gorilla")
    if err := ebiten.RunGame(NewGame()); err != nil {
        log.Fatal(err)
    }
}