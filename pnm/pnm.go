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

	"github.com/shogo82148/go-imaging/bitmap"
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
	return decodeP4(br, c)
}

func decodeP1(r io.Reader, c config) (image.Image, error) {
	return nil, nil
}

// decodeP4 decodes a Portable Bit Map image.
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

	// read height
	if err := skipWhitespace(br); err != nil {
		return config{}, err
	}
	height, err := readInt(br)
	if err != nil {
		return config{}, err
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
	case ' ', '\t', '\r', '\n':
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
		case ' ', '\t', '\r', '\n':
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
