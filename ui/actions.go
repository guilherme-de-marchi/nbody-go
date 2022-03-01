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
	ebiten.KeyN: SwitchShowName,
	ebiten.KeyZ: SetZoom,
	ebiten.KeyR: NewRandomUniverse,
	ebiten.KeyG: SetGconst,
	ebiten.KeyO: SetObjects,
}

// Key: W : Offset.Y -= Offset desloc
func MoveScreenUp(g *Game) {
	g.Offset.Y -= g.OffsetDesloc
}

// Key: S : Offset.Y += Offset desloc
func MoveScreenDown(g *Game) {
	g.Offset.Y += g.OffsetDesloc
}

// Key: A : Offset.X -= Offset desloc
func MoveScreenLeft(g *Game) {
	g.Offset.X -= g.OffsetDesloc
}

// Key: D : Offset.X -= Offset desloc
func MoveScreenRight(g *Game) {
	g.Offset.X += g.OffsetDesloc
}

func SwitchShowName(g *Game) {
	if inpututil.IsKeyJustPressed(ebiten.KeyN) {
		g.EditOpt.ShowName = !g.EditOpt.ShowName
	}
}

/*
Keys:
	Z + ArrowUp : Zoom /= Zoom desloc.
	Z + ArrowDown : Zoom *= Zoom desloc.
*/
func SetZoom(g *Game) {
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
		g.Zoom /= g.ZoomDesloc
	} else if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
		g.Zoom *= g.ZoomDesloc
	}
}

// Key: R : Generates e new random universe
func NewRandomUniverse(g *Game) {
	g.Universe = simul.NewRandomUniverse(
		g.Universe.Size,
		g.Universe.Gconst,
		g.MassR,
		g.RadR,
		g.ObjectQtt,
	)
}

/*
Keys:
	G + ArrowUp : Gravitational constant *= gConst desloc.
	G + ArrowDown : Gravitational constant /= gConst desloc.
*/
func SetGconst(g *Game) {
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
		g.Universe.Gconst *= g.GconstDesloc
	} else if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
		g.Universe.Gconst /= g.GconstDesloc
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
			g.MassR,
			g.RadR,
			g.ObjectsDesloc,
		)
		g.Universe.AddObjects(objs...)
	} else if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
		if len(g.Universe.Objects) == 0 {
			return
		}

		if r := len(g.Universe.Objects) - g.ObjectsDesloc; r <= 0 {
			g.Universe.Objects = []*simul.Object{}
		} else {
			g.Universe.Objects = g.Universe.Objects[:r]
		}
	}
}
