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

var screenWidth, screenHeight int = 1920 / 4, 1080 / 4

var rows, cols int = 3, 3
var numOfTextures = rows * cols
var textureWidth = float32(screenWidth) / float32(rows) * 0.9
var textureHeight = float32(screenHeight) / float32(cols) * 0.9
var paddingWidth = float32(screenWidth) / float32(rows+1) * 0.1
var paddingHeight = float32(screenHeight) / float32(cols+1) * 0.1

type texture struct {
	equation *textureEquation
	image    *ebiten.Image
}

func NewTexture(index int) *texture {
	equation := NewTextureEquation()
	image := generateTexture(equation, int(textureWidth), int(textureHeight))
	t := &texture{equation, image}
	return t
}

func (t *texture) mutate() {
	t.equation.mutate()
	image := generateTexture(t.equation, int(textureWidth), int(textureHeight))
	t.image = image
}

type textureEquation struct {
	r eqt.BaseNode
	g eqt.BaseNode
	b eqt.BaseNode
}

func (t *textureEquation) String() string {
	return "(EquationImage \n" + t.r.String() + "\n" + t.g.String() + "\n" + t.b.String() + ")"
}

func NewTextureEquation() *textureEquation {
	opNodeCount := 20

	t := &textureEquation{}
	t.r = randomEquation(opNodeCount)
	t.g = randomEquation(opNodeCount)
	t.b = randomEquation(opNodeCount)

	return t
}

func randomEquation(nodeCount int) eqt.BaseNode {
	node := eqt.RandomOpNode()

	for i := 0; i < nodeCount; i++ {
		node.AddRandomNode(eqt.RandomOpNode())
	}

	for node.AddLeafNode(eqt.RandomLeafNode()) {
	}

	return node
}

func (t *textureEquation) mutate() {
	n := rand.Intn(3)
	switch n {
	case 0:
		node := eqt.PickRandomNode(t.r)
		mutatedNode := eqt.Mutate(node)
		if node == t.r {
			t.r = mutatedNode
		}
	case 1:
		node := eqt.PickRandomNode(t.g)
		mutatedNode := eqt.Mutate(node)
		if node == t.g {
			t.g = mutatedNode
		}
	case 2:
		node := eqt.PickRandomNode(t.b)
		mutatedNode := eqt.Mutate(node)
		if node == t.b {
			t.b = mutatedNode
		}
	}
}

func generateTexture(t *textureEquation, width, height int) *ebiten.Image {
	texture := ebiten.NewImage(width, height)

	for y := 0; y < height; y++ {
		fy := float32(y)/float32(height)*2 - 1
		for x := 0; x < width; x++ {
			fx := float32(x)/float32(width)*2 - 1

			r := t.r.Eval(fx, fy)*255 + 127
			g := t.g.Eval(fy, fx)*255 + 127
			b := t.b.Eval(fx, fy)*255 + 127
			a := 255

			pixel := &color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
			texture.Set(x, y, pixel)
		}
	}

	return texture
}

// Game implements ebiten.Game interface.
type Game struct {
	textures []*texture
}

func NewGame() *Game {
	textures := make([]*texture, numOfTextures)
	for i := range numOfTextures {
		textures[i] = NewTexture(i)
	}
	return &Game{textures: textures}
}

// Update proceeds the game state.
// Update is called every tick (1/60 [s] by default).
func (g *Game) Update() error {
	// Write your game's logical update.
	keySpace := ebiten.KeySpace
	if inpututil.IsKeyJustPressed(keySpace) {
		for _, tex := range g.textures {
			tex.mutate()
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return ebiten.Termination
	}

	return nil
}

// Draw draws the game screen.
// Draw is called every frame (typically 1/60[s] for 60Hz display).
func (g *Game) Draw(screen *ebiten.Image) {
	// Write your game's rendering.
	for row := 0; row < rows; row++ {
		rowf := float32(row)
		for col := 0; col < cols; col++ {
			colf := float32(col)
			op := &ebiten.DrawImageOptions{}
			xOffset := paddingWidth*(colf+1) + textureWidth*colf
			yOffset := paddingHeight*(rowf+1) + textureHeight*rowf
			op.GeoM.Translate(float64(xOffset), float64(yOffset))
			index := row*cols + col
			screen.DrawImage(g.textures[index].image, op)
		}
	}
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
