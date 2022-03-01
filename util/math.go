package util

func PosToPx(pos, r, offset [2]float64) (float64, float64) {
	return pos[0]*r[0] - offset[0], pos[1]*r[1] - offset[1]
}

func PxToPos(px, r, offset [2]float64) (float64, float64) {
	return (px[0] + offset[0]) / r[0], (px[1] + offset[1]) / r[1]
}
