package main

import (
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

var rows, cols int = 3, 3
var numOfTextures = rows * cols
var textureWidth = float32(screenWidth/cols) * 0.9
var textureHeight = float32((screenHeight-50)/rows) * 0.9
var paddingWidth = float32(screenWidth) * 0.1 / float32(cols+1)
var paddingHeight = float32(screenHeight) * 0.1 / float32(rows+1)
var borderWidth = min(textureWidth/20, 2)

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

func crossover(a *textureEquation, b *textureEquation) *textureEquation {
	aq := copyTextureEquation(a)
	aColor := aq.pickRandomColor()
	aNode := eqt.PickRandomNode(aColor)

	bColor := b.pickRandomColor()
	bNode := eqt.PickRandomNode(bColor)

	eqt.ReplaceNode(aNode, bNode)
	return aq
}

func evolve(selectedEquations []*textureEquation) []*textureEquation {
	eqs := make([]*textureEquation, numOfTextures)

	n := len(selectedEquations)
	i := 0
	for i < numOfTextures {
		a := selectedEquations[rand.Intn(n)]
		b := selectedEquations[rand.Intn(n)]
		eqs[i] = crossover(a, b)
		i++
	}

	for _, eq := range eqs {
		n := rand.Intn(4)
		for i := 0; i < n; i++ {
			eq.mutate()
		}
	}
	return eqs
}

func (t *texture) applyEquation(e *textureEquation) {
	t.equation = e
	image := generateTexture(t.equation, int(textureWidth), int(textureHeight))
	t.image = image
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

func (t *texture) draw(screen *ebiten.Image) {
	if t.selected {
		border := ebiten.NewImage(int(textureWidth+borderWidth*2), int(textureHeight+borderWidth*2))
		border.Fill(color.RGBA{255, 255, 0, 255})
		borderOp := &ebiten.DrawImageOptions{}
		borderOp.GeoM.Translate(float64(t.x)-float64(borderWidth), float64(t.y)-float64(borderWidth))
		screen.DrawImage(border, borderOp)
	}

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(t.x), float64(t.y))
	screen.DrawImage(t.image, op)
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
	opNodeCount := rand.Intn(100) + 1

	t := &textureEquation{}
	t.r = randomEquation(opNodeCount)
	t.g = randomEquation(opNodeCount)
	t.b = randomEquation(opNodeCount)

	return t
}

func (t *textureEquation) pickRandomColor() eqt.BaseNode {
	n := rand.Intn(3)
	switch n {
	case 0:
		return t.r
	case 1:
		return t.g
	case 2:
		return t.b
	}
	panic("pick random failed")
}

func copyTextureEquation(t *textureEquation) *textureEquation {
	result := &textureEquation{eqt.CopyTree(t.r), eqt.CopyTree(t.g), eqt.CopyTree(t.b)}
	return result
}

func randomEquation(opNodeCount int) eqt.BaseNode {
	if opNodeCount < 1 {
		return nil
	}

	node := eqt.RandomOpNode()

	for i := 1; i < opNodeCount; i++ {
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
				if tex != nil && tex.selected {
					tex.mutate()
					tex.selected = false
				}
			}()
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return ebiten.Termination
	}

	if g.button.IsClicked() {
		selectedEquations := make([]*textureEquation, 0)
		for _, t := range g.textures {
			if t != nil && t.selected {
				selectedEquations = append(selectedEquations, t.equation)
			}
		}
		if len(selectedEquations) > 0 {
			eqs := evolve(selectedEquations)
			for i := range g.textures {
				g.textures[i].applyEquation(eqs[i])
				g.textures[i].selected = false
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
			tex.draw(screen)
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
