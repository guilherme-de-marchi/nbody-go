package main

import (
	"fmt"
	"image/color"
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	SCREEN_WIDTH  = 800
	SCREEN_HEIGHT = 800
)

var (
	square *ebiten.Image
)

type RandOpt struct {
	ObjectQtt   int
	MassR, radR [2]float64
}

type EditOpt struct {
	ObjectsDesloc    int
	GconstDesloc     float64
	Zoom, ZoomDesloc float64
	Offset           Coordinates2D
	OffsetDesloc     float64
}

type Game struct {
	Universe *Universe

	RandOpt
	EditOpt

	keys []ebiten.Key
}

func (g *Game) Update() error {
	g.keys = inpututil.AppendPressedKeys(g.keys[:0])

	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		g.Universe = NewRandomUniverse(
			g.Universe.Size,
			g.Universe.Gconst,
			g.MassR,
			g.radR,
			g.ObjectQtt,
		)
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
		if inpututil.KeyPressDuration(ebiten.KeyG) > 0 {
			// Increase gravitational constant of the universe
			g.Universe.Gconst *= g.GconstDesloc
		} else if inpututil.KeyPressDuration(ebiten.KeyO) > 0 {
			// Generate more random objects
			g.Universe.AddObjects(GetRandomObjects(
				[2]float64{0, g.Universe.Size.X},
				[2]float64{0, g.Universe.Size.Y},
				[2]float64{1, 1000},
				[2]float64{1, 1},
				g.ObjectsDesloc,
			)...)
		} else {
			g.Zoom /= g.ZoomDesloc
		}

		return nil
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
		if inpututil.KeyPressDuration(ebiten.KeyG) > 0 {
			// Decrease gravitational constant of the universe
			g.Universe.Gconst /= g.GconstDesloc
		} else if inpututil.KeyPressDuration(ebiten.KeyO) > 0 {
			// Remove some objects from universe
			if len(g.Universe.Objects) == 0 {
				return nil
			}

			if r := len(g.Universe.Objects) - g.ObjectsDesloc; r <= 0 {
				g.Universe.Objects = []*Object{}
			} else {
				g.Universe.Objects = g.Universe.Objects[:r]
			}
		} else {
			g.Zoom *= g.ZoomDesloc
		}

		return nil
	}

	for _, k := range g.keys {
		if k == ebiten.KeyW {
			g.Offset.Y -= g.OffsetDesloc
		}
		if k == ebiten.KeyS {
			g.Offset.Y += g.OffsetDesloc
		}
		if k == ebiten.KeyA {
			g.Offset.X -= g.OffsetDesloc
		}
		if k == ebiten.KeyD {
			g.Offset.X += g.OffsetDesloc
		}
	}

	g.Universe.ApplyGravity()
	// printObjects(g.Universe)
	// time.Sleep(time.Millisecond * 10)
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Clear()
	screen.Fill(color.Black)

	rx := SCREEN_WIDTH / (g.Universe.Size.X * g.EditOpt.Zoom)
	ry := SCREEN_HEIGHT / (g.Universe.Size.Y * g.EditOpt.Zoom)
	// log.Println("R X:", rx)
	// log.Println("R Y:", ry)
	for _, obj := range g.Universe.Objects {
		square = ebiten.NewImage(1, 1)
		square.Fill(color.White)

		opts := &ebiten.DrawImageOptions{}
		px := obj.Pos.X*rx - g.EditOpt.Offset.X
		py := obj.Pos.Y*ry - g.EditOpt.Offset.Y
		// log.Println("Pos X:", px)
		// log.Println("Pos Y:", py)
		opts.GeoM.Translate(px, py)
		screen.DrawImage(square, opts)
	}

	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Zoom: %v", g.EditOpt.Zoom), 0, 0)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Offset: %vx  %vy", g.EditOpt.Offset.X, g.EditOpt.Offset.Y), 0, 15)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Universe size: %vx  %vy", g.Universe.Size.X, g.Universe.Size.Y), 0, 30)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Amount of objects: %v", len(g.Universe.Objects)), 0, 45)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Gravitational constant: %v", g.Universe.Gconst), 0, 60)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return SCREEN_WIDTH, SCREEN_HEIGHT
}

func printObjects(u *Universe) {
	for {
		fmt.Println("\n###############")
		for _, obj := range u.Objects {
			fmt.Println(obj.Pos)
		}
		fmt.Println("###############")
	}
}

func main() {
	// objs := []*Object{
	// 	NewObject(Coordinates2D{10, 10}, 10000000000, 1),
	// 	NewObject(Coordinates2D{50, 50}, 10, 1),
	// }
	// universe := NewUniverse(Coordinates2D{SCREEN_WIDTH, SCREEN_HEIGHT}, objs...)

	rand.Seed(time.Now().UnixNano())
	universe := NewRandomUniverse(
		Coordinates2D{SCREEN_WIDTH, SCREEN_HEIGHT},
		G,
		[2]float64{1, 1000},
		[2]float64{1, 1},
		10,
	)

	ebiten.SetWindowSize(SCREEN_WIDTH, SCREEN_HEIGHT)
	ebiten.SetWindowTitle("Gravity Simulator")

	g := &Game{
		Universe: universe,
		RandOpt: RandOpt{
			10,
			[2]float64{1, 1000},
			[2]float64{1, 1},
		},
		EditOpt: EditOpt{
			ObjectsDesloc: 10,
			GconstDesloc:  2,
			Zoom:          1,
			ZoomDesloc:    2,
			Offset:        Coordinates2D{0, 0},
			OffsetDesloc:  10,
		},
		keys: []ebiten.Key{},
	}
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
