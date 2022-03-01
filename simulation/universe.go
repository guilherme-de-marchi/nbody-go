package simulation

import (
	"image/color"

	"github.com/Guilherme-De-Marchi/GravitySimulator/util"
)

type Universe struct {
	Size    Coordinates2D `json:"size,omitempty"`
	Gconst  float64       `json:"gravitational_const,omitempty"`
	Objects []*Object
}

func NewUniverse(size Coordinates2D, gConst float64, objs ...*Object) *Universe {
	return &Universe{
		Size:    size,
		Gconst:  gConst,
		Objects: objs,
	}
}

func NewRandomUniverse(size Coordinates2D, gConst float64, massR, radR [2]float64, qtt int) *Universe {
	objs := GetRandomObjects(
		size,
		massR,
		radR,
		qtt,
	)
	return NewUniverse(size, gConst, objs...)
}

func (u *Universe) AddObjects(obj ...*Object) {
	u.Objects = append(u.Objects, obj...)
}

func (u *Universe) ApplyGravity() {
	for _, obj := range u.Objects {
		for _, tar := range u.Objects {
			if tar == obj {
				continue
			}
			f := obj.GetGravitationalForce(tar, u.Gconst)
			// log.Println("resulting force:", f, "\n")
			obj.ApplyForce(f, tar)
		}
	}
}

// Fix this
/*
Returns the gradient matrix and the highest value found
*/
func (u *Universe) GetViewGravityGradient(size Coordinates2D, r, offset [2]float64) ([][]float64, float64) {
	obj := NewObject("-", color.RGBA{}, Coordinates2D{}, 1, 0) // irrelevant object
	var totalf float64

	gradient := make([][]float64, int(size.Y))
	var high float64
	for i := range gradient {
		gradient[i] = make([]float64, int(size.X))
		for j := range gradient[i] {
			totalf = 0
			obj.Pos.X, obj.Pos.Y = util.PxToPos([2]float64{float64(j), float64(i)}, r, offset)
			for _, tar := range u.Objects {
				totalf += obj.GetGravitationalForce(tar, u.Gconst).Magnitude
			}
			gradient[i][j] = totalf

			if totalf > high {
				high = totalf
			}
		}
	}

	return gradient, high
}

func (u *Universe) GetTotalGravityGradient(size Coordinates2D) ([][]float64, float64) {
	obj := NewObject("-", color.RGBA{}, Coordinates2D{}, 1, 0) // irrelevant object
	var totalf float64

	gradient := make([][]float64, int(size.Y))
	var high float64
	for i := range gradient {
		gradient[i] = make([]float64, int(size.X))
		for j := range gradient[i] {
			totalf = 0
			obj.Pos.X, obj.Pos.Y = float64(j), float64(i)
			for _, tar := range u.Objects {
				totalf += obj.GetGravitationalForce(tar, u.Gconst).Magnitude
			}
			gradient[i][j] = totalf

			if totalf > high {
				high = totalf
			}
		}
	}

	return gradient, high
}
