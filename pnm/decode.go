// Package pnm support Portable Any Map image format.
//
// [netpbm]: https://en.wikipedia.org/wiki/Netpbm
package pnm

import (
	"bufio"
	"errors"
	"fmt"
	"image"
	"image/color"
	"io"
	"math"
	"math/bits"

	"github.com/shogo82148/go-imaging/bitmap"
	"github.com/shogo82148/go-imaging/graymap"
	"github.com/shogo82148/go-imaging/pixmap"
)

func init() {
	image.RegisterFormat("pbm ascii", "P1", Decode, DecodeConfig)
	image.RegisterFormat("pgm ascii", "P2", Decode, DecodeConfig)
	image.RegisterFormat("ppm ascii", "P3", Decode, DecodeConfig)
	image.RegisterFormat("pbm binary", "P4", Decode, DecodeConfig)
	image.RegisterFormat("pgm binary", "P5", Decode, DecodeConfig)
	image.RegisterFormat("ppm binary", "P6", Decode, DecodeConfig)
}

type config struct {
	MagicNumber uint16
	Width       int
	Height      int
}

func Decode(r io.Reader) (image.Image, error) {
	br := bufio.NewReader(r)
	c, err := decodeConfig(br)
	if err != nil {
		return nil, err
	}
	switch c.MagicNumber {
	case 0x5031: // P1
		return decodeP1(br, c)
	case 0x5032: // P2
		return decodeP2(br, c)
	case 0x5033: // P3
		return decodeP3(br, c)
	case 0x5034: // P4
		return decodeP4(br, c)
	case 0x5035: // P5
		return decodeP5(br, c)
	case 0x5036: // P6
		return decodeP6(br, c)
	}
	return nil, errors.New("pnm: unsupported format")
}

// decodeP1 decodes a plain Portable Bit Map image.
// See https://netpbm.sourceforge.net/doc/pbm.html
func decodeP1(br *bufio.Reader, c config) (image.Image, error) {
	img := bitmap.New(image.Rect(0, 0, c.Width, c.Height))
	for y := 0; y < c.Height; y++ {
		for x := 0; x < c.Width; x++ {
			if err := skipWhitespace(br); err != nil {
				return nil, err
			}
			b, err := br.ReadByte()
			if err != nil {
				return nil, err
			}
			switch b {
			case '0':
				img.Set(x, y, bitmap.White)
			case '1':
				img.Set(x, y, bitmap.Black)
			default:
				return nil, fmt.Errorf("pnm: unexpected char: %c", b)
			}
		}
	}
	return img, nil
}

// decodeP2 decodes a plain Portable Gray Map image.
// See https://netpbm.sourceforge.net/doc/pgm.html
func decodeP2(br *bufio.Reader, c config) (image.Image, error) {
	if err := skipWhitespace(br); err != nil {
		return nil, err
	}
	maxVal, err := readInt(br)
	if err != nil {
		return nil, err
	}
	if maxVal == 0 || maxVal > 0xffff {
		return nil, fmt.Errorf("pnm: unsupported max value: %d", maxVal)
	}

	img := graymap.New(image.Rect(0, 0, c.Width, c.Height), graymap.Model(maxVal))
	for y := 0; y < c.Height; y++ {
		for x := 0; x < c.Width; x++ {
			if err := skipWhitespace(br); err != nil {
				return nil, err
			}
			val, err := readInt(br)
			if err != nil {
				return nil, err
			}
			if val > maxVal {
				return nil, fmt.Errorf("pnm: value %d is greater than max value %d", val, maxVal)
			}
			img.Set(x, y, graymap.Color{Y: uint16(val), Max: graymap.Model(maxVal)})
		}
	}
	return img, nil
}

// decodeP3 decodes a plain Portable Pix Map image.
// See https://netpbm.sourceforge.net/doc/ppm.html
func decodeP3(br *bufio.Reader, c config) (image.Image, error) {
	if err := skipWhitespace(br); err != nil {
		return nil, err
	}
	maxVal, err := readInt(br)
	if err != nil {
		return nil, err
	}
	if maxVal == 0 || maxVal > 0xffff {
		return nil, fmt.Errorf("pnm: unsupported max value: %d", maxVal)
	}

	img := pixmap.New(image.Rect(0, 0, c.Width, c.Height), pixmap.Model(maxVal))
	for y := 0; y < c.Height; y++ {
		for x := 0; x < c.Width; x++ {
			if err := skipWhitespace(br); err != nil {
				return nil, err
			}
			r, err := readInt(br)
			if err != nil {
				return nil, err
			}
			if r > maxVal {
				return nil, fmt.Errorf("pnm: value %d is greater than max value %d", r, maxVal)
			}
			if err := skipWhitespace(br); err != nil {
				return nil, err
			}
			g, err := readInt(br)
			if err != nil {
				return nil, err
			}
			if g > maxVal {
				return nil, fmt.Errorf("pnm: value %d is greater than max value %d", g, maxVal)
			}
			if err := skipWhitespace(br); err != nil {
				return nil, err
			}
			b, err := readInt(br)
			if err != nil {
				return nil, err
			}
			if b > maxVal {
				return nil, fmt.Errorf("pnm: value %d is greater than max value %d", b, maxVal)
			}
			img.Set(x, y, pixmap.Color{
				R:   uint16(r),
				G:   uint16(g),
				B:   uint16(b),
				Max: pixmap.Model(maxVal),
			})
		}
	}
	return img, nil
}

// decodeP4 decodes a raw Portable Bit Map image.
// See https://netpbm.sourceforge.net/doc/pbm.html
func decodeP4(br *bufio.Reader, c config) (image.Image, error) {
	if err := skipOneWhitespace(br); err != nil {
		return nil, err
	}

	stride := (c.Width + 7) / 8
	buf := make([]byte, stride*c.Height)
	if _, err := io.ReadFull(br, buf); err != nil {
		return nil, err
	}

	return &bitmap.Image{
		Pix:    buf,
		Stride: stride,
		Rect:   image.Rect(0, 0, c.Width, c.Height),
	}, nil
}

// decodeP5 decodes a raw Portable Gray Map image.
// See https://netpbm.sourceforge.net/doc/pgm.html
func decodeP5(br *bufio.Reader, c config) (image.Image, error) {
	if err := skipWhitespace(br); err != nil {
		return nil, err
	}
	maxVal, err := readInt(br)
	if err != nil {
		return nil, err
	}
	if maxVal == 0 || maxVal > 0xffff {
		return nil, fmt.Errorf("pnm: unsupported max value: %d", maxVal)
	}

	stride := c.Width
	if maxVal >= 256 {
		stride *= 2
	}
	buf := make([]byte, stride*c.Height)
	if _, err := io.ReadFull(br, buf); err != nil {
		return nil, err
	}

	return &graymap.Image{
		Pix:    buf,
		Stride: stride,
		Rect:   image.Rect(0, 0, c.Width, c.Height),
		Max:    graymap.Model(maxVal),
	}, nil
}

// decodeP6 decodes a raw Portable Pix Map image.
// See https://netpbm.sourceforge.net/doc/ppm.html
func decodeP6(br *bufio.Reader, c config) (image.Image, error) {
	if err := skipWhitespace(br); err != nil {
		return nil, err
	}
	maxVal, err := readInt(br)
	if err != nil {
		return nil, err
	}
	if maxVal == 0 || maxVal > 0xffff {
		return nil, fmt.Errorf("pnm: unsupported max value: %d", maxVal)
	}

	if err := skipOneWhitespace(br); err != nil {
		return nil, err
	}

	stride := c.Width * 3
	if maxVal >= 256 {
		stride *= 2
	}
	buf := make([]byte, stride*c.Height)
	if _, err := io.ReadFull(br, buf); err != nil {
		return nil, err
	}

	return &pixmap.Image{
		Pix:    buf,
		Stride: stride,
		Rect:   image.Rect(0, 0, c.Width, c.Height),
		Max:    pixmap.Model(maxVal),
	}, nil
}

func DecodeConfig(r io.Reader) (image.Config, error) {
	br := bufio.NewReader(r)
	c, err := decodeConfig(br)
	if err != nil {
		return image.Config{}, err
	}

	var m color.Model
	switch c.MagicNumber {
	case 0x5031, 0x5034: // P1, P4
		m = bitmap.ColorModel
	case 0x5032, 0x5035: // P2, P5
		if err := skipWhitespace(br); err != nil {
			return image.Config{}, err
		}
		maxVal, err := readInt(br)
		if err != nil {
			return image.Config{}, err
		}
		if maxVal == 0 || maxVal > 0xffff {
			return image.Config{}, fmt.Errorf("pnm: unsupported max value: %d", maxVal)
		}
		m = graymap.Model(maxVal)
	case 0x5033, 0x5036: // P3, P6
		if err := skipWhitespace(br); err != nil {
			return image.Config{}, err
		}
		maxVal, err := readInt(br)
		if err != nil {
			return image.Config{}, err
		}
		if maxVal == 0 || maxVal > 0xffff {
			return image.Config{}, fmt.Errorf("pnm: unsupported max value: %d", maxVal)
		}
		m = pixmap.Model(maxVal)
	default:
		return image.Config{}, errors.New("pnm: unsupported format")
	}
	return image.Config{
		ColorModel: m,
		Width:      c.Width,
		Height:     c.Height,
	}, nil
}

func decodeConfig(br *bufio.Reader) (config, error) {
	b1, err := br.ReadByte()
	if err != nil {
		return config{}, nil
	}
	b2, err := br.ReadByte()
	if err != nil {
		return config{}, nil
	}
	magic := uint16(b1)<<8 | uint16(b2)

	// read width
	if err := skipWhitespace(br); err != nil {
		return config{}, err
	}
	width, err := readInt(br)
	if err != nil {
		return config{}, err
	}
	if width <= 0 {
		return config{}, fmt.Errorf("pnm: invalid width: %d", width)
	}

	// read height
	if err := skipWhitespace(br); err != nil {
		return config{}, err
	}
	height, err := readInt(br)
	if err != nil {
		return config{}, err
	}
	if height <= 0 {
		return config{}, fmt.Errorf("pnm: invalid height: %d", height)
	}

	hi, lo := bits.Mul64(uint64(width), uint64(height))
	if hi != 0 || lo > 1<<30 {
		return config{}, errors.New("pnm: too large image")
	}

	return config{
		MagicNumber: magic,
		Width:       width,
		Height:      height,
	}, nil
}

func skipOneWhitespace(br *bufio.Reader) error {
	b, err := br.ReadByte()
	if err != nil {
		return err
	}
	switch b {
	case ' ', '\t', '\r', '\n', '\v', '\f':
		return nil
	default:
		return fmt.Errorf("pnm: unexpected char: %c", b)
	}
}

func skipWhitespace(br *bufio.Reader) error {
	for {
		b, err := br.ReadByte()
		if err != nil {
			return err
		}
		switch b {
		case ' ', '\t', '\r', '\n', '\v', '\f':
			// ignore whitespace
		case '#':
			// comment, skip to '\r' or '\n'
			if err := skipComment(br); err != nil {
				return err
			}
		default:
			if err := br.UnreadByte(); err != nil {
				return err
			}
			return nil
		}
	}
}

func skipComment(br *bufio.Reader) error {
	for {
		b, err := br.ReadByte()
		if err != nil {
			return err
		}
		if b == '\r' || b == '\n' {
			return nil
		}
	}
}

func readInt(br *bufio.Reader) (int, error) {
	const cutoff = math.MaxInt/10 + 1
	var ret int
	b, err := br.ReadByte()
	if err != nil {
		return 0, err
	}
	if '0' <= b && b <= '9' {
		ret = int(b - '0')
	} else {
		return 0, fmt.Errorf("pnm: unexpected char: %c", b)
	}

	for {
		b, err := br.ReadByte()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return 0, err
		}
		if '0' <= b && b <= '9' {
			if ret >= cutoff {
				return math.MaxInt, errors.New("pnm: integer overflow")
			}
			ret = ret*10 + int(b-'0')
		} else {
			if err := br.UnreadByte(); err != nil {
				return 0, err
			}
			break
		}
	}
	return ret, nil
}
