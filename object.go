package main

type Object struct {
	Pos                 Coordinates2D
	Mass, Accel, Radius float64
	Vel, Momentum       Vector2
}

func NewObject(pos Coordinates2D, mass, radius float64) *Object {
	return &Object{
		Pos:    pos,
		Mass:   mass,
		Radius: radius,
	}
}

func (obj *Object) GetDistance(tar *Object) float64 {
	return CalcDistance(obj.Pos.X, tar.Pos.X, obj.Pos.Y, tar.Pos.Y)
}

func (obj *Object) GetAbsDistance(tar *Object) float64 {
	return CalcAbsDistance(obj.Pos.X, tar.Pos.X, obj.Pos.Y, tar.Pos.Y)
}

func (obj *Object) GetGravitationalForce(tar *Object, gConst float64) Vector2 {
	d := obj.GetDistance(tar)
	return Vector2{
		Direction: obj.GetVectorDirection(tar),
		Magnitude: CalcGravitationalForce(obj.Mass, tar.Mass, d, gConst),
	}
}

// func (obj *Object) GetMomentum() Vector2 {
// 	return CalcMomentum(obj.Accel, obj.Mass)
// }

/*
Reduces the distance between two objects to 1
and returns the proportional x and y changes
*/
func (obj *Object) GetVectorDirection(tar *Object) Coordinates2D {
	d := obj.GetDistance(tar)
	lx := tar.Pos.X - obj.Pos.X
	ly := tar.Pos.Y - obj.Pos.Y

	// log.Println("Direction X:", CalcProportionalLeg(d, 1, lx))
	// log.Println("Direction Y:", CalcProportionalLeg(d, 1, ly))
	return Coordinates2D{
		X: CalcProportionalLeg(d, 1, lx),
		Y: CalcProportionalLeg(d, 1, ly),
	}
}

/*
Changes the object velocity and acceleration
based on the force applied
*/
func (obj *Object) ApplyForce(f Vector2, tar *Object) {
	obj.Vel.Direction = f.Direction
	obj.Accel = obj.GetResultingAcceleration(f, tar)
	obj.Vel.Magnitude += obj.Accel

	obj.SetPos(obj.GetResultingPos())
}

func (obj *Object) GetResultingAcceleration(f Vector2, tar *Object) float64 {
	if obj.GetDistance(tar) <= obj.Radius+tar.Radius {
		return 0
	}
	return CalcAcceleration(f, obj.Mass)
}

func (obj *Object) GetResultingPos() Coordinates2D {
	return CalcResultingPosition(obj.Pos, obj.Vel, obj.Accel)
}

func (obj *Object) SetPos(pos Coordinates2D) {
	obj.Pos = pos
}

func GetRandomObjects(posxR, posyR, massR, radR [2]float64, qtt int) []*Object {
	objs := make([]*Object, qtt)
	for i := range objs {
		objs[i] = NewObject(
			Coordinates2D{RandFloatRange(posxR[0], posxR[1]), RandFloatRange(posyR[0], posyR[1])},
			RandFloatRange(massR[0], massR[1]),
			RandFloatRange(massR[0], massR[1]),
		)
	}
	return objs
}
