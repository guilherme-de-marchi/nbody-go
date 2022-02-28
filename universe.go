package main

type Universe struct {
	Size    Coordinates2D
	Gconst  float64
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
		[2]float64{0, size.X},
		[2]float64{0, size.Y},
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
