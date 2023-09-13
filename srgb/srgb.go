//go:generate go run internal/gen/main.go

// Package srgb handles [sRGB] colors.
//
// [sRGB]: https://en.wikipedia.org/wiki/SRGB
package srgb

import (
	"image"
	"image/color"
	"math"

	"github.com/shogo82148/go-imaging/fp16"
	"github.com/shogo82148/go-imaging/fp16/fp16color"
	"github.com/shogo82148/go-imaging/internal/parallels"
)

// Decode decodes an sRGB color encoded image to a linear color image.
func Decode(img image.Image) *fp16.NRGBAh {
	switch img := img.(type) {
	case *image.RGBA:
		return decodeRGBA(img)
	}
	return decode(img)
}

func decodeRGBA(img *image.RGBA) *fp16.NRGBAh {
	bounds := img.Bounds()
	ret := fp16.NewNRGBAh(bounds)
	parallels.Parallel(bounds.Min.Y, bounds.Max.Y, func(y int) {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := img.RGBAAt(x, y)
			fr := float64(c.R) / 0xff
			fg := float64(c.G) / 0xff
			fb := float64(c.B) / 0xff
			fa := float64(c.A) / 0xff
			if fa != 0 {
				fr /= fa
				fg /= fa
				fb /= fa
			}
			fr = encodedToLinear(fr)
			fg = encodedToLinear(fg)
			fb = encodedToLinear(fb)
			ret.SetNRGBAhAt(x, y, fp16color.NewNRGBAh(fr, fg, fb, fa))
		}
	})
	return ret
}

func decode(img image.Image) *fp16.NRGBAh {
	bounds := img.Bounds()
	ret := fp16.NewNRGBAh(bounds)
	parallels.Parallel(bounds.Min.Y, bounds.Max.Y, func(y int) {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := img.At(x, y)
			r, g, b, a := c.RGBA()
			fr := float64(r) / 0xffff
			fg := float64(g) / 0xffff
			fb := float64(b) / 0xffff
			fa := float64(a) / 0xffff
			if fa != 0 {
				fr /= fa
				fg /= fa
				fb /= fa
			}
			fr = encodedToLinear(fr)
			fg = encodedToLinear(fg)
			fb = encodedToLinear(fb)
			ret.SetNRGBAhAt(x, y, fp16color.NewNRGBAh(fr, fg, fb, fa))
		}
	})
	return ret
}

// Encode encodes a linear color image to an sRGB color encoded image.
func Encode(img *fp16.NRGBAh) *image.NRGBA64 {
	bounds := img.Bounds()
	ret := image.NewNRGBA64(bounds)
	parallels.Parallel(bounds.Min.Y, bounds.Max.Y, func(y int) {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			rgba := img.NRGBAhAt(x, y)
			a := uint16(math.RoundToEven(rgba.A.Float64() * 0xffff))
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
