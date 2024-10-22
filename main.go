package main

import (
	"image"
	"image/color"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Sprite struct {
	img  *ebiten.Image
	X, Y float64
}

type Villain struct {
	*Sprite
}
type Game struct {
	Hero    *Sprite
	villain []*Villain
}

func (g *Game) Update() error {

	if ebiten.IsKeyPressed(ebiten.KeyRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
		g.Hero.X += 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
		g.Hero.X -= 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) || ebiten.IsKeyPressed(ebiten.KeyW) {
		g.Hero.Y -= 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) || ebiten.IsKeyPressed(ebiten.KeyS) {
		g.Hero.Y += 2
	}

	for _, sprite := range g.villain {

		if sprite.X <= g.Hero.X {
			sprite.X += 1
		} else if sprite.X >= g.Hero.X {
			sprite.X -= 1
		}

		if sprite.Y < g.Hero.Y {
			sprite.Y += 1
		} else if sprite.Y > g.Hero.Y {
			sprite.Y -= 1
		}

	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{120, 180, 255, 255})
	// screen.Fill(color.White)
	opts := ebiten.DrawImageOptions{}

	opts.GeoM.Translate(g.Hero.X, g.Hero.Y)

	screen.DrawImage(g.Hero.img.SubImage(
		image.Rect(0, 0, 16, 16)).(*ebiten.Image),
		&opts)
	opts.GeoM.Reset()

	// renduring multiple villains
	for _, sprite := range g.villain {
		opts.GeoM.Translate(sprite.X, sprite.Y)

		screen.DrawImage(
			sprite.img.SubImage(
				image.Rect(0, 0, 16, 16)).(*ebiten.Image),
			&opts)
		opts.GeoM.Reset()
	}

	// border
	g.Border(g.Hero.X, g.Hero.Y)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ebiten.WindowSize()
}

func (g *Game) Border(width, hight float64) {
	g.Hero.X = math.Max(g.Hero.X, 0.0)

	g.Hero.Y = math.Max(g.Hero.Y, 0.0)

	g.Hero.X = math.Min(g.Hero.X, 620)
	g.Hero.Y = math.Min(g.Hero.Y, 460)

}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Runner Game!")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	hero, _, err := ebitenutil.NewImageFromFile("asset/images/hero.png")
	if err != nil {
		log.Fatal(err)
	}

	villain, _, err := ebitenutil.NewImageFromFile("asset/images/villain.png")
	if err != nil {
		log.Fatal(err)
	}
	game := Game{
		Hero: &Sprite{
			img: hero,
			X:   50.0,
			Y:   50.0,
		},
		villain: []*Villain{

			{
				&Sprite{
					img: villain,
					X:   175.0,
					Y:   175.0,
				},
			},

			{
				&Sprite{
					img: villain,
					X:   50.0,
					Y:   50.0,
				},
			},

			{
				&Sprite{
					img: villain,
					X:   100.0,
					Y:   100.0,
				},
			},
		},
	}

	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}
}
