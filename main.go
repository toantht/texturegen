package main

import (
	"fmt"
	"image/color"
	"log"
	"time"

	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	eqt "github.com/toantht/texturegen/equation"
	"github.com/toantht/texturegen/gui"
)

var screenWidth, screenHeight int = 1920 / 4, 1080 / 4

var rows, cols int = 10, 20
var numOfTextures = rows * cols
var textureWidth = float32(screenWidth/cols) * 0.9
var textureHeight = float32((screenHeight-50)/rows) * 0.9
var paddingWidth = float32(screenWidth) * 0.1 / float32(cols+1)
var paddingHeight = float32(screenHeight) * 0.1 / float32(rows+1)

type texture struct {
	index    int
	equation *textureEquation
	image    *ebiten.Image
	x, y     int
	selected bool
}

func NewTexture(index int) *texture {
	equation := NewTextureEquation()
	image := generateTexture(equation, int(textureWidth), int(textureHeight))
	col := index % cols
	row := index / cols
	x := paddingWidth*float32(col+1) + textureWidth*float32(col)
	y := paddingHeight*float32(row+1) + textureHeight*float32(row)

	t := &texture{index, equation, image, int(x), int(y), false}
	return t
}

func (t *texture) mutate() {
	t.equation.mutate()
	image := generateTexture(t.equation, int(textureWidth), int(textureHeight))
	t.image = image
}

func (t *texture) update() {
	mx, my := ebiten.CursorPosition()
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0) {
		if mx >= t.x && mx <= t.x+int(textureWidth) && my >= t.y && my <= t.y+int(textureHeight) {
			t.selected = !t.selected
		}
	}
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
	opNodeCount := rand.Intn(100)

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
	texturesChannel chan *texture
	textures        []*texture
	button          *gui.Button
}

func NewGame() *Game {
	texChan := make(chan *texture)
	textures := make([]*texture, numOfTextures)

	button := gui.NewButton((screenWidth-80)/2, (screenHeight - 40), 80, 30)

	for i := range numOfTextures {
		go func(i int) {
			texChan <- NewTexture(i)
		}(i)
	}
	return &Game{textures: textures, texturesChannel: texChan, button: button}
}

// Update proceeds the game state.
// Update is called every tick (1/60 [s] by default).
func (g *Game) Update() error {
	// Write your game's logical update.
	keySpace := ebiten.KeySpace
	if inpututil.IsKeyJustPressed(keySpace) {
		for _, tex := range g.textures {
			go func() {
				if tex != nil {
					tex.mutate()
				}
			}()
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return ebiten.Termination
	}

	if g.button.IsClicked() {
		fmt.Println("Clicked")
		for _, t := range g.textures {
			if t != nil && t.selected {
				t.mutate()
			}
		}
	}

	for _, t := range g.textures {
		if t != nil {
			t.update()
		}
	}

	return nil
}

// Draw draws the game screen.
// Draw is called every frame (typically 1/60[s] for 60Hz display).
func (g *Game) Draw(screen *ebiten.Image) {
	// Write your game's rendering.
	select {
	case tex, ok := <-g.texturesChannel:
		if ok {
			g.textures[tex.index] = tex
		}
	default:
	}

	for _, tex := range g.textures {
		if tex != nil {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(tex.x), float64(tex.y))
			screen.DrawImage(tex.image, op)
		}
	}
	g.button.Draw(screen)
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
