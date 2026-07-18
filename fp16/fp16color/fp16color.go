package fp16color

import (
	"image/color"

	"github.com/shogo82148/floats"
)

// RGBAh represents an alpha-premultiplied 64-bit color,
// having 16 bits float for each of red, green, blue and alpha.
type RGBAh struct {
	R, G, B, A floats.Float16
}

func NewRGBAh(r, g, b, a float64) RGBAh {
	fr := floats.NewFloat16(r)
	fg := floats.NewFloat16(g)
	fb := floats.NewFloat16(b)
	fa := floats.NewFloat16(a)
	return RGBAh{fr, fg, fb, fa}
}

func (c RGBAh) RGBA() (r, g, b, a uint32) {
	r = uint32(c.R.Float64().BuiltIn() * 0xffff)
	g = uint32(c.G.Float64().BuiltIn() * 0xffff)
	b = uint32(c.B.Float64().BuiltIn() * 0xffff)
	a = uint32(c.A.Float64().BuiltIn() * 0xffff)
	return
}

var RGBAhModel color.Model = color.ModelFunc(rgbahModel)

func rgbahModel(c color.Color) color.Color {
	if _, ok := c.(RGBAh); ok {
		return c
	}
	r, g, b, a := c.RGBA()
	return NRGBAh{
		R: floats.NewFloat16(float64(r) / 0xffff),
		G: floats.NewFloat16(float64(g) / 0xffff),
		B: floats.NewFloat16(float64(b) / 0xffff),
		A: floats.NewFloat16(float64(a) / 0xffff),
	}
}

var _ color.Color = NRGBAh{}

// NRGBAh represents a non-alpha-premultiplied 64-bit color,
// having 16 bits float for each of red, green, blue and alpha.
type NRGBAh struct {
	R, G, B, A floats.Float16
}

func NewNRGBAh(r, g, b, a float64) NRGBAh {
	fr := floats.NewFloat16(r)
	fg := floats.NewFloat16(g)
	fb := floats.NewFloat16(b)
	fa := floats.NewFloat16(a)
	return NRGBAh{fr, fg, fb, fa}
}

func (c NRGBAh) RGBA() (r, g, b, a uint32) {
	fa := c.A.Float64().BuiltIn()
	fr := c.R.Float64().BuiltIn() * fa
	fg := c.G.Float64().BuiltIn() * fa
	fb := c.B.Float64().BuiltIn() * fa
	r = uint32(fr * 0xffff)
	g = uint32(fg * 0xffff)
	b = uint32(fb * 0xffff)
	a = uint32(fa * 0xffff)
	return
}

var NRGBAhModel color.Model = color.ModelFunc(nrgbahModel)

func nrgbahModel(c color.Color) color.Color {
	if _, ok := c.(NRGBAh); ok {
		return c
	}
	r, g, b, a := c.RGBA()
	if a == 0 {
		return NRGBAh{0, 0, 0, 0}
	}
	// Since Color.RGBA returns an alpha-premultiplied color, we should have r <= a && g <= a && b <= a.
	fa := float64(a) / 0xffff
	factor := 1 / float64(a)
	fr := float64(r) * factor
	fg := float64(g) * factor
	fb := float64(b) * factor
	return NRGBAh{
		R: floats.NewFloat16(fr),
		G: floats.NewFloat16(fg),
		B: floats.NewFloat16(fb),
		A: floats.NewFloat16(fa),
	}
}
