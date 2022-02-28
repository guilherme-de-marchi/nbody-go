package main

import (
	"encoding/json"
	"fmt"
	"image/color"
	"log"
	"math/rand"
	"os"
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
	simulConf SimulConfig
	square    *ebiten.Image
)

type SimulConfig struct {
	GenerationType string   `json:"generation_type,omitempty"`
	Universe       Universe `json:"universe,omitempty"`
	RandOpt        RandOpt  `json:"random_options,omitempty"`
	EditOpt        EditOpt  `json:"edit_options,omitempty"`
}

type RandOpt struct {
	MassR     [2]float64 `json:"mass_range,omitempty"`
	RadR      [2]float64 `json:"object_radius_range,omitempty"`
	ObjectQtt int        `json:"object_quantity,omitempty"`
}

type EditOpt struct {
	ObjectsDesloc int           `json:"object_quantity_desloc,omitempty"`
	GconstDesloc  float64       `json:"gravitational_const_desloc,omitempty"`
	Zoom          float64       `json:"initial_zoom,omitempty"`
	ZoomDesloc    float64       `json:"zoom_desloc,omitempty"`
	Offset        Coordinates2D `json:"initial_offset,omitempty"`
	OffsetDesloc  float64       `json:"offset_desloc,omitempty"`
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
			g.RadR,
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
	confJ, err := os.ReadFile("./config.json")
	if err != nil {
		log.Fatal("[INTERNAL ERROR]: ", err)
	}
	json.Unmarshal(confJ, &simulConf)

	var universe *Universe
	if simulConf.GenerationType == "randomized" {
		rand.Seed(time.Now().UnixNano())
		universe = NewRandomUniverse(
			simulConf.Universe.Size,
			G,
			simulConf.RandOpt.MassR,
			simulConf.RandOpt.RadR,
			simulConf.RandOpt.ObjectQtt,
		)
	} else if simulConf.GenerationType == "prefab" {
		log.Fatal("Generation type 'prefab' not implemented")
	} else {
		log.Fatal("Invalid value for field 'generation_type'")
	}

	ebiten.SetWindowSize(SCREEN_WIDTH, SCREEN_HEIGHT)
	ebiten.SetWindowTitle("Gravity Simulator")

	g := &Game{
		Universe: universe,
		RandOpt:  simulConf.RandOpt,
		EditOpt:  simulConf.EditOpt,
		keys:     []ebiten.Key{},
	}
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
