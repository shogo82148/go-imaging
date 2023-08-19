package pnm

import (
	"fmt"
	"image"
	"image/draw"
	"io"

	"github.com/shogo82148/go-imaging/bitmap"
	"github.com/shogo82148/go-imaging/graymap"
	"github.com/shogo82148/go-imaging/pixmap"
)

type Type int

const (
	TypeAuto Type = 0
	TypePBM  Type = 1
	TypePGM  Type = 2
	TypePPM  Type = 3
)

var defaultEncoder = &Encoder{}

// Encode writes the image m to w in PNM format.
func Encode(w io.Writer, m image.Image) error {
	return defaultEncoder.Encode(w, m)
}

// Encoder is a Portable Any Map image encoder.
type Encoder struct {
	// Plain specifies whether to encode in plain(ASCII) format.
	Plain bool

	// Type specifies the image type.
	Type Type

	// Max specifies the maximum value of the image color.
	// If Max is 0, the maximum value is determined by the image type.
	// If Type is TypePBM, Max is ignored.
	Max uint16
}

func (enc *Encoder) Encode(w io.Writer, m image.Image) error {
	typ := enc.Type
	if typ == TypeAuto {
		// auto detect type
		switch m.(type) {
		case *bitmap.Image:
			typ = TypePBM
		case *graymap.Image, *image.Gray, *image.Gray16:
			typ = TypePGM
		case *pixmap.Image, *image.RGBA64, *image.NRGBA64, *image.RGBA, *image.NRGBA:
			typ = TypePPM
		default:
			typ = TypePPM
		}
	}
	switch typ {
	case TypePBM:
		if enc.Plain {
			return enc.encodeP1(w, m)
		} else {
			return enc.encodeP4(w, m)
		}
	case TypePGM:
		if enc.Plain {
			return enc.encodeP2(w, m)
		} else {
			return enc.encodeP5(w, m)
		}
	case TypePPM:
		if enc.Plain {
			return enc.encodeP3(w, m)
		} else {
			return enc.encodeP6(w, m)
		}
	default:
		return fmt.Errorf("pnm: unknown encoding type: %d", typ)
	}
}

// encodeP1 encodes a plain Portable Bit Map image.
// See https://netpbm.sourceforge.net/doc/pbm.html
func (enc *Encoder) encodeP1(w io.Writer, m image.Image) error {
	bounds := m.Bounds()
	if _, err := fmt.Fprintf(w, "P1\n%d %d\n", bounds.Dx(), bounds.Dy()); err != nil {
		return err
	}
	for y := bounds.Min.Y; y < m.Bounds().Max.Y; y++ {
		for x := bounds.Min.X; x < m.Bounds().Max.X; x++ {
			var d int
			c := m.At(x, y)
			if bitmap.ColorModel.Convert(c).(bitmap.Color) {
				d = 1
			} else {
				d = 0
			}
			if _, err := fmt.Fprintf(w, "%d", d); err != nil {
				return err
			}
		}
		if _, err := fmt.Fprintln(w); err != nil {
			return err
		}
	}
	return nil
}

// encodeP2 encodes a plain Portable Gray Map image.
// See https://netpbm.sourceforge.net/doc/pgm.html
func (enc *Encoder) encodeP2(w io.Writer, m image.Image) error {
	maxValue := graymap.Model(enc.Max)
	if maxValue == 0 {
		switch m := m.(type) {
		case *graymap.Image:
			maxValue = m.Max
		case *image.Gray, *image.RGBA, *image.NRGBA:
			maxValue = 0xff
		case *image.Gray16, *image.RGBA64, *image.NRGBA64:
			maxValue = 0xffff
		default:
			maxValue = 0xff
		}
	}

	bounds := m.Bounds()
	if _, err := fmt.Fprintf(w, "P2\n%d %d\n%d\n", bounds.Dx(), bounds.Dy(), maxValue); err != nil {
		return err
	}
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		c := m.At(bounds.Min.X, y)
		gc := maxValue.Convert(c).(graymap.Color)
		if _, err := fmt.Fprintf(w, "%d", gc.Y); err != nil {
			return err
		}
		for x := bounds.Min.X + 1; x < bounds.Max.X; x++ {
			c := m.At(x, y)
			gc := maxValue.Convert(c).(graymap.Color)
			if _, err := fmt.Fprintf(w, " %d", gc.Y); err != nil {
				return err
			}
		}
		if _, err := fmt.Fprintln(w); err != nil {
			return err
		}
	}
	return nil
}

// encodeP3 encodes a plain Portable Pix Map image.
// See https://netpbm.sourceforge.net/doc/ppm.html
func (enc *Encoder) encodeP3(w io.Writer, m image.Image) error {
	maxValue := pixmap.Model(enc.Max)
	if maxValue == 0 {
		switch m := m.(type) {
		case *pixmap.Image:
			maxValue = m.Max
		case *image.Gray, *image.RGBA, *image.NRGBA:
			maxValue = 0xff
		case *image.Gray16, *image.RGBA64, *image.NRGBA64:
			maxValue = 0xffff
		default:
			maxValue = 0xff
		}
	}

	bounds := m.Bounds()
	if _, err := fmt.Fprintf(w, "P3\n%d %d\n%d\n", bounds.Dx(), bounds.Dy(), maxValue); err != nil {
		return err
	}
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		c := m.At(bounds.Min.X, y)
		pc := maxValue.Convert(c).(pixmap.Color)
		if _, err := fmt.Fprintf(w, "%d %d %d", pc.R, pc.G, pc.B); err != nil {
			return err
		}
		for x := bounds.Min.X + 1; x < bounds.Max.X; x++ {
			c := m.At(x, y)
			pc := maxValue.Convert(c).(pixmap.Color)
			if _, err := fmt.Fprintf(w, " %d %d %d", pc.R, pc.G, pc.B); err != nil {
				return err
			}
		}
		if _, err := fmt.Fprintln(w); err != nil {
			return err
		}
	}

	return nil
}

// encodeP4 encodes a raw Portable Bit Map image.
// See https://netpbm.sourceforge.net/doc/pbm.html
func (enc *Encoder) encodeP4(w io.Writer, m image.Image) error {
	img := bitmap.New(m.Bounds())
	clone(img, m)

	bounds := m.Bounds()
	if _, err := fmt.Fprintf(w, "P4\n%d %d\n", bounds.Dx(), bounds.Dy()); err != nil {
		return err
	}
	if _, err := w.Write(img.Pix); err != nil {
		return err
	}
	return nil
}

// encodePt encodes a raw Portable Gray Map image.
// See https://netpbm.sourceforge.net/doc/pgm.html
func (enc *Encoder) encodeP5(w io.Writer, m image.Image) error {
	maxValue := graymap.Model(enc.Max)
	if maxValue == 0 {
		switch m := m.(type) {
		case *graymap.Image:
			maxValue = m.Max
		case *image.Gray, *image.RGBA, *image.NRGBA:
			maxValue = 0xff
		case *image.Gray16, *image.RGBA64, *image.NRGBA64:
			maxValue = 0xffff
		default:
			maxValue = 0xff
		}
	}

	bounds := m.Bounds()
	img := graymap.New(bounds, maxValue)
	clone(img, m)
	if _, err := fmt.Fprintf(w, "P5\n%d %d\n%d\n", bounds.Dx(), bounds.Dy(), maxValue); err != nil {
		return err
	}
	if _, err := w.Write(img.Pix); err != nil {
		return err
	}
	return nil
}

func (enc *Encoder) encodeP6(w io.Writer, m image.Image) error {
	maxValue := pixmap.Model(enc.Max)
	if maxValue == 0 {
		switch m := m.(type) {
		case *pixmap.Image:
			maxValue = m.Max
		case *image.Gray, *image.RGBA, *image.NRGBA:
			maxValue = 0xff
		case *image.Gray16, *image.RGBA64, *image.NRGBA64:
			maxValue = 0xffff
		default:
			maxValue = 0xff
		}
	}

	bounds := m.Bounds()
	img := pixmap.New(bounds, maxValue)
	clone(img, m)
	if _, err := fmt.Fprintf(w, "P6\n%d %d\n%d\n", bounds.Dx(), bounds.Dy(), maxValue); err != nil {
		return err
	}
	if _, err := w.Write(img.Pix); err != nil {
		return err
	}
	return nil
}

func clone(dst draw.Image, src image.Image) {
	bounds := dst.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			dst.Set(x, y, src.At(x, y))
		}
	}
}
