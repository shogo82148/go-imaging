//go:generate go run internal/gen/main.go

// Package srgb handles [sRGB] colors.
//
// [sRGB]: https://en.wikipedia.org/wiki/SRGB
package srgb

import (
	"image"
	"image/color"

	"github.com/shogo82148/go-imaging/internal/parallels"
)

type nrgba64At interface {
	image.Image
	NRGBA64At(x, y int) color.NRGBA64
}

type nrgbaAt interface {
	image.Image
	NRGBAAt(x, y int) color.NRGBA
}

type rgba64At interface {
	image.Image
	RGBA64At(x, y int) color.RGBA64
}

type rgbaAt interface {
	image.Image
	RGBAAt(x, y int) color.RGBA
}

// Decode decodes an sRGB color encoded image to a linear color image.
func Decode(img image.Image) *image.NRGBA64 {
	switch img := img.(type) {
	case nrgba64At:
		return decodeNRGBA64(img)
	case nrgbaAt:
		return decodeNRGBA(img)
	case rgba64At:
		return decodeRGBA64(img)
	case rgbaAt:
		return decodeRGBA(img)
	}
	return decode(img)
}

func decodeNRGBA64(img nrgba64At) *image.NRGBA64 {
	bounds := img.Bounds()
	ret := image.NewNRGBA64(bounds)
	parallels.Parallel(bounds.Min.Y, bounds.Max.Y, func(y int) {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			rgba := img.NRGBA64At(x, y)
			ret.SetNRGBA64(x, y, color.NRGBA64{
				R: encodedToLinearTable16[rgba.R],
				G: encodedToLinearTable16[rgba.G],
				B: encodedToLinearTable16[rgba.B],
				A: rgba.A,
			})
		}
	})

	return ret
}

func decodeNRGBA(img nrgbaAt) *image.NRGBA64 {
	bounds := img.Bounds()
	ret := image.NewNRGBA64(bounds)
	parallels.Parallel(bounds.Min.Y, bounds.Max.Y, func(y int) {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			rgba := img.NRGBAAt(x, y)
			a := uint16(rgba.A)
			a |= a << 8
			ret.SetNRGBA64(x, y, color.NRGBA64{
				R: encodedToLinearTable8[rgba.R],
				G: encodedToLinearTable8[rgba.G],
				B: encodedToLinearTable8[rgba.B],
				A: a,
			})
		}
	})

	return ret
}

func decodeRGBA64(img rgba64At) *image.NRGBA64 {
	bounds := img.Bounds()
	ret := image.NewNRGBA64(bounds)
	parallels.Parallel(bounds.Min.Y, bounds.Max.Y, func(y int) {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			rgba := img.RGBA64At(x, y)
			ret.SetNRGBA64(x, y, convertRGBA64ToNRGBA64(rgba))
		}
	})

	return ret
}

func decodeRGBA(img rgbaAt) *image.NRGBA64 {
	bounds := img.Bounds()
	ret := image.NewNRGBA64(bounds)
	parallels.Parallel(bounds.Min.Y, bounds.Max.Y, func(y int) {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			rgba := img.RGBAAt(x, y)
			ret.SetNRGBA64(x, y, convertRGBAToNRGBA64(rgba))
		}
	})

	return ret
}

func decode(img image.Image) *image.NRGBA64 {
	bounds := img.Bounds()
	ret := image.NewNRGBA64(bounds)
	parallels.Parallel(bounds.Min.Y, bounds.Max.Y, func(y int) {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := img.At(x, y)
			ret.SetNRGBA64(x, y, convertToNRGBA64(c))
		}
	})

	return ret
}

func convertRGBA64ToNRGBA64(c color.RGBA64) color.NRGBA64 {
	if c.A == 0xffff {
		return color.NRGBA64{encodedToLinearTable16[c.R], encodedToLinearTable16[c.G], encodedToLinearTable16[c.B], 0xffff}
	}
	if c.A == 0 {
		return color.NRGBA64{0, 0, 0, 0}
	}
	// Since Color.RGBA returns an alpha-premultiplied color, we should have r <= a && g <= a && b <= a.
	r := (c.R * 0xffff) / c.A
	g := (c.G * 0xffff) / c.A
	b := (c.B * 0xffff) / c.A
	return color.NRGBA64{encodedToLinearTable16[r], encodedToLinearTable16[g], encodedToLinearTable16[b], c.A}
}

func convertRGBAToNRGBA64(c color.RGBA) color.NRGBA64 {
	if c.A == 0xff {
		return color.NRGBA64{encodedToLinearTable8[c.R], encodedToLinearTable8[c.G], encodedToLinearTable8[c.B], 0xffff}
	}
	if c.A == 0 {
		return color.NRGBA64{0, 0, 0, 0}
	}

	r := uint32(c.R)
	g := uint32(c.G)
	b := uint32(c.B)
	a := uint32(c.A)

	// Since Color.RGBA returns an alpha-premultiplied color, we should have r <= a && g <= a && b <= a.
	r = (r * 0xffff) / a
	g = (g * 0xffff) / a
	b = (g * 0xffff) / a
	a = a << 8
	return color.NRGBA64{encodedToLinearTable16[r], encodedToLinearTable16[g], encodedToLinearTable16[b], uint16(a)}
}

func convertToNRGBA64(c color.Color) color.NRGBA64 {
	if c, ok := c.(color.NRGBA64); ok {
		return color.NRGBA64{encodedToLinearTable16[c.R], encodedToLinearTable16[c.G], encodedToLinearTable16[c.B], c.A}
	}
	r, g, b, a := c.RGBA()
	if a == 0xffff {
		return color.NRGBA64{encodedToLinearTable16[r], encodedToLinearTable16[g], encodedToLinearTable16[b], 0xffff}
	}
	if a == 0 {
		return color.NRGBA64{0, 0, 0, 0}
	}
	// Since Color.RGBA returns an alpha-premultiplied color, we should have r <= a && g <= a && b <= a.
	r = (r * 0xffff) / a
	g = (g * 0xffff) / a
	b = (b * 0xffff) / a
	return color.NRGBA64{encodedToLinearTable16[r], encodedToLinearTable16[g], encodedToLinearTable16[b], uint16(a)}

}

// Encode encodes a linear color image to an sRGB color encoded image.
func Encode(img *image.NRGBA64) *image.NRGBA64 {
	bounds := img.Bounds()
	ret := image.NewNRGBA64(bounds)
	parallels.Parallel(bounds.Min.Y, bounds.Max.Y, func(y int) {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			rgba := img.NRGBA64At(x, y)
			ret.SetRGBA64(x, y, color.RGBA64{
				R: linearToEncodedTable16[rgba.R],
				G: linearToEncodedTable16[rgba.G],
				B: linearToEncodedTable16[rgba.B],
				A: rgba.A,
			})
		}
	})

	return ret
}
