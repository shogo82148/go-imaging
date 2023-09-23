package png

import (
	"hash/crc32"
	"image"
	"io"
)

// ImageWithMeta is a PNG image with metadata.
type ImageWithMeta struct {
	image.Image

	// Gamma is the gamma value of the image.
	// If Gamma is 0, the image has no gamma information.
	Gamma float64
}

// Decode reads a PNG image from r and returns it as an image.Image.
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
