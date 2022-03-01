package ui

import (
	"fmt"
	"image/color"

	simul "github.com/Guilherme-De-Marchi/GravitySimulator/simulation"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	SCREEN_WIDTH  = 800
	SCREEN_HEIGHT = 800
)

var square *ebiten.Image

type Game simul.Simulation

func (g *Game) Init() error {
	ebiten.SetWindowSize(SCREEN_WIDTH, SCREEN_HEIGHT)
	ebiten.SetWindowTitle("Gravity Simulator")

	if err := ebiten.RunGame(g); err != nil {
		return err
	}

	return nil
}

func (g *Game) Update() error {
	g.Keys = inpututil.AppendPressedKeys(g.Keys[:0])

	for _, k := range g.Keys {
		if f, ok := KeyMap[k]; ok {
			f(g)
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

		if g.EditOpt.ShowName {
			ebitenutil.DebugPrintAt(screen, obj.Name, int(px+2), int(py-12))
		}
	}

	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("TPS: %0.2f", ebiten.CurrentFPS()), 0, 0)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Zoom: %v", 1/g.EditOpt.Zoom), 0, 15)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Offset: %vx  %vy", g.EditOpt.Offset.X, g.EditOpt.Offset.Y), 0, 30)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Show objects name: %v", g.EditOpt.ShowName), 0, 45)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Universe size: %vx  %vy", g.Universe.Size.X, g.Universe.Size.Y), 0, 60)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Amount of objects: %v", len(g.Universe.Objects)), 0, 75)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Gravitational constant: %v", g.Universe.Gconst), 0, 90)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return SCREEN_WIDTH, SCREEN_HEIGHT
}
