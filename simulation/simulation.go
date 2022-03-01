package simulation

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

type RandOpt struct {
	MassR     [2]float64 `json:"mass_range,omitempty"`
	RadR      [2]float64 `json:"object_radius_range,omitempty"`
	ObjectQtt int        `json:"object_quantity,omitempty"`
}

type EditOpt struct {
	ShowPauseScreen      bool `json:"show_pause_screen,omitempty"`
	ShowDebug            bool `json:"show_debug,omitempty"`
	ShowObject           bool `json:"show_object,omitempty"`
	ShowObjectName       bool `json:"show_object_name,omitempty"`
	ShowWinGravityGrad   bool `json:"show_windown_gravity_gradient,omitempty"`
	ShowTotalGravityGrad bool `json:"show_total_gravity_gradient,omitempty"`

	ObjectsDesloc int           `json:"object_quantity_desloc,omitempty"`
	GconstDesloc  float64       `json:"gravitational_const_desloc,omitempty"`
	Zoom          float64       `json:"initial_zoom,omitempty"`
	ZoomDesloc    float64       `json:"zoom_desloc,omitempty"`
	Offset        Coordinates2D `json:"initial_offset,omitempty"`
	OffsetDesloc  float64       `json:"offset_desloc,omitempty"`
}

type Simulation struct {
	Universe *Universe

	RandOpt RandOpt
	EditOpt EditOpt

	GradientImage *image.RGBA
	Keys          []ebiten.Key
}

func NewSimulation(u *Universe, randOpt RandOpt, editOpt EditOpt) *Simulation {
	return &Simulation{
		Universe: u,
		RandOpt:  randOpt,
		EditOpt:  editOpt,
		Keys:     []ebiten.Key{},
	}
}
