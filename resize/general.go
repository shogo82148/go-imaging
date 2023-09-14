package resize

import (
	"github.com/shogo82148/float16"
	"github.com/shogo82148/go-imaging/fp16"
	"github.com/shogo82148/go-imaging/fp16/fp16color"
	"github.com/shogo82148/go-imaging/internal/parallels"
)

// General resizes the image using General (Bicubic B:1, C:0) interpolation.
func General(dst, src *fp16.NRGBAh) {
	dstBounds := dst.Bounds()
	dstDx := dstBounds.Dx()
	//dstDy := dstBounds.Dy()
	srcBounds := src.Bounds()
	srcDx := srcBounds.Dx()
	srcDy := srcBounds.Dy()

	//tmp := fp16.NewNRGBAh(image.Rect(0, 0, dstDx, srcDy))

	// resize horizontally
	parallels.Parallel(0, srcDy, func(y int) {
		for x := 0; x < dstDx; x++ {
			var c fp16color.NRGBAh
			srcX, remX := mulDiv(x, srcDx-1, dstDx-1)
			d := float64(remX) / float64(dstDx-1)
			c0 := nrgbhAt(src, srcBounds.Min.X+srcX-1, srcBounds.Min.Y+y)
			c1 := nrgbhAt(src, srcBounds.Min.X+srcX+0, srcBounds.Min.Y+y)
			c2 := nrgbhAt(src, srcBounds.Min.X+srcX+1, srcBounds.Min.Y+y)
			c3 := nrgbhAt(src, srcBounds.Min.X+srcX+2, srcBounds.Min.Y+y)

			// https://qiita.com/yoya/items/f167b2598fec98679422
			// https://legacy.imagemagick.org/Usage/filter/#mitchell
			// https://www.cs.utexas.edu/~fussell/courses/cs384g-fall2013/lectures/mitchell/Mitchell.pdf
			a0 := -1.0 / 6.0 * (d - 1) * (d - 1) * (d - 1)
			a1 := (-1.0/2.0*d-1)*d*d + 2.0/3.0
			a2 := d*(d*(1.0/2.0*d-5.0/2.0)+7.0/2.0) - 5.0/6.0
			a3 := 1.0 / 6.0 * d * d * d

			c.R = general(c0.R, c1.R, c2.R, c3.R, a0, a1, a2, a3)
			c.G = general(c0.G, c1.G, c2.G, c3.G, a0, a1, a2, a3)
			c.B = general(c0.B, c1.B, c2.B, c3.B, a0, a1, a2, a3)
			c.A = general(c0.A, c1.A, c2.A, c3.A, a0, a1, a2, a3)
			dst.SetNRGBAh(x+dstBounds.Min.X, y+dstBounds.Min.Y, c)
		}
	})
}

// nrgbhAt is similar to src.NRGBAhAt(x, y), but it returns the color at the
// nearest point if the point is out of bounds.
func nrgbhAt(img *fp16.NRGBAh, x, y int) fp16color.NRGBAh {
	bounds := img.Bounds()
	x = max(bounds.Min.X, min(bounds.Max.X, x))
	y = max(bounds.Min.Y, min(bounds.Max.Y, y))
	return img.NRGBAhAt(x, y)
}

func general(c0, c1, c2, c3 float16.Float16, a0, a1, a2, a3 float64) float16.Float16 {
	return float16.FromFloat64(c0.Float64()*a0 + c1.Float64()*a1 + c2.Float64()*a2 + c3.Float64()*a3)
}
