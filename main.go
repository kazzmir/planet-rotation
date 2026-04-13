package main

import (
    _ "fmt"
    "os"
    "image"
    "image/color" 
    "sync"
    "math/rand/v2"
    _ "image/png"
    _ "image/jpeg"
    "github.com/hajimehoshi/ebiten/v2"
    "github.com/hajimehoshi/ebiten/v2/inpututil"
    "github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

func draw(screen *ebiten.Image, planetImage *ebiten.Image, cloudImage *ebiten.Image, timeSeconds float64, shader *ebiten.Shader) {

    /*
    var opts2 ebiten.DrawImageOptions
    opts2.GeoM.Translate(10, 200)
    opts2.GeoM.Scale(0.5, 0.5)
    screen.DrawImage(cloudImage, &opts2)
    */

    bounds := planetImage.Bounds()
    w, h := bounds.Dx(), bounds.Dy()

    opts := &ebiten.DrawRectShaderOptions{}
    opts.GeoM.Translate(10, 100)
    opts.GeoM.Scale(0.5, 0.5)

    // fmt.Printf("Resolution: %v x %v\n", w, h)

    rotationSpeed := timeSeconds / 600.0

    opts.Uniforms = map[string]interface{}{
        "Rotation":       float32(rotationSpeed),
    }
    opts.Images[0] = planetImage

    screen.DrawRectShader(w, h, shader, opts)

    opts.Blend = ebiten.BlendLighter
    opts.Images[0] = cloudImage
    opts.ColorScale.ScaleAlpha(0.2)
    opts.Uniforms["Rotation"] = float32(rotationSpeed * 1.5)
    screen.DrawRectShader(w, h, shader, opts)
}

type Game struct {
    Counter uint64
    Shader *ebiten.Shader
    Planet *ebiten.Image
    CloudImage *ebiten.Image
    Init sync.Once
    drawClouds int
}

func makeCloudImage(bounds image.Rectangle) *ebiten.Image {
    w, h := bounds.Dx(), bounds.Dy()
    img := ebiten.NewImage(w, h)
    // img.Fill(color.RGBA{B: 255, A: 255})

    /*
    if 2 > 1 {
        img2, _, _ := ebitenutil.NewImageFromFile("rect2.png")
        if 2 > 1 {
            return img2
        }
        var opts ebiten.DrawImageOptions
        img.DrawImage(img2, &opts)

        return img
    }
    */

    clouds := []string{"cloud1.png", "cloud-a.png"}

    for _, cloudFile := range clouds {
        cloud, _, err := ebitenutil.NewImageFromFile(cloudFile)
        if err != nil {
            panic(err)
        }

        cloudBounds := cloud.Bounds()

        for range 10 {
            var opts ebiten.DrawImageOptions
            opts.GeoM.Translate(float64(w - cloudBounds.Dx()) * rand.Float64(), float64(h - cloudBounds.Dy()) * rand.Float64())
            img.DrawImage(cloud, &opts)
        }
    }

    return img
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

    // planetImage, _, err := ebitenutil.NewImageFromFile("rect2.png")
    planetImage, _, err := ebitenutil.NewImageFromFile("mars.jpg")
    if err != nil {
        panic(err)
    }

    cloudImage := makeCloudImage(planetImage.Bounds())

    return &Game{
        Counter: 0,
        Shader: shader,
        Planet: planetImage,
        drawClouds: 800,
        CloudImage: cloudImage,
    }
}

func (g *Game) Update() error {
    g.Counter += 1

    keys := inpututil.AppendJustPressedKeys(nil)
    for _, k := range keys {
        if k == ebiten.KeyEscape || k == ebiten.KeyCapsLock {
            return ebiten.Termination
        }
    }

    return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
    screen.Fill(color.NRGBA{R: 32, G: 0, B: 0, A: 255})

    // if g.drawClouds > 0 {
        draw(screen, g.Planet, g.CloudImage, float64(g.Counter), g.Shader)
    // }

    /*
    var opts ebiten.DrawImageOptions
    opts.GeoM.Translate(10, 200)
    opts.GeoM.Scale(0.5, 0.5)
    screen.DrawImage(g.CloudImage, &opts)
    */
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
    // return g.Planet.Bounds().Dx(), g.Planet.Bounds().Dy()
    return outsideWidth, outsideHeight
}

func main() {

    ebiten.SetWindowSize(1024, 768)
    ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
    ebiten.RunGame(MakeGame())
}
