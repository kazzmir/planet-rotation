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
    // "github.com/hajimehoshi/ebiten/v2/vector"
)

func draw(screen *ebiten.Image, x float64, y float64, scale float64, planetImage *ebiten.Image, cloudImage *ebiten.Image, timeSeconds float64, shader *ebiten.Shader) {

    /*
    var opts2 ebiten.DrawImageOptions
    opts2.GeoM.Translate(10, 200)
    opts2.GeoM.Scale(0.5, 0.5)
    screen.DrawImage(cloudImage, &opts2)
    */

    bounds := planetImage.Bounds()
    w, h := bounds.Dx(), bounds.Dy()

    opts := &ebiten.DrawRectShaderOptions{}
    opts.GeoM.Translate(-float64(w) * 0.5, -float64(h) * 0.5)
    opts.GeoM.Scale(scale, scale)
    opts.GeoM.Translate(x, y)

    // fmt.Printf("Resolution: %v x %v\n", w, h)

    rotationSpeed := timeSeconds / 600.0

    opts.Uniforms = map[string]interface{}{
        "Rotation":       float32(rotationSpeed),
    }
    opts.Images[0] = planetImage

    screen.DrawRectShader(w, h, shader, opts)

    if cloudImage != nil {
        opts.Blend = ebiten.BlendLighter
        opts.Images[0] = cloudImage
        opts.ColorScale.ScaleAlpha(0.2)
        opts.Uniforms["Rotation"] = float32(rotationSpeed * 1.5)
        screen.DrawRectShader(w, h, shader, opts)
    }
}

type Planet int

const (
    Earth Planet = iota
    Mars
)

func (planet Planet) Next() Planet {
    switch planet {
    case Earth:
        return Mars
    case Mars:
        return Earth
    }

    return Earth
}

func (planet Planet) Previous() Planet {
    switch planet {
    case Earth:
        return Mars
    case Mars:
        return Earth
    }

    return Earth
}

type Game struct {
    Counter uint64
    Shader *ebiten.Shader
    Planet *ebiten.Image
    CurrentPlanet Planet
    CloudImage *ebiten.Image
    Init sync.Once
    drawClouds bool
    Scale float64
}

func makeCloudImage(bounds image.Rectangle) *ebiten.Image {
    w, h := bounds.Dx(), bounds.Dy()
    img := ebiten.NewImage(w, h)
    // img.Fill(color.RGBA{B: 255, A: 255})

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

    planetImage, _, err := ebitenutil.NewImageFromFile("earth.jpg")
    // planetImage, _, err := ebitenutil.NewImageFromFile("mars.jpg")
    if err != nil {
        panic(err)
    }

    cloudImage := makeCloudImage(planetImage.Bounds())

    return &Game{
        Counter: 0,
        Shader: shader,
        Planet: planetImage,
        CurrentPlanet: Earth,
        drawClouds: true,
        CloudImage: cloudImage,
        Scale: 0.5,
    }
}

func (g *Game) Update() error {
    g.Counter += 1

    keys := inpututil.AppendJustPressedKeys(nil)
    for _, k := range keys {
        if k == ebiten.KeyEscape || k == ebiten.KeyCapsLock {
            return ebiten.Termination
        }
        if k == ebiten.KeySpace {
            g.drawClouds = !g.drawClouds
        }

        change := false
        if k == ebiten.KeyLeft {
            g.CurrentPlanet = g.CurrentPlanet.Previous()
            change = true
        } else if k == ebiten.KeyRight {
            g.CurrentPlanet = g.CurrentPlanet.Next()
            change = true
        }

        if change {
            switch g.CurrentPlanet {
                case Earth:
                    g.Planet, _, _ = ebitenutil.NewImageFromFile("earth.jpg")
                    g.CloudImage = makeCloudImage(g.Planet.Bounds())
                case Mars:
                    g.Planet, _, _ = ebitenutil.NewImageFromFile("mars.jpg")
                    g.CloudImage = makeCloudImage(g.Planet.Bounds())
            }
        }
    }

    _, yWheel := ebiten.Wheel()
    if yWheel * yWheel > 0 {
        g.Scale += yWheel * 0.05
    }

    return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
    screen.Fill(color.NRGBA{R: 32, G: 0, B: 0, A: 255})

    x := float64(screen.Bounds().Dx() / 2)
    y := float64(screen.Bounds().Dy() / 2)

    /*
    vector.StrokeLine(screen, float32(x), 0, float32(x), float32(screen.Bounds().Dy()), 1, color.White, false)
    vector.StrokeLine(screen, 0, float32(y), float32(screen.Bounds().Dx()), float32(y), 1, color.White, false)
    */

    cloud := g.CloudImage
    if !g.drawClouds {
        cloud = nil
    }

    draw(screen, x, y, g.Scale, g.Planet, cloud, float64(g.Counter), g.Shader)
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
