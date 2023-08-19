package pixmap

import (
	"image"
	"image/color"
	"image/draw"
)

var _ color.Color = Color{}

// Color represents a gray color in any depth.
type Color struct {
	R   uint16
	G   uint16
	B   uint16
	Max Model
}

func (c Color) RGBA() (r, g, b, a uint32) {
	r = (uint32(c.R) * 0xffff) / uint32(c.Max)
	g = (uint32(c.G) * 0xffff) / uint32(c.Max)
	b = (uint32(c.B) * 0xffff) / uint32(c.Max)
	a = 0xffff
	return
}

// Model represents a gray color model.
type Model uint16

var _ color.Model = Model(0)

func (m Model) Convert(c color.Color) color.Color {
	var r, g, b uint32
	if c, ok := c.(Color); ok {
		r = (uint32(c.R) * uint32(m)) / uint32(c.Max)
		g = (uint32(c.G) * uint32(m)) / uint32(c.Max)
		b = (uint32(c.B) * uint32(m)) / uint32(c.Max)
	} else {
		r, g, b, _ = c.RGBA()
		r = (r * uint32(m)) / 0xffff
		g = (g * uint32(m)) / 0xffff
		b = (b * uint32(m)) / 0xffff
	}
	return Color{R: uint16(r), G: uint16(g), B: uint16(b), Max: m}
}

// Image represents a Pix Map image.
type Image struct {
	Pix    []uint8
	Stride int
	Rect   image.Rectangle
	Max    Model
}

var _ image.Image = (*Image)(nil)
var _ draw.Image = (*Image)(nil)

// New returns a new Image with the given bounds and maximum value.
func New(r image.Rectangle, m Model) *Image {
	stride := r.Dx() * 3
	if m >= 256 {
		stride *= 2 // 16bit
	}
	return &Image{
		Pix:    make([]uint8, r.Dy()*stride),
		Stride: stride,
		Rect:   r,
		Max:    m,
	}
}

// ColorModel implements [image.Image].
func (img *Image) ColorModel() color.Model {
	return img.Max
}

// Bounds implements [image.Image].
func (img *Image) Bounds() image.Rectangle {
	return img.Rect
}

func (img *Image) At(x, y int) color.Color {
	return img.PixAt(x, y)
}

func (img *Image) PixAt(x, y int) Color {
	if !(image.Point{x, y}.In(img.Rect)) {
		return Color{
			R:   uint16(img.Max),
			G:   uint16(img.Max),
			B:   uint16(img.Max),
			Max: img.Max,
		}
	}
	if img.Max >= 256 {
		offset := (y-img.Rect.Min.Y)*img.Stride + (x-img.Rect.Min.X)*2*3
		return Color{
			R:   uint16(img.Pix[offset+0])<<8 | uint16(img.Pix[offset+1]),
			G:   uint16(img.Pix[offset+2])<<8 | uint16(img.Pix[offset+3]),
			B:   uint16(img.Pix[offset+4])<<8 | uint16(img.Pix[offset+5]),
			Max: img.Max,
		}
	} else {
		offset := (y-img.Rect.Min.Y)*img.Stride + (x-img.Rect.Min.X)*3
		return Color{
			R:   uint16(img.Pix[offset+0]),
			G:   uint16(img.Pix[offset+1]),
			B:   uint16(img.Pix[offset+2]),
			Max: img.Max,
		}
	}
}

func (img *Image) Set(x, y int, c color.Color) {
	img.SetPix(x, y, img.Max.Convert(c).(Color))
}

func (img *Image) SetPix(x, y int, c Color) {
	if !(image.Point{x, y}.In(img.Rect)) {
		return
	}

	if img.Max >= 256 {
		offset := (y-img.Rect.Min.Y)*img.Stride + (x-img.Rect.Min.X)*2*3
		img.Pix[offset+0] = uint8(c.R >> 8)
		img.Pix[offset+1] = uint8(c.R)
		img.Pix[offset+2] = uint8(c.G >> 8)
		img.Pix[offset+3] = uint8(c.G)
		img.Pix[offset+4] = uint8(c.B >> 8)
		img.Pix[offset+5] = uint8(c.B)
	} else {
		offset := (y-img.Rect.Min.Y)*img.Stride + (x-img.Rect.Min.X)*3
		img.Pix[offset+0] = uint8(c.R)
		img.Pix[offset+1] = uint8(c.G)
		img.Pix[offset+2] = uint8(c.B)
	}
}
