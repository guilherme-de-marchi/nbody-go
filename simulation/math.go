package simulation

import (
	"math"
)

// Gravitational Constant
var G = 6.674 * math.Pow(10, -11)

type Coordinates2D struct {
	X float64 `json:"x,omitempty"`
	Y float64 `json:"y,omitempty"`
}

type Vector2 struct {
	Direction Coordinates2D
	Magnitude float64
}

/*
Using the equation for universal gravitation.
F = G*((m1*m2)/r**2)
m1 and m2 are the masses of the objects;
r is the distance between the centers of their masses;
k is the gravitational constant. (may use var G from this pkg);
*/
func CalcGravitationalForce(m1, m2, r, k float64) float64 {
	// log.Println("CalcGravitationalForce:", k*((m1*m2)/math.Pow(r, 2)))
	return k * ((m1 * m2) / math.Pow(r, 2))
}

/*
Using Pythagorean theorem.
h = √(x2-x1^2 + y2-y1^2)
*/
func CalcDistance(x1, x2, y1, y2 float64) float64 {
	// log.Println("CalcDistance:", math.Sqrt(math.Pow(x2-x1, 2)+math.Pow(y2-y1, 2)))
	return math.Sqrt(math.Pow(x2-x1, 2) + math.Pow(y2-y1, 2))
}

/*
Using Pythagorean theorem.
h = √(|x2-x1|^2 + |y2-y1|^2)
*/
func CalcAbsDistance(x1, x2, y1, y2 float64) float64 {
	// log.Println("CalcAbsDistance:", math.Sqrt(math.Pow(math.Abs(x2-x1), 2)+math.Pow(math.Abs(y2-y1), 2)))
	return math.Sqrt(math.Pow(math.Abs(x2-x1), 2) + math.Pow(math.Abs(y2-y1), 2))
}

/*
Using second Newton's law.
a=F/m
F is the force applied on the object;
m is the mass of the object;
a is the acceletarion of the object;
*/
func CalcAcceleration(f Vector2, m float64) float64 {
	// log.Println("Force causing Acceleration:", f.Magnitude)
	// log.Println("Mass:", m)
	// log.Println("Acceleration:", f.Magnitude/m)
	return f.Magnitude / m
}

/*
p=a*m
m is the mass of the object;
a is the acceletarion of the object;
*/
func CalcMomentum(a Vector2, m float64) Vector2 {
	return Vector2{
		Direction: a.Direction,
		Magnitude: a.Magnitude * m,
	}
}

/*
c2=h2*c1/h1
*/
func CalcProportionalLeg(h1, h2, c1 float64) float64 {
	// log.Println("CalcProportionalLeg:", h2*c1/h1)
	return h2 * c1 / h1
}

func CalcResultingPosition(pos Coordinates2D, vel Vector2, accel float64) Coordinates2D {
	// log.Println("Pos X:", pos.X)
	// log.Println("Pos Y:", pos.Y)
	// log.Println("Direction X:", vel.Direction.X)
	// log.Println("Direction Y:", vel.Direction.Y)
	// log.Println("Resulting Pos X:", pos.X+vel.Direction.X*vel.Magnitude)
	// log.Println("Resulting Pos Y:", pos.Y+vel.Direction.Y*vel.Magnitude)
	// log.Println("Velocity magnitude:", vel.Magnitude)
	return Coordinates2D{
		X: pos.X + vel.Direction.X*vel.Magnitude,
		Y: pos.Y + vel.Direction.Y*vel.Magnitude,
	}
}

func CalcDensity(m, r float64) float64 {
	return m / (math.Pi * math.Pow(r, 2))
}
