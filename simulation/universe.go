package simulation

import (
	"image/color"
	"math"

	"github.com/Guilherme-De-Marchi/nbody-go/util"
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
func (u *Universe) GetViewGravityGradient(exp float64, size, r, offset [2]float64) ([][]float64, float64) {
	obj := NewObject("-", color.RGBA{}, Coordinates2D{}, 10000, 0) // irrelevant object
	var totalf float64

	gradient := make([][]float64, int(size[1]))
	var high float64
	for i := range gradient {
		gradient[i] = make([]float64, int(size[0]))
		for j := range gradient[i] {
			totalf = 0
			obj.Pos.X, obj.Pos.Y = util.PxToPos([2]float64{float64(j), float64(i)}, r, offset)
			for _, tar := range u.Objects {
				totalf += obj.GetGravitationalForce(tar, u.Gconst).Magnitude
			}
			gradient[i][j] = math.Pow(totalf, exp)

			if totalf > high {
				high = totalf
			}
		}
	}

	return gradient, high
}

func (u *Universe) GetTotalGravityGradient(step, exp float64) ([][]float64, float64) {
	obj := NewObject("-", color.RGBA{}, Coordinates2D{}, 1, 0) // irrelevant object
	var totalf float64

	tX := int(u.Size.X / step)
	tY := int(u.Size.Y / step)
	gradient := make([][]float64, tY)
	var high float64
	for i := 0; i < tY; i++ {
		gradient[i] = make([]float64, tX)
		for j := 0; j < tX; j++ {
			totalf = 0
			for _, tar := range u.Objects {
				totalf += obj.GetGravitationalForce(tar, u.Gconst).Magnitude
			}
			gradient[i][j] = math.Pow(totalf, exp)

			if totalf > high {
				high = totalf
			}
			obj.Pos.X += step
		}
		obj.Pos.Y += step
	}

	return gradient, high
}
