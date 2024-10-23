package main

import (
	"image"
	"image/color"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Sprite struct {
	img          *ebiten.Image
	X, Y, Dx, Dy float64
}

type Villain struct {
	*Sprite
}
type Game struct {
	Hero      *Sprite
	villain   []*Villain
	caught    bool
	colliders []image.Rectangle

	// Paused  bool
}

func CheckCollisionHorizontal(sprite *Sprite, coliders []image.Rectangle) {
	for _, col := range coliders {
		if col.Overlaps(
			image.Rect(
				int(sprite.X),
				int(sprite.Y),
				int(sprite.X)+16.0,
				int(sprite.Y)+16.0),
		) {
			if sprite.Dx > 0.0 {
				sprite.X = float64(col.Min.X) - 16.0
			} else if sprite.Dx < 0.0 {
				sprite.X = float64(col.Max.X)
			}
		}
	}
}

func CheckCollisionVertical(sprite *Sprite, coliders []image.Rectangle) {
	for _, col := range coliders {
		if col.Overlaps(
			image.Rect(
				int(sprite.X),
				int(sprite.Y),
				int(sprite.X)+16.0,
				int(sprite.Y)+16.0),
		) {
			if sprite.Dy > 0.0 {
				sprite.Y = float64(col.Min.Y) - 16.0
			} else if sprite.Dy < 0.0 {
				sprite.Y = float64(col.Max.Y)
			}
		}
	}
}

func (g *Game) Update() error {
	// if ebiten.IsKeyPressed(ebiten.KeyTab) {
	// 	g.Paused = !g.Paused

	// }

	// if g.Paused {
	// 	return nil
	// }
	g.Hero.Dx = 0.0
	g.Hero.Dy = 0.0

	g.caught = false

	if ebiten.IsKeyPressed(ebiten.KeyRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
		g.Hero.Dx += 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
		g.Hero.Dx -= 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) || ebiten.IsKeyPressed(ebiten.KeyW) {
		g.Hero.Dy -= 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) || ebiten.IsKeyPressed(ebiten.KeyS) {
		g.Hero.Dy += 2
	}

	g.Hero.X += g.Hero.Dx
	CheckCollisionHorizontal(g.Hero, g.colliders)
	g.Hero.Y += g.Hero.Dy
	CheckCollisionVertical(g.Hero, g.colliders)

	// villain
	for _, sprite := range g.villain {
		sprite.Dx = 0.0
		sprite.Dy = 0.0

		if sprite.X < g.Hero.X {
			sprite.Dx += 1

		} else if sprite.X > g.Hero.X {
			sprite.Dx -= 1
		}
		if sprite.Y < g.Hero.Y {
			sprite.Y += 1
		} else if sprite.Y > g.Hero.Y {
			sprite.Y -= 1
		}
		sprite.X += sprite.Dx
		sprite.Y += sprite.Dy
		CheckCollisionHorizontal(sprite.Sprite, g.colliders)
		CheckCollisionVertical(sprite.Sprite, g.colliders)

		// caught
		// if math.Abs(sprite.X-g.Hero.X) < 15 && math.Abs(sprite.Y-g.Hero.Y) < 15 {
		// 	g.caught = true
		// }
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// screen.Fill(color.RGBA{120, 180, 255, 255})
	screen.Fill(color.White)
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
	opts.GeoM.Reset()
	// border
	g.Border(g.Hero.X, g.Hero.Y)
	if g.caught {
		ebitenutil.DebugPrint(screen, "you get cought")
	}
	// if g.Paused {
	// 	ebitenutil.DebugPrint(screen, "you have paused")
	// }

	for _, col := range g.colliders {
		vector.StrokeRect(screen,
			float32(col.Min.X),
			float32(col.Min.Y),
			float32(col.Dx()),
			float32(col.Dy()),
			1.0,
			color.RGBA{255, 0, 0, 255},
			true)

	}

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
		colliders: []image.Rectangle{
			image.Rect(100, 180, 116, 116),
			image.Rect(100, 280, 116, 216),
			// 	image.Rect(200, 680, 116, 116),
			// 	image.Rect(700, 380, 116, 116),
			// 	image.Rect(900, 80, 116, 116),
		},
	}

	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}
}
