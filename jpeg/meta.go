package jpeg

import (
	"image"
	"io"

	"github.com/shogo82148/go-imaging/icc"
)

type ImageWithMeta struct {
	image.Image
	ICCProfile *icc.Profile
}

func DecodeWithMeta(r io.Reader) (*ImageWithMeta, error) {
	var d decoder
	return d.decodeWithMeta(r)
}

func (d *decoder) decodeWithMeta(r io.Reader) (*ImageWithMeta, error) {
	d.r = r

	// Check for the Start Of Image marker.
	if err := d.readFull(d.tmp[:2]); err != nil {
		return nil, err
	}
	if d.tmp[0] != 0xff || d.tmp[1] != soiMarker {
		return nil, FormatError("missing SOI marker")
	}

	// Process the remaining segments until the End Of Image marker.
	for {
		err := d.readFull(d.tmp[:2])
		if err != nil {
			return nil, err
		}
		for d.tmp[0] != 0xff {
			// Strictly speaking, this is a format error. However, libjpeg is
			// liberal in what it accepts. As of version 9, next_marker in
			// jdmarker.c treats this as a warning (JWRN_EXTRANEOUS_DATA) and
			// continues to decode the stream. Even before next_marker sees
			// extraneous data, jpeg_fill_bit_buffer in jdhuff.c reads as many
			// bytes as it can, possibly past the end of a scan's data. It
			// effectively puts back any markers that it overscanned (e.g. an
			// "\xff\xd9" EOI marker), but it does not put back non-marker data,
			// and thus it can silently ignore a small number of extraneous
			// non-marker bytes before next_marker has a chance to see them (and
			// print a warning).
			//
			// We are therefore also liberal in what we accept. Extraneous data
			// is silently ignored.
			//
			// This is similar to, but not exactly the same as, the restart
			// mechanism within a scan (the RST[0-7] markers).
			//
			// Note that extraneous 0xff bytes in e.g. SOS data are escaped as
			// "\xff\x00", and so are detected a little further down below.
			d.tmp[0] = d.tmp[1]
			d.tmp[1], err = d.readByte()
			if err != nil {
				return nil, err
			}
		}
		marker := d.tmp[1]
		if marker == 0 {
			// Treat "\xff\x00" as extraneous data.
			continue
		}
		for marker == 0xff {
			// Section B.1.1.2 says, "Any marker may optionally be preceded by any
			// number of fill bytes, which are bytes assigned code X'FF'".
			marker, err = d.readByte()
			if err != nil {
				return nil, err
			}
		}
		if marker == eoiMarker { // End Of Image.
			break
		}
		if rst0Marker <= marker && marker <= rst7Marker {
			// Figures B.2 and B.16 of the specification suggest that restart markers should
			// only occur between Entropy Coded Segments and not after the final ECS.
			// However, some encoders may generate incorrect JPEGs with a final restart
			// marker. That restart marker will be seen here instead of inside the processSOS
			// method, and is ignored as a harmless error. Restart markers have no extra data,
			// so we check for this before we read the 16-bit length of the segment.
			continue
		}

		// Read the 16-bit length of the segment. The value includes the 2 bytes for the
		// length itself, so we subtract 2 to get the number of remaining bytes.
		if err = d.readFull(d.tmp[:2]); err != nil {
			return nil, err
		}
		n := int(d.tmp[0])<<8 + int(d.tmp[1]) - 2
		if n < 0 {
			return nil, FormatError("short segment length")
		}

		switch marker {
		case sof0Marker, sof1Marker, sof2Marker:
			d.baseline = marker == sof0Marker
			d.progressive = marker == sof2Marker
			err = d.processSOF(n)
		case dhtMarker:
			err = d.processDHT(n)
		case dqtMarker:
			err = d.processDQT(n)
		case sosMarker:
			err = d.processSOS(n)
		case driMarker:
			err = d.processDRI(n)
		case app0Marker:
			err = d.processApp0Marker(n)
		case app2Marker:
			err = d.processApp2Marker(n)
		case app14Marker:
			err = d.processApp14Marker(n)
		default:
			if app0Marker <= marker && marker <= app15Marker || marker == comMarker {
				err = d.ignore(n)
			} else if marker < 0xc0 { // See Table B.1 "Marker code assignments".
				err = FormatError("unknown marker")
			} else {
				err = UnsupportedError("unknown marker")
			}
		}
		if err != nil {
			return nil, err
		}
	}

	if d.progressive {
		if err := d.reconstructProgressiveImage(); err != nil {
			return nil, err
		}
	}

	var img image.Image
	var err error
	if d.img1 != nil {
		img = d.img1
	} else if d.img3 != nil {
		if d.blackPix != nil {
			img, err = d.applyBlack()
		} else if d.isRGB() {
			img, err = d.convertToRGB()
		} else {
			img = d.img3
		}
	}
	if err != nil {
		return nil, err
	}
	if img == nil {
		return nil, FormatError("missing SOS marker")
	}

	var iccProfile *icc.Profile
	if d.iccProfileLen > 0 {
		r := &multiBlockReader{blocks: d.iccProfile[:]}
		iccProfile, err = icc.Decode(r)
		if err != nil {
			return nil, err
		}
	}
	return &ImageWithMeta{
		Image:      img,
		ICCProfile: iccProfile,
	}, nil
}

func (d *decoder) processApp2Marker(n int) error {
	l := len("ICC_PROFILE") + 1 + 2 // +1 for the null terminator, +2 for sub-block index and total sub-blocks
	if n < l {
		return d.ignore(n)
	}

	if err := d.readFull(d.tmp[:l]); err != nil {
		return err
	}
	if string(d.tmp[:l-2]) != "ICC_PROFILE\x00" {
		return d.ignore(n - l)
	}
	buf := make([]byte, n-l)
	if err := d.readFull(buf); err != nil {
		return err
	}
	idx := int(d.tmp[l-2])
	d.iccProfile[idx] = buf
	d.iccProfileLen += len(buf)
	return nil
}

type multiBlockReader struct {
	blocks [][]byte
	idx    int // current block index
	off    int // current offset in the block
}

func (r *multiBlockReader) Read(p []byte) (n int, err error) {
	if r.idx >= len(r.blocks) {
		return 0, io.EOF
	}
	n = copy(p, r.blocks[r.idx][r.off:])
	r.off += n
	if r.off >= len(r.blocks[r.idx]) {
		r.idx++
		r.off = 0
	}
	return n, nil
}
