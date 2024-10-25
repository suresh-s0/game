package main

import (
	"image"
	"image/color"
	"log"
	"math"
	"math/rand"

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
	villains  []*Villain
	caught    bool
	colliders []image.Rectangle
}

func CheckCollisionHorizontal(sprite *Sprite, colliders []image.Rectangle) {
	for _, col := range colliders {
		if col.Overlaps(
			image.Rect(
				int(sprite.X),
				int(sprite.Y),
				int(sprite.X)+16,
				int(sprite.Y)+16,
			),
		) {
			if sprite.Dx > 0 {
				sprite.X = float64(col.Min.X) - 16
			} else if sprite.Dx < 0 {
				sprite.X = float64(col.Max.X)
			}
		}
	}
}

func CheckCollisionVertical(sprite *Sprite, colliders []image.Rectangle) {
	for _, col := range colliders {
		if col.Overlaps(
			image.Rect(
				int(sprite.X),
				int(sprite.Y),
				int(sprite.X)+16,
				int(sprite.Y)+16,
			),
		) {
			if sprite.Dy > 0 {
				sprite.Y = float64(col.Min.Y) - 16
			} else if sprite.Dy < 0 {
				sprite.Y = float64(col.Max.Y)
			}
		}
	}
}
func generateRandomColliders(num int, maxX, maxY int) []image.Rectangle {
	colliders := make([]image.Rectangle, 0, num)

	for i := 0; i < num; i++ {
		var newRect image.Rectangle
		valid := false

		for !valid {

			x := rand.Intn(maxX)
			y := rand.Intn(maxY)

			width := rand.Intn(30) + 20
			height := rand.Intn(150) + 20

			// Ensure the rectangle is within the window bounds
			if x+width > maxX {
				width = maxX - x
			}
			if y+height > maxY {
				height = maxY - y
			}

			// Create the rectangle
			newRect = image.Rect(x, y, x+width, y+height)

			// Check for overlap with existing colliders
			valid = true
			for _, col := range colliders {
				if newRect.Overlaps(col) {
					valid = false
					break
				}
			}
		}

		// Add the non-overlapping rectangle to the list
		colliders = append(colliders, newRect)
	}

	return colliders
}

func (g *Game) Border() {
	g.Hero.X = math.Max(g.Hero.X, 0.0)
	g.Hero.Y = math.Max(g.Hero.Y, 0.0)
	g.Hero.X = math.Min(g.Hero.X, 620)
	g.Hero.Y = math.Min(g.Hero.Y, 460)
}

func (g *Game) Update() error {
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

	for i, villain := range g.villains {
		villain.Dx = 0.0
		villain.Dy = 0.0

		// Calculate distance from Hero
		xDistance := villain.X - g.Hero.X
		yDistance := villain.Y - g.Hero.Y
		speed := 1.0

		// Move villain towards Hero
		if xDistance > 1.0 {
			villain.Dx -= speed
		} else if xDistance < -1.0 {
			villain.Dx += speed
		}

		if yDistance > 1.0 {
			villain.Dy -= speed
		} else if yDistance < -1.0 {
			villain.Dy += speed
		}

		// Check distance to other villains to avoid overlap
		for j, otherVillain := range g.villains {
			if i != j {
				// Calculate the distance between two villains
				villainDistanceX := villain.X - otherVillain.X
				villainDistanceY := villain.Y - otherVillain.Y
				villainDistance := math.Sqrt(villainDistanceX*villainDistanceX + villainDistanceY*villainDistanceY)

				// Ensure they maintain a minimum distance of 20 pixels
				minDistance := 20.0
				if villainDistance < minDistance {
					// Adjust direction to maintain distance
					if villainDistanceX > 0 {
						villain.Dx += speed / 2 // move right
					} else {
						villain.Dx -= speed / 2 // move left
					}
					if villainDistanceY > 0 {
						villain.Dy += speed / 2 // move down
					} else {
						villain.Dy -= speed / 2 // move up
					}
				}
			}
		}

		// Update villain position
		villain.X += villain.Dx
		CheckCollisionHorizontal(villain.Sprite, g.colliders)
		villain.Y += villain.Dy
		CheckCollisionVertical(villain.Sprite, g.colliders)

		// Check if Hero is caught
		if math.Abs(villain.X-g.Hero.X) < 15 && math.Abs(villain.Y-g.Hero.Y) < 15 {
			g.caught = true
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.White)

	opts := ebiten.DrawImageOptions{}

	opts.GeoM.Translate(g.Hero.X, g.Hero.Y)

	screen.DrawImage(g.Hero.img.SubImage(
		image.Rect(0, 0, 16, 16)).(*ebiten.Image),
		&opts)
	opts.GeoM.Reset()

	for _, villain := range g.villains {
		opts.GeoM.Translate(villain.X, villain.Y)
		screen.DrawImage(villain.img.SubImage(image.Rect(0, 0, 16, 16)).(*ebiten.Image), &opts)
		opts.GeoM.Reset()
	}

	g.Border()
	if g.caught {
		ebitenutil.DebugPrint(screen, "You got caught!")
	}

	for _, col := range g.colliders {
		vector.DrawFilledRect(screen,
			float32(col.Min.X),
			float32(col.Min.Y),
			float32(col.Dx()),
			float32(col.Dy()),
			color.RGBA{0, 0, 5, 200},
			true,
		)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ebiten.WindowSize()
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
			X:   300.0,
			Y:   400.0,
		},
		villains: []*Villain{
			{&Sprite{img: villain, X: 0.0, Y: 0.0}},
			{&Sprite{img: villain, X: 50.0, Y: 150.0}},
			{&Sprite{img: villain, X: 500.0, Y: 300.0}},
		},
		colliders: generateRandomColliders(25, 640, 480),
	}

	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}
}
