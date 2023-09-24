package png

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"hash/crc32"
	"image"
	"image/color"
	"io"
	"math"
	"strconv"

	"github.com/shogo82148/go-imaging/icc"
)

// ImageWithMeta is a PNG image with metadata.
type ImageWithMeta struct {
	image.Image

	// Gamma is the gamma value of the image.
	// If Gamma is 0, the image has no gamma information.
	Gamma float64

	// SRGB is the sRGB information of the image.
	// If SRGB is nil, the image has no sRGB information.
	SRGB *SRGB

	// ICCProfileName is the name of the ICC profile of the image.
	ICCProfileName string

	// ICCProfile is the ICC profile of the image.
	// If ICCProfile is nil, the image has no ICC profile.
	ICCProfile *icc.Profile
}

type SRGB struct {
	RenderingIntent RenderingIntent
}

type RenderingIntent int

const (
	// RenderingIntentPerceptual is for images preferring good adaptation
	// to the output device gamut at the expense of colorimetric accuracy, such as photographs.
	RenderingIntentPerceptual RenderingIntent = 0

	// RenderingIntentRelative is for images requiring color
	// appearance matching (relative to the output device white point), such as logos.
	RenderingIntentRelative RenderingIntent = 1

	// RenderingIntentSaturation is for images preferring preservation
	// of saturation at the expense of hue and lightness, such as charts and graphs.
	RenderingIntentSaturation RenderingIntent = 2

	// RenderingIntentAbsolute is for images requiring preservation of absolute colorimetry,
	// such as previews of images destined for a different output device (proofs).
	RenderingIntentAbsolute RenderingIntent = 3
)

func (ri RenderingIntent) String() string {
	switch ri {
	case RenderingIntentPerceptual:
		return "Perceptual"
	case RenderingIntentRelative:
		return "Relative"
	case RenderingIntentSaturation:
		return "Saturation"
	case RenderingIntentAbsolute:
		return "Absolute"
	default:
		return "Unknown RenderingIntent: " + strconv.Itoa(int(ri))
	}
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
	img.SRGB = d.srgb
	if d.icc != nil {
		img.ICCProfileName = d.profileName
		img.ICCProfile = d.icc
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
	if m.SRGB != nil {
		e.writeSRGB(m.SRGB)
	}
	if m.ICCProfile != nil {
		e.writeICCP(m.ICCProfileName, m.ICCProfile)
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

func (e *encoder) writeSRGB(srgb *SRGB) {
	e.tmp[0] = byte(srgb.RenderingIntent)
	e.writeChunk(e.tmp[:1], "sRGB")
}

// iccProfileLatin1ToUTF8 converts a string from Latin-1 to UTF-8.
// Leading, trailing, and consecutive spaces are not permitted.
func iccProfileLatin1ToUTF8(s string) (string, bool) {
	if len(s) == 0 {
		return "", true
	}
	if s[0] == ' ' || s[len(s)-1] == ' ' {
		return "", false
	}

	runes := make([]rune, 0, len(s))
	sp := false
	for _, ch := range []byte(s) {
		if (ch < 32 || ch > 126) && (ch < 161 || ch > 255) {
			return "", false
		}
		if sp && ch == ' ' {
			return "", false
		}
		sp = ch == ' '
		runes = append(runes, rune(ch))
	}
	return string(runes), true
}

func iccProfileUTF8ToLatin1(s string) string {
	buf := make([]byte, 0, len(s))
	var sp bool
	for _, ch := range s {
		if (ch < 32 || ch > 126) && (ch < 161 || ch > 255) {
			ch = '.'
		}
		if sp && ch == ' ' {
			continue
		}
		sp = ch == ' '
		buf = append(buf, byte(ch))
	}
	return string(bytes.TrimSpace(buf))
}

func (e *encoder) writeICCP(name string, profile *icc.Profile) {
	if e.err != nil {
		return
	}

	buf := bytes.NewBuffer(e.tmp[:0])
	buf.WriteString(iccProfileUTF8ToLatin1(name))
	buf.WriteByte(0x00) // null terminator
	buf.WriteByte(0x00) // compression method: zlib

	// Compress the ICC profile
	w, err := zlib.NewWriterLevel(buf, zlib.BestCompression)
	if err != nil {
		e.err = err
		return
	}
	e.err = profile.Encode(w)
	if e.err != nil {
		return
	}
	e.err = w.Close()
	if e.err != nil {
		return
	}
	e.writeChunk(buf.Bytes(), "iCCP")
}
