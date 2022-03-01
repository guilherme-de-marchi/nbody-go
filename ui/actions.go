package ui

import (
	simul "github.com/Guilherme-De-Marchi/GravitySimulator/simulation"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type ActionFunc func(*Game)

var KeyMap = map[ebiten.Key]ActionFunc{
	ebiten.KeyW: MoveScreenUp,
	ebiten.KeyS: MoveScreenDown,
	ebiten.KeyA: MoveScreenLeft,
	ebiten.KeyD: MoveScreenRight,

	ebiten.KeyEscape: ShowPauseScreen,
	ebiten.Key1:      ShowDebug,
	ebiten.Key2:      ShowObject,
	ebiten.Key3:      ShowObjectName,
	ebiten.Key4:      ShowWinGravityGrad,
	ebiten.Key5:      ShowTotalGravityGrad,

	ebiten.KeyZ: SetZoom,
	ebiten.KeyR: NewRandomUniverse,
	ebiten.KeyG: SetGconst,
	ebiten.KeyO: SetObjects,
	ebiten.KeyE: SetGradExp,
}

// Key: W : Offset.Y -= Offset desloc
func MoveScreenUp(g *Game) {
	g.EditOpt.Offset.Y -= g.EditOpt.OffsetDesloc
}

// Key: S : Offset.Y += Offset desloc
func MoveScreenDown(g *Game) {
	g.EditOpt.Offset.Y += g.EditOpt.OffsetDesloc
}

// Key: A : Offset.X -= Offset desloc
func MoveScreenLeft(g *Game) {
	g.EditOpt.Offset.X -= g.EditOpt.OffsetDesloc
}

// Key: D : Offset.X -= Offset desloc
func MoveScreenRight(g *Game) {
	g.EditOpt.Offset.X += g.EditOpt.OffsetDesloc
}

// Key: Escape : Show pause screen (on/off)
func ShowPauseScreen(g *Game) {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		g.EditOpt.ShowPauseScreen = !g.EditOpt.ShowPauseScreen
	}
}

// Key: 1 : Show debug (on/off)
func ShowDebug(g *Game) {
	if inpututil.IsKeyJustPressed(ebiten.Key1) {
		g.EditOpt.ShowDebug = !g.EditOpt.ShowDebug
	}
}

// Key: 2 : Show object (on/off)
func ShowObject(g *Game) {
	if inpututil.IsKeyJustPressed(ebiten.Key2) {
		g.EditOpt.ShowObject = !g.EditOpt.ShowObject
	}
}

// Key: 3 : Show object name (on/off)
func ShowObjectName(g *Game) {
	if inpututil.IsKeyJustPressed(ebiten.Key3) {
		g.EditOpt.ShowObjectName = !g.EditOpt.ShowObjectName
	}
}

// Key: 4 : Show gravity gradient of window (on/off)
func ShowWinGravityGrad(g *Game) {
	if inpututil.IsKeyJustPressed(ebiten.Key4) {
		g.EditOpt.ShowWinGravityGrad = !g.EditOpt.ShowWinGravityGrad
	}
}

// Key: 5 : Show total gravity gradient (on/off)
func ShowTotalGravityGrad(g *Game) {
	if inpututil.IsKeyJustPressed(ebiten.Key5) {
		g.EditOpt.ShowTotalGravityGrad = !g.EditOpt.ShowTotalGravityGrad
	}
}

/*
Keys:
	Z + ArrowUp : Zoom /= Zoom desloc.
	Z + ArrowDown : Zoom *= Zoom desloc.
*/
func SetZoom(g *Game) {
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
		g.EditOpt.Zoom /= g.EditOpt.ZoomDesloc
	} else if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
		g.EditOpt.Zoom *= g.EditOpt.ZoomDesloc
	}
}

// Key: R : Generates e new random universe
func NewRandomUniverse(g *Game) {
	g.Universe = simul.NewRandomUniverse(
		g.Universe.Size,
		g.Universe.Gconst,
		g.RandOpt.MassR,
		g.RandOpt.RadR,
		g.RandOpt.ObjectQtt,
	)
}

/*
Keys:
	G + ArrowUp : Gravitational constant *= gConst desloc.
	G + ArrowDown : Gravitational constant /= gConst desloc.
*/
func SetGconst(g *Game) {
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
		g.Universe.Gconst *= g.EditOpt.GconstDesloc
	} else if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
		g.Universe.Gconst /= g.EditOpt.GconstDesloc
	}
}

/*
Keys:
	O + ArrowUp : Add n random objects to the universe.
	O + ArrowDown : Remove last n objects added to the universe.

	n = Objects desloc.
*/
func SetObjects(g *Game) {
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
		objs := simul.GetRandomObjects(
			g.Universe.Size,
			g.RandOpt.MassR,
			g.RandOpt.RadR,
			g.EditOpt.ObjectsDesloc,
		)
		g.Universe.AddObjects(objs...)
	} else if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
		if len(g.Universe.Objects) == 0 {
			return
		}

		if r := len(g.Universe.Objects) - g.EditOpt.ObjectsDesloc; r <= 0 {
			g.Universe.Objects = []*simul.Object{}
		} else {
			g.Universe.Objects = g.Universe.Objects[:r]
		}
	}
}

/*
Keys:
	E + ArrowUp : Gradient exp += Gradient exp desloc.
	E + ArrowDown : Gradient exp -= Gradient exp desloc.
*/
func SetGradExp(g *Game) {
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
		g.EditOpt.GradExp += g.EditOpt.GradExpDesloc
	} else if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
		g.EditOpt.GradExp -= g.EditOpt.GradExpDesloc
	}
}
