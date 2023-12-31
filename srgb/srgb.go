//go:generate go run ./internal/cmd/gen

// Package srgb handles [sRGB] colors.
//
// [sRGB]: https://en.wikipedia.org/wiki/SRGB
package srgb

import (
	"image"
	"image/color"
	"math"

	"github.com/shogo82148/float16"
	"github.com/shogo82148/go-imaging/fp16"
	"github.com/shogo82148/go-imaging/fp16/fp16color"
	"github.com/shogo82148/go-imaging/internal/parallels"
)

var one = float16.FromFloat64(1)

// DecodeTone decodes an sRGB color encoded image to a linear color image.
func DecodeTone(img image.Image) *fp16.NRGBAh {
	switch img := img.(type) {
	case *image.RGBA:
		return decodeToneRGBA(img)
	case *image.RGBA64:
		return decodeToneRGBA64(img)
	case *image.NRGBA:
		return decodeToneNRGBA(img)
	case *image.NRGBA64:
		return decodeToneNRGBA64(img)
	case *image.Gray:
		return decodeToneGray(img)
	case *image.Gray16:
		return decodeToneGray16(img)
	case *image.Alpha:
		return decodeToneAlpha(img)
	case *image.Alpha16:
		return decodeToneAlpha16(img)
	case *image.Paletted:
		return decodeTonePaletted(img)
	}
	return decodeTone(img)
}

func decodeToneRGBA(img *image.RGBA) *fp16.NRGBAh {
	bounds := img.Bounds()
	ret := fp16.NewNRGBAh(bounds)
	parallels.Parallel(bounds.Min.Y, bounds.Max.Y, func(y int) {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := img.RGBAAt(x, y)
			if c.A == 0 {
				fr := encodedToLinearTable[c.R]
				fg := encodedToLinearTable[c.G]
				fb := encodedToLinearTable[c.B]
				ret.SetNRGBAh(x, y, fp16color.NRGBAh{R: fr, G: fg, B: fb, A: 0})
			} else if c.A == 0xff {
				fr := encodedToLinearTable[c.R]
				fg := encodedToLinearTable[c.G]
				fb := encodedToLinearTable[c.B]
				ret.SetNRGBAh(x, y, fp16color.NRGBAh{R: fr, G: fg, B: fb, A: one})
			} else {
				fr := float64(c.R) / 0xff
				fg := float64(c.G) / 0xff
				fb := float64(c.B) / 0xff
				fa := float64(c.A) / 0xff
				fr /= fa
				fg /= fa
				fb /= fa
				fr = encodedToLinear(fr)
				fg = encodedToLinear(fg)
				fb = encodedToLinear(fb)
				ret.SetNRGBAh(x, y, fp16color.NewNRGBAh(fr, fg, fb, fa))
			}
		}
	})
	return ret
}

func decodeToneRGBA64(img *image.RGBA64) *fp16.NRGBAh {
	bounds := img.Bounds()
	ret := fp16.NewNRGBAh(bounds)
	parallels.Parallel(bounds.Min.Y, bounds.Max.Y, func(y int) {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := img.RGBA64At(x, y)
			if c.A == 0 {
				fr := encodedToLinearTable16[c.R]
				fg := encodedToLinearTable16[c.G]
				fb := encodedToLinearTable16[c.B]
				ret.SetNRGBAh(x, y, fp16color.NRGBAh{R: fr, G: fg, B: fb, A: 0})
			} else if c.A == 0xffff {
				fr := encodedToLinearTable16[c.R]
				fg := encodedToLinearTable16[c.G]
				fb := encodedToLinearTable16[c.B]
				ret.SetNRGBAh(x, y, fp16color.NRGBAh{R: fr, G: fg, B: fb, A: one})
			} else {
				fr := float64(c.R) / 0xffff
				fg := float64(c.G) / 0xffff
				fb := float64(c.B) / 0xffff
				fa := float64(c.A) / 0xffff
				fr /= fa
				fg /= fa
				fb /= fa
				fr = encodedToLinear(fr)
				fg = encodedToLinear(fg)
				fb = encodedToLinear(fb)
				ret.SetNRGBAh(x, y, fp16color.NewNRGBAh(fr, fg, fb, fa))
			}
		}
	})
	return ret
}

func decodeToneNRGBA(img *image.NRGBA) *fp16.NRGBAh {
	bounds := img.Bounds()
	ret := fp16.NewNRGBAh(bounds)
	parallels.Parallel(bounds.Min.Y, bounds.Max.Y, func(y int) {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := img.NRGBAAt(x, y)
			fr := encodedToLinearTable[c.R]
			fg := encodedToLinearTable[c.G]
			fb := encodedToLinearTable[c.B]
			fa := float16.FromFloat64(float64(c.A) / 0xff)
			ret.SetNRGBAh(x, y, fp16color.NRGBAh{R: fr, G: fg, B: fb, A: fa})
		}
	})
	return ret
}

func decodeToneNRGBA64(img *image.NRGBA64) *fp16.NRGBAh {
	bounds := img.Bounds()
	ret := fp16.NewNRGBAh(bounds)
	parallels.Parallel(bounds.Min.Y, bounds.Max.Y, func(y int) {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := img.NRGBA64At(x, y)
			fr := encodedToLinearTable16[c.R]
			fg := encodedToLinearTable16[c.G]
			fb := encodedToLinearTable16[c.B]
			fa := float16.FromFloat64(float64(c.A) / 0xffff)
			ret.SetNRGBAh(x, y, fp16color.NRGBAh{R: fr, G: fg, B: fb, A: fa})
		}
	})
	return ret
}

func decodeToneGray(img *image.Gray) *fp16.NRGBAh {
	bounds := img.Bounds()
	ret := fp16.NewNRGBAh(bounds)
	parallels.Parallel(bounds.Min.Y, bounds.Max.Y, func(y int) {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := img.GrayAt(x, y)
			fy := encodedToLinearTable[c.Y]
			ret.SetNRGBAh(x, y, fp16color.NRGBAh{R: fy, G: fy, B: fy, A: one})
		}
	})
	return ret
}

func decodeToneGray16(img *image.Gray16) *fp16.NRGBAh {
	bounds := img.Bounds()
	ret := fp16.NewNRGBAh(bounds)
	parallels.Parallel(bounds.Min.Y, bounds.Max.Y, func(y int) {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := img.Gray16At(x, y)
			fy := encodedToLinearTable16[c.Y]
			ret.SetNRGBAh(x, y, fp16color.NRGBAh{R: fy, G: fy, B: fy, A: one})
		}
	})
	return ret
}

func decodeToneAlpha(img *image.Alpha) *fp16.NRGBAh {
	bounds := img.Bounds()
	ret := fp16.NewNRGBAh(bounds)
	parallels.Parallel(bounds.Min.Y, bounds.Max.Y, func(y int) {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := img.AlphaAt(x, y)
			fa := float64(c.A) / 0xff
			ret.SetNRGBAh(x, y, fp16color.NRGBAh{R: one, G: one, B: one, A: float16.FromFloat64(fa)})
		}
	})
	return ret
}

func decodeToneAlpha16(img *image.Alpha16) *fp16.NRGBAh {
	bounds := img.Bounds()
	ret := fp16.NewNRGBAh(bounds)
	parallels.Parallel(bounds.Min.Y, bounds.Max.Y, func(y int) {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := img.Alpha16At(x, y)
			fa := float64(c.A) / 0xffff
			ret.SetNRGBAh(x, y, fp16color.NRGBAh{R: one, G: one, B: one, A: float16.FromFloat64(fa)})
		}
	})
	return ret
}

func decodeTonePaletted(img *image.Paletted) *fp16.NRGBAh {
	bounds := img.Bounds()
	ret := fp16.NewNRGBAh(bounds)
	palette := make([]fp16color.NRGBAh, len(img.Palette))
	for i, c := range img.Palette {
		r, g, b, a := c.RGBA()
		fr := float64(r) / 0xffff
		fg := float64(g) / 0xffff
		fb := float64(b) / 0xffff
		fa := float64(a) / 0xffff
		if a != 0 {
			fr /= fa
			fg /= fa
			fb /= fa
		}
		fr = encodedToLinear(fr)
		fg = encodedToLinear(fg)
		fb = encodedToLinear(fb)
		palette[i] = fp16color.NewNRGBAh(fr, fg, fb, fa)
	}
	parallels.Parallel(bounds.Min.Y, bounds.Max.Y, func(y int) {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := palette[img.ColorIndexAt(x, y)]
			ret.SetNRGBAh(x, y, c)
		}
	})
	return ret
}

func decodeTone(img image.Image) *fp16.NRGBAh {
	bounds := img.Bounds()
	ret := fp16.NewNRGBAh(bounds)
	parallels.Parallel(bounds.Min.Y, bounds.Max.Y, func(y int) {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := img.At(x, y)
			r, g, b, a := c.RGBA()
			if a == 0 {
				fr := encodedToLinearTable16[r]
				fg := encodedToLinearTable16[g]
				fb := encodedToLinearTable16[b]
				ret.SetNRGBAh(x, y, fp16color.NRGBAh{R: fr, G: fg, B: fb, A: 0})
			} else if a == 0xffff {
				fr := encodedToLinearTable16[r]
				fg := encodedToLinearTable16[g]
				fb := encodedToLinearTable16[b]
				ret.SetNRGBAh(x, y, fp16color.NRGBAh{R: fr, G: fg, B: fb, A: one})
			} else {
				fr := float64(r) / 0xffff
				fg := float64(g) / 0xffff
				fb := float64(b) / 0xffff
				fa := float64(a) / 0xffff
				fr /= fa
				fg /= fa
				fb /= fa
				fr = encodedToLinear(fr)
				fg = encodedToLinear(fg)
				fb = encodedToLinear(fb)
				ret.SetNRGBAh(x, y, fp16color.NewNRGBAh(fr, fg, fb, fa))
			}
		}
	})
	return ret
}

// EncodeTone encodes a linear color image to an 8bit depth sRGB color encoded image.
func EncodeTone(img *fp16.NRGBAh) *image.NRGBA {
	bounds := img.Bounds()
	ret := image.NewNRGBA(bounds)
	parallels.Parallel(bounds.Min.Y, bounds.Max.Y, func(y int) {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			rgba := img.NRGBAhAt(x, y)
			fa := max(0, min(1, rgba.A.Float64())) // clamp
			a := byte(math.RoundToEven(fa * 0xff))
			ret.SetNRGBA(x, y, color.NRGBA{
				R: linearToEncodedTable[rgba.R],
				G: linearToEncodedTable[rgba.G],
				B: linearToEncodedTable[rgba.B],
				A: a,
			})
		}
	})

	return ret
}

// EncodeTone64 encodes a linear color image to a 16bit depth sRGB color encoded image.
func EncodeTone64(img *fp16.NRGBAh) *image.NRGBA64 {
	bounds := img.Bounds()
	ret := image.NewNRGBA64(bounds)
	parallels.Parallel(bounds.Min.Y, bounds.Max.Y, func(y int) {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			rgba := img.NRGBAhAt(x, y)
			fa := max(0, min(1, rgba.A.Float64())) // clamp
			a := uint16(math.RoundToEven(fa * 0xffff))
			ret.SetNRGBA64(x, y, color.NRGBA64{
				R: linearToEncodedTable16[rgba.R],
				G: linearToEncodedTable16[rgba.G],
				B: linearToEncodedTable16[rgba.B],
				A: a,
			})
		}
	})

	return ret
}

func encodedToLinear(v float64) float64 {
	// https://en.wikipedia.org/wiki/SRGB#From_sRGB_to_CIE_XYZ
	if v <= 0.0031308*12.92 {
		return v / 12.92
	}
	return math.Pow((v+0.055)/1.055, 2.4)
}
