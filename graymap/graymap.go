package graymap

import (
	"image"
	"image/color"
	"image/draw"
)

var _ color.Color = Color{}

// Color represents a gray color in any depth.
type Color struct {
	Y   uint16
	Max Model
}

func (c Color) RGBA() (r, g, b, a uint32) {
	y := uint32(c.Y)
	y = (y * 0xffff) / uint32(c.Max)
	return y, y, y, 0xffff
}

// Model represents a gray color model.
type Model uint16

var _ color.Model = Model(0)

func (m Model) Convert(c color.Color) color.Color {
	if c, ok := c.(Color); ok {
		y := uint32(c.Y)
		y = (y * uint32(m)) / uint32(c.Max)
		return Color{Y: uint16(y), Max: m}
	}

	r, g, b, _ := c.RGBA()

	// These coefficients (the fractions 0.299, 0.587 and 0.114) are the same
	// as those given by the JFIF specification and used by func RGBToYCbCr in
	// ycbcr.go.
	//
	// Note that 19595 + 38470 + 7471 equals 65536.
	y := (19595*r + 38470*g + 7471*b + 1<<15) >> 16

	y = (y * uint32(m)) / 0xffff
	return Color{Y: uint16(y), Max: m}
}

type Image struct {
	Pix    []uint8
	Stride int
	Rect   image.Rectangle
	Max    Model
}

var _ image.Image = (*Image)(nil)
var _ draw.Image = (*Image)(nil)

func New(r image.Rectangle, m Model) *Image {
	stride := r.Dx()
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
	return img.GrayAt(x, y)
}

func (img *Image) GrayAt(x, y int) Color {
	if !(image.Point{x, y}.In(img.Rect)) {
		return Color{Y: uint16(img.Max), Max: img.Max}
	}
	if img.Max >= 256 {
		offset := (y-img.Rect.Min.Y)*img.Stride + (x-img.Rect.Min.X)*2
		return Color{
			Y:   uint16(img.Pix[offset+0])<<8 | uint16(img.Pix[offset+1]),
			Max: img.Max,
		}
	} else {
		offset := (y-img.Rect.Min.Y)*img.Stride + (x - img.Rect.Min.X)
		return Color{
			Y:   uint16(img.Pix[offset]),
			Max: img.Max,
		}
	}
}

func (img *Image) Set(x, y int, c color.Color) {
	img.SetGray(x, y, img.Max.Convert(c).(Color))
}

func (img *Image) SetGray(x, y int, c Color) {
	if !(image.Point{x, y}.In(img.Rect)) {
		return
	}
	cy := uint32(c.Y)
	if c.Max != img.Max {
		cy = (cy * uint32(img.Max)) / uint32(c.Max)
	}
	if img.Max >= 256 {
		offset := (y-img.Rect.Min.Y)*img.Stride + (x-img.Rect.Min.X)*2
		img.Pix[offset+0] = uint8(cy >> 8)
		img.Pix[offset+1] = uint8(cy)
	} else {
		offset := (y-img.Rect.Min.Y)*img.Stride + (x - img.Rect.Min.X)
		img.Pix[offset] = uint8(cy)
	}
}
