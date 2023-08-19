package pnm

import (
	"fmt"
	"image"
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
			typ = TypePBM
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

func (enc *Encoder) encodeP2(w io.Writer, m image.Image) error {
	return nil
}

func (enc *Encoder) encodeP3(w io.Writer, m image.Image) error {
	return nil
}

func (enc *Encoder) encodeP4(w io.Writer, m image.Image) error {
	return nil
}

func (enc *Encoder) encodeP5(w io.Writer, m image.Image) error {
	return nil
}

func (enc *Encoder) encodeP6(w io.Writer, m image.Image) error {
	return nil
}
