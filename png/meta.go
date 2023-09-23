package png

import (
	"encoding/binary"
	"hash/crc32"
	"image"
	"image/color"
	"io"
	"math"
	"strconv"
)

// ImageWithMeta is a PNG image with metadata.
type ImageWithMeta struct {
	image.Image

	// Gamma is the gamma value of the image.
	// If Gamma is 0, the image has no gamma information.
	Gamma float64
}

// DecodeWithMeta reads a PNG image from r and returns it as an image.Image.
// The type of Image returned depends on the PNG contents.
func DecodeWithMeta(r io.Reader) (*ImageWithMeta, error) {
	d := &decoder{
		r:   r,
		crc: crc32.NewIEEE(),
	}
	if err := d.checkHeader(); err != nil {
		if err == io.EOF {
			err = io.ErrUnexpectedEOF
		}
		return nil, err
	}
	for d.stage != dsSeenIEND {
		if err := d.parseChunk(false); err != nil {
			if err == io.EOF {
				err = io.ErrUnexpectedEOF
			}
			return nil, err
		}
	}

	img := &ImageWithMeta{
		Image: d.img,
	}
	if d.gamma != 0 {
		img.Gamma = float64(d.gamma) / 100000
	}
	return img, nil
}

// EncodeWithMeta writes the Image m to w in PNG format. Any Image may be
// encoded, but images that are not image.NRGBA might be encoded lossily.
func EncodeWithMeta(w io.Writer, m *ImageWithMeta) error {
	var e Encoder
	return e.EncodeWithMeta(w, m)
}

// Encode writes the Image m to w in PNG format.
func (enc *Encoder) EncodeWithMeta(w io.Writer, m *ImageWithMeta) error {
	// Obviously, negative widths and heights are invalid. Furthermore, the PNG
	// spec section 11.2.2 says that zero is invalid. Excessively large images are
	// also rejected.
	mw, mh := int64(m.Bounds().Dx()), int64(m.Bounds().Dy())
	if mw <= 0 || mh <= 0 || mw >= 1<<32 || mh >= 1<<32 {
		return FormatError("invalid image size: " + strconv.FormatInt(mw, 10) + "x" + strconv.FormatInt(mh, 10))
	}

	var e *encoder
	if enc.BufferPool != nil {
		buffer := enc.BufferPool.Get()
		e = (*encoder)(buffer)

	}
	if e == nil {
		e = &encoder{}
	}
	if enc.BufferPool != nil {
		defer enc.BufferPool.Put((*EncoderBuffer)(e))
	}

	e.enc = enc
	e.w = w
	e.m = m

	var pal color.Palette
	// cbP8 encoding needs PalettedImage's ColorIndexAt method.
	if _, ok := m.Image.(image.PalettedImage); ok {
		pal, _ = m.ColorModel().(color.Palette)
	}
	if pal != nil {
		if len(pal) <= 2 {
			e.cb = cbP1
		} else if len(pal) <= 4 {
			e.cb = cbP2
		} else if len(pal) <= 16 {
			e.cb = cbP4
		} else {
			e.cb = cbP8
		}
	} else {
		switch m.ColorModel() {
		case color.GrayModel:
			e.cb = cbG8
		case color.Gray16Model:
			e.cb = cbG16
		case color.RGBAModel, color.NRGBAModel, color.AlphaModel:
			if opaque(m) {
				e.cb = cbTC8
			} else {
				e.cb = cbTCA8
			}
		default:
			if opaque(m) {
				e.cb = cbTC16
			} else {
				e.cb = cbTCA16
			}
		}
	}

	_, e.err = io.WriteString(w, pngHeader)
	e.writeIHDR()
	if m.Gamma != 0 {
		e.writeGAMA(m.Gamma)
	}
	if pal != nil {
		e.writePLTEAndTRNS(pal)
	}
	e.writeIDATs()
	e.writeIEND()
	return e.err
}

func (e *encoder) writeGAMA(gamma float64) {
	g := uint32(math.RoundToEven(gamma * 100000))
	binary.BigEndian.PutUint32(e.tmp[:4], g)
	e.writeChunk(e.tmp[:4], "gAMA")
}
