package main

import (
    "os"
    _ "image/png"
    "github.com/hajimehoshi/ebiten/v2"
    "github.com/hajimehoshi/ebiten/v2/inpututil"
    "github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

func draw(screen *ebiten.Image, planetImage *ebiten.Image, timeSeconds uint64, shader *ebiten.Shader) {

    bounds := planetImage.Bounds()
    w, h := bounds.Dx(), bounds.Dy()

    opts := &ebiten.DrawRectShaderOptions{}
    opts.Uniforms = map[string]interface{}{
        "Time":       float32(float64(timeSeconds) / 200.0),
        "Resolution": []float32{float32(w), float32(h)},
    }
    opts.Images[0] = planetImage

    screen.DrawRectShader(w, h, shader, opts)
}

type Game struct {
    Counter uint64
    Shader *ebiten.Shader
    Planet *ebiten.Image
}

func MakeGame() *Game {
    data, err := os.ReadFile("planet.kage")
    if err != nil {
        panic(err)
    }
    shader, err := ebiten.NewShader(data)
    if err != nil {
        panic(err)
    }

    planetImage, _, err := ebitenutil.NewImageFromFile("world1.png")
    if err != nil {
        panic(err)
    }

    return &Game{
        Counter: 0,
        Shader: shader,
        Planet: planetImage,
    }
}

func (g *Game) Update() error {
    g.Counter++

    keys := inpututil.AppendJustPressedKeys(nil)
    for _, k := range keys {
        if k == ebiten.KeyEscape || k == ebiten.KeyCapsLock {
            return ebiten.Termination
        }
    }

    return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
    draw(screen, g.Planet, g.Counter, g.Shader)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
    return g.Planet.Bounds().Dx(), g.Planet.Bounds().Dy()
}

func main() {

    ebiten.RunGame(MakeGame())
}
