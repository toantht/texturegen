package main

import (
	"image/color"
	"log"
	"time"

	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	eqt "github.com/toantht/texturegen/equation"
)

const screenWidth, screenHeight int = 1920 / 4, 1080 / 4

type EquationFunc func(x, y float32) float32

func randomEquation(nodeCount int) eqt.BaseNode {
	node := eqt.RandomOpNode()

	for i := 0; i < nodeCount; i++ {
		node.AddRandomNode(eqt.RandomOpNode())
	}

	for node.AddLeafNode(eqt.RandomLeafNode()) {
	}

	return node
}

func generateTexture(width, height int) *ebiten.Image {
	texture := ebiten.NewImage(width, height)

	nodeCount := 8
	rEquation := randomEquation(nodeCount)
	gEquation := randomEquation(nodeCount)
	bEquation := randomEquation(nodeCount)

	for y := 0; y < height; y++ {
		fy := float32(y)/float32(height)*2 - 1
		for x := 0; x < width; x++ {
			fx := float32(x)/float32(width)*2 - 1

			r := rEquation.Eval(fx, fy)*255 + 127
			g := gEquation.Eval(fy, fx)*255 + 127
			b := bEquation.Eval(fx, fy)*255 + 127
			a := 255

			pixel := &color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
			texture.Set(x, y, pixel)
		}
	}

	return texture
}

// Game implements ebiten.Game interface.
type Game struct {
	texture *ebiten.Image
}

func NewGame() *Game {
	texture := generateTexture(screenWidth, screenHeight)
	g := &Game{texture: texture}
	return g
}

// Update proceeds the game state.
// Update is called every tick (1/60 [s] by default).
func (g *Game) Update() error {
	// Write your game's logical update.
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.texture = generateTexture(screenWidth, screenHeight)
	}
	return nil
}

// Draw draws the game screen.
// Draw is called every frame (typically 1/60[s] for 60Hz display).
func (g *Game) Draw(screen *ebiten.Image) {
	// Write your game's rendering.
	screen.DrawImage(g.texture, nil)
}

// Layout takes the outside size (e.g., the window size) and returns the (logical) screen size.
// If you don't have to adjust the screen size with the outside size, just return a fixed size.
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("TextureGen")

	rand.Seed(time.Now().UnixNano())

	game := NewGame()
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
