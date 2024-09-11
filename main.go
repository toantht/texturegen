package main

import (
	"fmt"
	"image/color"
	"log"
	"math"
	"time"

	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	eqt "github.com/toantht/texturegen/equation"
)

const screenWidth, screenHeight int = 1920 / 4, 1080 / 4

type EquationFunc func(x, y float32) float32

func randomEquation(depth int) EquationFunc {
	if depth <= 0 {
		switch rand.Intn(3) {
		case 0:
			return func(x, y float32) float32 { return x }
		case 1:
			return func(x, y float32) float32 { return y }
		default:
			constant := rand.Float32()*2 - 1
			return func(x, y float32) float32 { return constant }
		}
	}

	switch rand.Intn(4) {
	case 0:
		f1 := randomEquation(depth - 1)
		f2 := randomEquation(depth - 1)
		return func(x, y float32) float32 {
			return (f1(x, y) + f2(x, y)) / 2
		}
	case 1:
		f1 := randomEquation(depth - 1)
		f2 := randomEquation(depth - 1)
		return func(x, y float32) float32 {
			return f1(x, y) * f2(x, y)
		}
	case 2:
		f := randomEquation(depth - 1)
		return func(x, y float32) float32 {
			r64 := float64(f(x, y))
			return float32(math.Sin(r64))
		}
	case 3:
		f := randomEquation(depth - 1)
		return func(x, y float32) float32 {
			r64 := float64(f(x, y))
			return float32(math.Cos(r64))
		}
	}
	panic("randomize equation failed")
}

func generateTexture(width, height int) *ebiten.Image {
	texture := ebiten.NewImage(width, height)

	opX := eqt.NewOpX()
	opY := eqt.NewOpY()

	sin := eqt.NewOpSin()
	sin.Children[0] = opX

	cos := eqt.NewOpCos()
	cos.Children[0] = opY

	mult := eqt.NewOpMult()
	mult.Children[0] = sin
	mult.Children[1] = opY

	atan2 := eqt.NewOpAtan2()
	atan2.Children[0] = opX
	atan2.Children[1] = cos

	plus := eqt.NewOpPlus()
	plus.Children[0] = mult
	plus.Children[1] = atan2
	fmt.Printf("plus: %v\n", plus)

	constant := eqt.NewOpConstant()
	fmt.Printf("constant: %v\n", constant)

	for y := 0; y < height; y++ {
		fy := float32(y) / 20
		for x := 0; x < width; x++ {
			fx := float32(x) / 20

			r := plus.Eval(fx, fy) * 255
			g := plus.Eval(fy, fx) * 255
			b := constant.Eval(fx, fy) * 255
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
