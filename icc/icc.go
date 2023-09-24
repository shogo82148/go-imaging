package icc

import (
	"bytes"
	"encoding"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"
	"slices"
	"strconv"
	"time"
)

const iccHeaderSize = 128
const ICCMagicNumber Signature = 0x61637370 // 'acsp'

type Signature uint32

func printable(b byte) byte {
	if b < 0x20 || b > 0x7e {
		return '.'
	}
	return b
}

func (s Signature) String() string {
	var buf [4]byte
	buf[0] = printable(byte(s >> 24))
	buf[1] = printable(byte(s >> 16))
	buf[2] = printable(byte(s >> 8))
	buf[3] = printable(byte(s))
	return string(buf[:])
}

type DateTimeNumber struct {
	Year   uint16
	Month  uint16
	Day    uint16
	Hour   uint16
	Minute uint16
	Second uint16
}

func (dt DateTimeNumber) Time() time.Time {
	return time.Date(
		int(dt.Year),
		time.Month(dt.Month),
		int(dt.Day),
		int(dt.Hour),
		int(dt.Minute),
		int(dt.Second),
		0,
		time.UTC,
	)
}

type positionNumber struct {
	Offset uint32
	Size   uint32
}

type response16Number struct {
	Device   uint16
	Reserved uint16
	Attr     S15Fixed16Number
}

// S15Fixed16Number is a 16-bit signed integer with a 16-bit fraction.
type S15Fixed16Number int32

func (n S15Fixed16Number) Float64() float64 {
	return float64(n) / 0x10000
}

func S15Fixed16NumberFromFloat64(f float64) S15Fixed16Number {
	return S15Fixed16Number(math.RoundToEven(f * 0x10000))
}

func (n S15Fixed16Number) String() string {
	f := n.Float64()
	return strconv.FormatFloat(f, 'f', -1, 64)
}

// u16Fix16Number is a 16-bit unsigned integer with a 16-bit fraction.
type u16Fix16Number uint32

// u1Fixed15Number is a 1-bit unsigned integer with a 15-bit fraction.
type u1Fixed15Number uint16

// U8Fixed8Number is a 8-bit unsigned integer with a 8-bit fraction.
type U8Fixed8Number uint16

func (n U8Fixed8Number) Float64() float64 {
	return float64(n) / 0x100
}

type XYZNumber struct {
	X S15Fixed16Number
	Y S15Fixed16Number
	Z S15Fixed16Number
}

type Profile struct {
	ProfileHeader
	Tags []TagEntry
}

type Version uint32

func (v Version) Major() int {
	return int(v >> 24)
}

func (v Version) Minor() int {
	return int(v >> 20 & 0xf)
}

func (v Version) BugFix() int {
	return int(v >> 16 & 0xf)
}

func (v Version) String() string {
	return fmt.Sprintf("%d.%d.%d.0", v.Major(), v.Minor(), v.BugFix())
}

type TagEntry struct {
	Tag        Tag
	TagContent TagContent
}

// readN reads n bytes from r.
func readN(r io.Reader, n int64) ([]byte, error) {
	if n < 0 {
		return nil, errors.New("icc: invalid size")
	}
	if n == 0 {
		return []byte{}, nil
	}
	if n < 16*1024 {
		// optimize for small size
		buf := make([]byte, n)
		if _, err := io.ReadFull(r, buf); err != nil && err != io.EOF {
			return nil, err
		}
		return buf, nil
	}

	data, err := io.ReadAll(io.LimitReader(r, n))
	if err != nil {
		return nil, err
	}
	if int64(len(data)) < n {
		return nil, io.ErrUnexpectedEOF
	}
	return data, nil
}

func Decode(r io.Reader) (*Profile, error) {
	var header ProfileHeader
	if err := binary.Read(r, binary.BigEndian, &header); err != nil {
		return nil, err
	}
	if header.Magic != ICCMagicNumber {
		return nil, errors.New("icc: invalid magic number")
	}
	if header.Size < iccHeaderSize {
		return nil, errors.New("icc: invalid profile size")
	}

	data, err := readN(r, int64(header.Size)-iccHeaderSize)
	if err != nil {
		return nil, err
	}
	br := bytes.NewReader(data)

	var tagCount uint32
	if err := binary.Read(br, binary.BigEndian, &tagCount); err != nil {
		return nil, err
	}
	table := make([]tagTable, tagCount)
	if err := binary.Read(br, binary.BigEndian, &table); err != nil {
		return nil, err
	}

	tags := make([]TagEntry, tagCount)
	for i, t := range table {
		if t.Offset+t.Size > header.Size {
			return nil, errors.New("icc: invalid tag table")
		}
		if t.Offset > t.Offset+t.Size {
			// overflow
			return nil, errors.New("icc: invalid tag table")
		}

		tagData := data[t.Offset-iccHeaderSize : t.Offset+t.Size-iccHeaderSize]
		tagType := TagType(binary.BigEndian.Uint32(tagData))

		var content TagContent
		switch tagType {
		case TagTypeCurve:
			var tag TagContentCurve
			if err := tag.UnmarshalBinary(tagData); err != nil {
				return nil, err
			}
			content = &tag
		case TagTypeParametricCurve:
			var tag TagContentParametricCurve
			if err := tag.UnmarshalBinary(tagData); err != nil {
				return nil, err
			}
			content = &tag
		default:
			content = &TagContentRaw{
				Data: slices.Clone(tagData),
			}
		}
		tags[i] = TagEntry{
			Tag:        t.Signature,
			TagContent: content,
		}
	}

	return &Profile{
		ProfileHeader: header,
		Tags:          tags,
	}, nil
}

// alignWriter is a writer that aligns the data to 4 bytes.
type alignWriter struct {
	w   io.Writer
	n   int64
	buf [4]byte
}

func (w *alignWriter) Write(data []byte) (n int, err error) {
	n, err = w.w.Write(data)
	w.n += int64(n)
	return
}

func (w *alignWriter) Align() error {
	n := (w.n + 0x03) &^ 0x03
	_, err := w.Write(w.buf[:n-w.n])
	return err
}

func (p *Profile) Encode(w io.Writer) error {
	// calculate the size of the profile header
	offset := uint32(128)              // for the profile header
	offset += 4                        // for the tag count
	offset += uint32(len(p.Tags) * 12) // for the tag table

	tagTable := make([]tagTable, len(p.Tags))
	tagContents := make([][]byte, len(p.Tags))
	for i, tag := range p.Tags {
		// encode the tag content
		data, err := tag.TagContent.MarshalBinary()
		if err != nil {
			return err
		}
		tagContents[i] = data

		// set the tag table
		tagTable[i].Signature = tag.Tag
		tagTable[i].Offset = offset
		tagTable[i].Size = uint32(len(data))

		offset += uint32(len(data))
		offset = (offset + 0x03) &^ 0x03 // align to 4 bytes
	}

	header := p.ProfileHeader
	header.Size = offset
	header.Magic = ICCMagicNumber

	aw := &alignWriter{w: w}
	if err := binary.Write(aw, binary.BigEndian, header); err != nil {
		return err
	}
	if err := binary.Write(aw, binary.BigEndian, uint32(len(p.Tags))); err != nil {
		return err
	}
	if err := binary.Write(aw, binary.BigEndian, tagTable); err != nil {
		return err
	}
	for _, data := range tagContents {
		if _, err := aw.Write(data); err != nil {
			return err
		}
		if err := aw.Align(); err != nil {
			return err
		}
	}
	return nil
}

func (p *Profile) Get(tag Tag) TagContent {
	for _, t := range p.Tags {
		if t.Tag == tag {
			return t.TagContent
		}
	}
	return nil
}

// ProfileHeader is the header of ICC profile.
type ProfileHeader struct {
	Size                   uint32
	CMMType                Signature
	Version                Version
	Class                  Class
	ColorSpace             ColorSpace
	ProfileConnectionSpace ColorSpace
	DateTime               DateTimeNumber
	Magic                  Signature
	Platform               Platform
	Flags                  uint32
	Manufacturer           Signature
	DeviceModel            Signature
	DeviceAttributes       uint64
	RenderingIntent        uint32
	XYZ                    XYZNumber
	ProfileCreator         uint32
	ProfileID              [2]uint64
	SpectralPCS            uint32
	SpectralPCSRange       [6]byte
	BiSpectralPCSRange     [6]byte
	MCSSignature           uint32
	SubClass               uint32
	Reserved               uint32
}

type tagTable struct {
	Signature Tag
	Offset    uint32
	Size      uint32
}

type TagContent interface {
	encoding.BinaryMarshaler
	encoding.BinaryUnmarshaler
	TagType() TagType
}

var _ TagContent = (*TagContentRaw)(nil)
var _ Curve = (*TagContentCurve)(nil)

type TagContentRaw struct {
	Data []byte
}

func (t *TagContentRaw) TagType() TagType {
	return TagType(binary.BigEndian.Uint32(t.Data))
}

func (t *TagContentRaw) MarshalBinary() ([]byte, error) {
	return t.Data, nil
}

func (t *TagContentRaw) UnmarshalBinary(data []byte) error {
	t.Data = slices.Clone(data)
	return nil
}
