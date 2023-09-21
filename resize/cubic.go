package resize

import (
	"math"
)

// https://qiita.com/yoya/items/f167b2598fec98679422
func cubicBCcoefficient(b, c float64) [8]float64 {
	p := 2 - 1.5*b - c
	q := -3 + 2*b + c
	r := 0.0
	s := 1 - (1.0/3)*b
	t := -(1.0/6)*b - c
	u := b + 5.0*c
	v := -2*b - 8*c
	w := (4.0/3)*b + 4*c
	return [8]float64{p, q, r, s, t, u, v, w}
}

func cubicBC(x float64, coeff [8]float64) float64 {
	var y float64
	p, q, r, s, t, u, v, w := coeff[0], coeff[1], coeff[2], coeff[3], coeff[4], coeff[5], coeff[6], coeff[7]
	x = math.Abs(x)
	if x < 1 {
		y = ((p*x+q)*x+r)*x + s
	} else if x < 2 {
		y = ((t*x+u)*x+v)*x + w
	}
	return y
}
