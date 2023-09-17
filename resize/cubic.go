package resize

import (
	"math"

	"github.com/shogo82148/float16"
	"github.com/shogo82148/go-imaging/fp16"
	"github.com/shogo82148/go-imaging/fp16/fp16color"
)

// https://qiita.com/yoya/items/f167b2598fec98679422
func cubicBCcoefficient(b, c float64) []float64 {
	p := 2 - 1.5*b - c
	q := -3 + 2*b + c
	r := 0.0
	s := 1 - (1.0/3)*b
	t := -(1.0/6)*b - c
	u := b + 5.0*c
	v := -2*b - 8*c
	w := (4.0/3)*b + 4*c
	return []float64{p, q, r, s, t, u, v, w}
}

func cubicBC(x float64, coeff []float64) float64 {
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

func general(c0, c1, c2, c3 float16.Float16, d float64, coeff []float64) float16.Float16 {
	a0 := cubicBC(1+d, coeff)
	a1 := cubicBC(d, coeff)
	a2 := cubicBC(1-d, coeff)
	a3 := cubicBC(2-d, coeff)
	return float16.FromFloat64(c0.Float64()*a0 + c1.Float64()*a1 + c2.Float64()*a2 + c3.Float64()*a3)
}

// nrgbhAt is similar to src.NRGBAhAt(x, y), but it returns the color at the
// nearest point if the point is out of bounds.
func nrgbhAt(img *fp16.NRGBAh, x, y int) fp16color.NRGBAh {
	bounds := img.Bounds()
	x = max(bounds.Min.X, min(bounds.Max.X-1, x))
	y = max(bounds.Min.Y, min(bounds.Max.Y-1, y))
	return img.NRGBAhAt(x, y)
}
