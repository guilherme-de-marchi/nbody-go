package ui

import (
	"fmt"
	"image"
	"image/color"
	"log"

	simul "github.com/Guilherme-De-Marchi/GravitySimulator/simulation"
	"github.com/Guilherme-De-Marchi/GravitySimulator/util"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	SCREEN_WIDTH  = 300
	SCREEN_HEIGHT = 300
)

var square *ebiten.Image

type Game simul.Simulation

func (g *Game) Init() error {
	g.GradientImage = image.NewRGBA(image.Rect(0, 0, int(g.Universe.Size.X), int(g.Universe.Size.Y)))

	ebiten.SetWindowSize(SCREEN_WIDTH, SCREEN_HEIGHT)
	ebiten.SetWindowTitle("Gravity Simulator")

	log.Println("[GAME] PRESS 'ESCAPE' TO SEE THE CONTROLS")
	if err := ebiten.RunGame(g); err != nil {
		return err
	}

	return nil
}

func (g *Game) Update() error {
	if g.EditOpt.ShowWinGravityGrad {
		g.UpdateWinGravityGrad()
	}

	g.Keys = inpututil.AppendPressedKeys(g.Keys[:0])

	for _, k := range g.Keys {
		if f, ok := KeyMap[k]; ok {
			f(g)
		}
	}

	if !g.EditOpt.ShowPauseScreen {
		g.Universe.ApplyGravity()
	}
	return nil
}

func (g *Game) UpdateWinGravityGrad() {
	rx := SCREEN_WIDTH / (g.Universe.Size.X * g.EditOpt.Zoom)
	ry := SCREEN_HEIGHT / (g.Universe.Size.Y * g.EditOpt.Zoom)
	gradient, high := g.Universe.GetViewGravityGradient(
		g.Universe.Size,
		[2]float64{rx, ry},
		[2]float64{g.EditOpt.Offset.X, g.EditOpt.Offset.Y},
	)

	var c color.RGBA
	rc := high / util.MAX_RGB_INT
	for i := 0; i < int(g.Universe.Size.Y); i++ {
		for j := 0; j < int(g.Universe.Size.X); j++ {
			f := gradient[i][j]
			c = util.IntToRgb(int(f / rc))
			g.GradientImage.Set(j, i, c)
		}
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Clear()
	screen.Fill(color.White)

	if g.EditOpt.ShowPauseScreen {
		g.DrawPauseScreen(screen)
		return
	}

	if g.EditOpt.ShowWinGravityGrad {
		g.DrawWinGravGrad(screen)
	}

	if g.EditOpt.ShowObject {
		g.DrawObject(screen)
	}

	if g.EditOpt.ShowDebug {
		g.DrawDebug(screen)
	}
}

func (g *Game) DrawPauseScreen(screen *ebiten.Image) {
	ebitenutil.DebugPrintAt(screen, "W : Move Up", 0, 0)
	ebitenutil.DebugPrintAt(screen, "A : Move Left", 0, 15)
	ebitenutil.DebugPrintAt(screen, "S : Move Down", 0, 30)
	ebitenutil.DebugPrintAt(screen, "D : Move Right", 0, 45)

	ebitenutil.DebugPrintAt(screen, "1 : Show Debug informations", 0, 75)
	ebitenutil.DebugPrintAt(screen, "2 : Show Objects", 0, 90)
	ebitenutil.DebugPrintAt(screen, "3 : Show Objects name", 0, 105)
	ebitenutil.DebugPrintAt(screen, "4 : Show Gravity Gradient (objects on windows) [DROPS TPS]", 0, 120)
	ebitenutil.DebugPrintAt(screen, "Escape : Show Pause Screen", 0, 135)

	ebitenutil.DebugPrintAt(screen, "R : Generate a New Random Universe", 0, 165)
	ebitenutil.DebugPrintAt(screen, "Z + ArrowUp : Increases Zoom", 0, 180)
	ebitenutil.DebugPrintAt(screen, "Z + ArrowDown : Decreases Zoom", 0, 195)
	ebitenutil.DebugPrintAt(screen, "G + ArrowUp : Increases Gravitational Constant", 0, 210)
	ebitenutil.DebugPrintAt(screen, "G + ArrowDown : Decreases Gravitational Constant", 0, 225)
	ebitenutil.DebugPrintAt(screen, "O + ArrowUp : Add N Objects", 0, 240)
	ebitenutil.DebugPrintAt(screen, "O + ArrowDown : Remove N Objects", 0, 255)
}

func (g *Game) DrawDebug(screen *ebiten.Image) {
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("TPS: %0.2f", ebiten.CurrentTPS()), 0, 0)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Zoom: %v", 1/g.EditOpt.Zoom), 0, 15)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Offset: %vx  %vy", g.EditOpt.Offset.X, g.EditOpt.Offset.Y), 0, 30)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Universe size: %vx  %vy", g.Universe.Size.X, g.Universe.Size.Y), 0, 45)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Amount of objects: %v", len(g.Universe.Objects)), 0, 60)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Gravitational constant: %v", g.Universe.Gconst), 0, 75)

	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Show objects: %v", g.EditOpt.ShowObject), 0, 105)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Show objects name: %v", g.EditOpt.ShowObjectName), 0, 120)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Show gravitational gradient: %v", g.EditOpt.ShowWinGravityGrad), 0, 135)
}

func (g *Game) DrawObject(screen *ebiten.Image) {
	rx := SCREEN_WIDTH / (g.Universe.Size.X * g.EditOpt.Zoom)
	ry := SCREEN_HEIGHT / (g.Universe.Size.Y * g.EditOpt.Zoom)

	for _, obj := range g.Universe.Objects {
		square = ebiten.NewImage(1, 1)
		square.Fill(obj.Color)

		opts := &ebiten.DrawImageOptions{}
		// Fix this
		px, py := util.PosToPx([2]float64{obj.Pos.X, obj.Pos.Y}, [2]float64{rx, ry}, [2]float64{g.EditOpt.Offset.X, g.EditOpt.Offset.Y})

		opts.GeoM.Translate(px, py)
		screen.DrawImage(square, opts)

		if g.EditOpt.ShowObjectName {
			g.DrawObjectName(screen, obj, px, py)
		}
	}
}

func (g *Game) DrawObjectName(screen *ebiten.Image, obj *simul.Object, px, py float64) {
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Name: %v", obj.Name), int(px+10), int(py-15))
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Mass: %v", obj.Mass), int(px+10), int(py-30))
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Radius: %v", obj.Radius), int(px+10), int(py-45))
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Velocity: %v", obj.Vel), int(px+10), int(py-60))
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Acceleration: %v", obj.Accel), int(px+10), int(py-75))
}

func (g *Game) DrawWinGravGrad(screen *ebiten.Image) {
	screen.ReplacePixels(g.GradientImage.Pix)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return SCREEN_WIDTH, SCREEN_HEIGHT
}
