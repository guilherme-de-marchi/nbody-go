package ui

import (
	"fmt"
	"image"
	"image/color"
	"log"

	simul "github.com/Guilherme-De-Marchi/nbody-go/simulation"
	"github.com/Guilherme-De-Marchi/nbody-go/util"
	"github.com/fogleman/gg"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	SCREEN_WIDTH  = 300
	SCREEN_HEIGHT = 300
)

var circle *ebiten.Image

type Game simul.Simulation

func (g *Game) Init() error {
	g.WinGradientImage = image.NewRGBA(image.Rect(0, 0, SCREEN_WIDTH, SCREEN_HEIGHT))
	g.TotalGradientImage = image.NewRGBA(image.Rect(0, 0, int(g.Universe.Size.X), int(g.Universe.Size.Y)))

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

	if g.EditOpt.ShowTotalGravityGrad {
		g.UpdateTotalGravityGrad()
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
		g.EditOpt.GradExp,
		[2]float64{SCREEN_WIDTH, SCREEN_HEIGHT},
		[2]float64{rx, ry},
		[2]float64{g.EditOpt.Offset.X, g.EditOpt.Offset.Y},
	)
	// log.Println(high)

	var c color.RGBA
	rc := high / util.MAX_RGB_INT
	for i := 0; i < SCREEN_HEIGHT; i++ {
		for j := 0; j < SCREEN_WIDTH; j++ {
			f := gradient[i][j]
			c = util.IntToRgb(int(f / rc))
			g.WinGradientImage.Set(j, i, c)
		}
	}
}

func (g *Game) UpdateTotalGravityGrad() {
	gradient, high := g.Universe.GetTotalGravityGradient(
		1,
		g.EditOpt.GradExp,
	)
	log.Println(high)

	var c color.RGBA
	rc := high / util.MAX_RGB_INT
	for i := 0; i < int(g.Universe.Size.Y); i++ {
		for j := 0; j < int(g.Universe.Size.X); j++ {
			f := gradient[i][j]
			c = util.IntToRgb(int(f / rc))
			g.TotalGradientImage.Set(j, i, c)
		}
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Clear()
	screen.Fill(color.Black)

	if g.EditOpt.ShowPauseScreen {
		g.DrawPauseScreen(screen)
		return
	}

	if g.EditOpt.ShowWinGravityGrad {
		g.DrawWinGravGrad(screen)
	}

	if g.EditOpt.ShowTotalGravityGrad {
		g.DrawTotalGravGrad(screen)
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
	ebitenutil.DebugPrintAt(screen, "E + ArrowUp : Increases Gradient Exp", 0, 270)
	ebitenutil.DebugPrintAt(screen, "E + ArrowDown : Decreases Gradient Exp", 0, 285)
}

func (g *Game) DrawDebug(screen *ebiten.Image) {
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("TPS: %0.2f", ebiten.CurrentTPS()), 0, 0)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Zoom: %v", 1/g.EditOpt.Zoom), 0, 15)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Offset: %vx  %vy", g.EditOpt.Offset.X, g.EditOpt.Offset.Y), 0, 30)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Universe size: %vx  %vy", g.Universe.Size.X, g.Universe.Size.Y), 0, 45)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Amount of objects: %v", len(g.Universe.Objects)), 0, 60)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Gravitational constant: %v", g.Universe.Gconst), 0, 75)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Gradient exp: %v", g.EditOpt.GradExp), 0, 90)

	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Show objects: %v", g.EditOpt.ShowObject), 0, 105)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Show objects name: %v", g.EditOpt.ShowObjectName), 0, 120)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Show gravitational gradient: %v", g.EditOpt.ShowWinGravityGrad), 0, 135)
}

func (g *Game) DrawObject(screen *ebiten.Image) {
	rx := SCREEN_WIDTH / (g.Universe.Size.X * g.EditOpt.Zoom)
	ry := SCREEN_HEIGHT / (g.Universe.Size.Y * g.EditOpt.Zoom)

	var ctx *gg.Context
	for _, obj := range g.Universe.Objects {
		// Fix this
		px, py := util.PosToPx(
			[2]float64{obj.Pos.X - obj.Radius, obj.Pos.Y - obj.Radius},
			[2]float64{rx, ry},
			[2]float64{g.EditOpt.Offset.X, g.EditOpt.Offset.Y},
		)

		lx := int(obj.Radius * 2 * rx)
		ly := int(obj.Radius * 2 * ry)
		if lx <= 0 || ly <= 0 {
			continue
		}

		ctx = gg.NewContext(lx, ly)
		ctx.DrawCircle(float64(lx/2), float64(ly/2), float64(lx/2))
		ctx.SetColor(obj.Color)
		ctx.Fill()

		circle = ebiten.NewImageFromImage(ctx.Image())
		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Translate(px, py)
		screen.DrawImage(circle, opts)

		if g.EditOpt.ShowObjectName {
			g.DrawObjectName(screen, obj, px, py)
		}
	}
}

func (g *Game) DrawObjectName(screen *ebiten.Image, obj *simul.Object, px, py float64) {
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Name: %v", obj.Name), int(px), int(py-15))
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Mass: %v", obj.Mass), int(px), int(py-30))
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Radius: %v", obj.Radius), int(px), int(py-45))
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Velocity: %v", obj.Vel), int(px), int(py-60))
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Acceleration: %v", obj.Accel), int(px), int(py-75))
}

func (g *Game) DrawWinGravGrad(screen *ebiten.Image) {
	screen.ReplacePixels(g.WinGradientImage.Pix)
}

func (g *Game) DrawTotalGravGrad(screen *ebiten.Image) {
	screen.ReplacePixels(g.TotalGradientImage.Pix)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return SCREEN_WIDTH, SCREEN_HEIGHT
}
