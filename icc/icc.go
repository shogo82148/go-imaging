package icc

import (
	"bytes"
	"encoding"
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"slices"
	"strconv"
	"time"
)

const ICCmagicNumber = 0x61637370 // 'acsp'

type dateTimeNumber struct {
	Year   uint16
	Month  uint16
	Day    uint16
	Hour   uint16
	Minute uint16
	Second uint16
}

func (dt dateTimeNumber) Time() time.Time {
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

type xyzNumber struct {
	X S15Fixed16Number
	Y S15Fixed16Number
	Z S15Fixed16Number
}

type Profile struct {
	Version    Version
	Class      Class
	ColorSpace ColorSpace
	Time       time.Time
	Tags       []TagEntry
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

type Class uint32

const (
	ClassInput      Class = 0x73636e72 // 'scnr'
	ClassDisplay    Class = 0x6d6e7472 // 'mntr'
	ClassOutput     Class = 0x70727472 // 'prtr'
	ClassLink       Class = 0x6c696e6b // 'link'
	ClassAbstract   Class = 0x61627374 // 'abst'
	ClassColorSpace Class = 0x73706163 // 'spac'
	ClassNamedColor Class = 0x6e6d636c // 'nmcl'
)

func printable(b byte) byte {
	if b < 0x20 || b > 0x7e {
		return '.'
	}
	return b
}

func (class Class) String() string {
	switch class {
	case ClassInput:
		return "Input device profile"
	case ClassDisplay:
		return "Display device profile"
	case ClassOutput:
		return "Output device profile"
	case ClassLink:
		return "DeviceLink profile"
	case ClassAbstract:
		return "Abstract profile"
	case ClassColorSpace:
		return "ColorSpace profile"
	case ClassNamedColor:
		return "NamedColor profile"
	default:
		return fmt.Sprintf(
			"Unknown Class(%08xh '%c%c%c%c')",
			uint32(class),
			printable(byte(class>>24)),
			printable(byte(class>>16)),
			printable(byte(class>>8)),
			printable(byte(class)),
		)
	}
}

type ColorSpace uint32

const (
	ColorSpaceXYZ   ColorSpace = 0x58595a20 // 'XYZ '
	ColorSpaceLab   ColorSpace = 0x4c616220 // 'Lab '
	ColorSpaceLuv   ColorSpace = 0x4c757620 // 'Luv '
	ColorSpaceYCbCr ColorSpace = 0x59436272 // 'YCbr'
	ColorSpaceYxy   ColorSpace = 0x59787920 // 'Yxy '
	ColorSpaceRGB   ColorSpace = 0x52474220 // 'RGB '
	ColorSpaceGray  ColorSpace = 0x47524159 // 'GRAY'
	ColorSpaceHSV   ColorSpace = 0x48535620 // 'HSV '
	ColorSpaceHLS   ColorSpace = 0x484c5320 // 'HLS '
	ColorSpaceCMYK  ColorSpace = 0x434d594b // 'CMYK'
	ColorSpaceCMY   ColorSpace = 0x434d5920 // 'CMY '
	ColorSpace2CLR  ColorSpace = 0x32434c52 // '2CLR'
	ColorSpace3CLR  ColorSpace = 0x33434c52 // '3CLR'
	ColorSpace4CLR  ColorSpace = 0x34434c52 // '4CLR'
	ColorSpace5CLR  ColorSpace = 0x35434c52 // '5CLR'
	ColorSpace6CLR  ColorSpace = 0x36434c52 // '6CLR'
	ColorSpace7CLR  ColorSpace = 0x37434c52 // '7CLR'
	ColorSpace8CLR  ColorSpace = 0x38434c52 // '8CLR'
	ColorSpace9CLR  ColorSpace = 0x39434c52 // '9CLR'
	ColorSpace10CLR ColorSpace = 0x41434c52 // 'ACLR'
	ColorSpace11CLR ColorSpace = 0x42434c52 // 'BCLR'
	ColorSpace12CLR ColorSpace = 0x43434c52 // 'CCLR'
	ColorSpace13CLR ColorSpace = 0x44434c52 // 'DCLR'
	ColorSpace14CLR ColorSpace = 0x45434c52 // 'ECLR'
	ColorSpace15CLR ColorSpace = 0x46434c52 // 'FCLR'
)

func (cs ColorSpace) String() string {
	switch cs {
	case ColorSpaceXYZ:
		return "CIEXYZ"
	case ColorSpaceLab:
		return "CIELab"
	case ColorSpaceLuv:
		return "CIELuv"
	case ColorSpaceYCbCr:
		return "YCbCr"
	case ColorSpaceYxy:
		return "CIEYxy"
	case ColorSpaceRGB:
		return "RGB"
	case ColorSpaceGray:
		return "Gray"
	case ColorSpaceHSV:
		return "HSV"
	case ColorSpaceHLS:
		return "HLS"
	case ColorSpaceCMYK:
		return "CMYK"
	case ColorSpaceCMY:
		return "CMY"
	case ColorSpace2CLR:
		return "2 color"
	case ColorSpace3CLR:
		return "3 color"
	case ColorSpace4CLR:
		return "4 color"
	case ColorSpace5CLR:
		return "5 color"
	case ColorSpace6CLR:
		return "6 color"
	case ColorSpace7CLR:
		return "7 color"
	case ColorSpace8CLR:
		return "8 color"
	case ColorSpace9CLR:
		return "9 color"
	case ColorSpace10CLR:
		return "10 color"
	case ColorSpace11CLR:
		return "11 color"
	case ColorSpace12CLR:
		return "12 color"
	case ColorSpace13CLR:
		return "13 color"
	case ColorSpace14CLR:
		return "14 color"
	case ColorSpace15CLR:
		return "15 color"
	default:
		return fmt.Sprintf(
			"Unknown(%08xh '%c%c%c%c')",
			uint32(cs),
			printable(byte(cs>>24)),
			printable(byte(cs>>16)),
			printable(byte(cs>>8)),
			printable(byte(cs)),
		)
	}
}

type TagEntry struct {
	Tag        Tag
	TagContent TagContent
}

func Decode(data []byte) (*Profile, error) {
	r := bytes.NewReader(data)
	var header profileHeader
	if err := binary.Read(r, binary.BigEndian, &header); err != nil {
		return nil, err
	}
	if header.Magic != ICCmagicNumber {
		return nil, errors.New("icc: invalid magic number")
	}

	var tagCount uint32
	if err := binary.Read(r, binary.BigEndian, &tagCount); err != nil {
		return nil, err
	}
	table := make([]tagTable, tagCount)
	if err := binary.Read(r, binary.BigEndian, &table); err != nil {
		return nil, err
	}

	tags := make([]TagEntry, tagCount)
	for i, t := range table {
		tagData := data[t.Offset : t.Offset+t.Size]
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
			content = &TagContentRaw{Data: tagData}
		}
		tags[i] = TagEntry{
			Tag:        t.Signature,
			TagContent: content,
		}
	}

	return &Profile{
		Class:      header.Class,
		ColorSpace: header.ColorSpace,
		Time:       header.DateTime.Time(),
		Tags:       tags,
	}, nil
}

func (p *Profile) Get(tag Tag) TagContent {
	for _, t := range p.Tags {
		if t.Tag == tag {
			return t.TagContent
		}
	}
	return nil
}

// profileHeader is the header of ICC profile.
type profileHeader struct {
	Size               uint32
	CMMType            uint32
	Version            Version
	Class              Class
	ColorSpace         ColorSpace
	PCS                uint32
	DateTime           dateTimeNumber
	Magic              uint32
	Platform           uint32
	Flags              uint32
	Manufacturer       uint32
	DeviceModel        uint32
	DeviceAttributes   uint64
	RenderingIntent    uint32
	XYZ                xyzNumber
	ProfileCreator     uint32
	ProfileID          [2]uint64
	SpectralPCS        uint32
	SpectralPCSRange   [6]byte
	BiSpectralPCSRange [6]byte
	MCSSignature       uint32
	SubClass           uint32
	Reserved           uint32
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

// Curve is a tone reproduction curve.
type Curve interface {
	EncodeTone(x float64) float64
	DecodeTone(x float64) float64
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

var _ TagContent = (*TagContentCurve)(nil)

type TagContentCurve struct {
	Data []uint16
}

type tagContentCurve struct {
	TagType  TagType
	Reserved uint32
	Count    uint32
}

func (t *TagContentCurve) TagType() TagType { return TagTypeCurve }

func (t *TagContentCurve) MarshalBinary() ([]byte, error) {
	buf := bytes.NewBuffer(make([]byte, 0, 12+len(t.Data)*2))
	curve := tagContentCurve{
		TagType: t.TagType(),
		Count:   uint32(len(t.Data)),
	}
	if err := binary.Write(buf, binary.BigEndian, curve); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.BigEndian, t.Data); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (t *TagContentCurve) UnmarshalBinary(data []byte) error {
	r := bytes.NewReader(data)
	var curve tagContentCurve
	if err := binary.Read(r, binary.BigEndian, &curve); err != nil {
		return err
	}
	t.Data = make([]uint16, curve.Count)
	if err := binary.Read(r, binary.BigEndian, &t.Data); err != nil {
		return err
	}
	return nil
}

func (t *TagContentCurve) EncodeTone(x float64) float64 {
	x = max(0, min(1, x)) // clip to [0.0, 1.0]
	if len(t.Data) == 0 {
		return x
	}
	if len(t.Data) == 1 {
		gamma := U8Fixed8Number(t.Data[0]).Float64()
		return math.Pow(x, gamma)
	}

	i, f := math.Modf(x * float64(len(t.Data)-1))
	i0 := int(i)
	y0 := float64(t.Data[i0]) / 0xffff
	if i0 == len(t.Data)-1 {
		return y0
	}

	// linear interpolation
	i1 := i0 + 1
	y1 := float64(t.Data[i1]) / 0xffff
	return y0 + f*(y1-y0)
}

func (t *TagContentCurve) DecodeTone(y float64) float64 {
	y = max(0, min(1, y)) // clip to [0.0, 1.0]
	if len(t.Data) == 0 {
		return y
	}
	if len(t.Data) == 1 {
		gamma := U8Fixed8Number(t.Data[0]).Float64()
		return math.Pow(y, 1/gamma)
	}

	i := uint16(y * 0xffff)
	idx, ok := slices.BinarySearch(t.Data, i)
	x0 := float64(idx) / float64(len(t.Data)-1)
	if ok || idx == len(t.Data)-1 {
		return x0
	}
	y0 := float64(t.Data[idx]) / 0xffff
	y1 := float64(t.Data[idx+1]) / 0xffff
	if y0 == y1 {
		return x0
	}
	f := (y - y0) / (y1 - y0)
	return x0 + f/float64(len(t.Data)-1)
}

type TagContentParametricCurve struct {
	FunctionType uint16
	Params       [8]S15Fixed16Number // this is not a slice because to avoid extra boundary check.
}

type tagContentParametricCurve struct {
	TagType      TagType
	_            uint32
	FunctionType uint16
	_            uint16
}

func (t *TagContentParametricCurve) TagType() TagType { return TagTypeParametricCurve }

func (t *TagContentParametricCurve) MarshalBinary() ([]byte, error) {
	return nil, errors.New("not implemented")
}

func (t *TagContentParametricCurve) UnmarshalBinary(data []byte) error {
	r := bytes.NewReader(data)
	var curve tagContentParametricCurve
	if err := binary.Read(r, binary.BigEndian, &curve); err != nil {
		return err
	}
	var params []S15Fixed16Number
	switch curve.FunctionType {
	case 0x0000:
		params = t.Params[:1]
	case 0x0001:
		params = t.Params[:3]
	case 0x0002:
		params = t.Params[:4]
	case 0x0003:
		params = t.Params[:5]
	case 0x0004:
		params = t.Params[:7]
	}
	if err := binary.Read(r, binary.BigEndian, params); err != nil {
		return err
	}
	return nil
}

func (t *TagContentParametricCurve) EncodeTone(x float64) float64 {
	x = max(0, min(1, x)) // clip to [0.0, 1.0]
	switch t.FunctionType {
	// Y = X^g
	case 0x0000:
		g := t.Params[0].Float64()
		return math.Pow(x, g)

	// CIE122-1966
	// Y = (aX + b)^g   if X >= -b/a
	// Y = 0            if X <  -b/a
	case 0x0001:
		g := t.Params[0].Float64()
		a := t.Params[1].Float64()
		b := t.Params[2].Float64()

		y := a*x + b
		if y < 0 {
			return 0
		}
		y = math.Pow(y, g)
		return y

	// IEC 61966‐3
	// Y = (aX + b)^g + c  if X >= -b/a
	// Y = c               if X <  -b/a
	case 0x0002:
		g := t.Params[0].Float64()
		a := t.Params[1].Float64()
		b := t.Params[2].Float64()
		c := t.Params[3].Float64()

		y := a*x + b
		if y < 0 {
			return c
		}
		y = math.Pow(y, g) + c
		return max(0, min(1, y)) // clip to [0.0, 1.0]

	// Y = (aX + b)^g     if X >= d
	// Y = cX             if X <  d
	case 0x0003:
		g := t.Params[0].Float64()
		a := t.Params[1].Float64()
		b := t.Params[2].Float64()
		c := t.Params[3].Float64()
		d := t.Params[4].Float64()

		if x < d {
			return c * x
		}
		y := math.Pow(a*x+b, g)
		return y

	// Y = (aX + b)^g + e  if X >= d
	// Y = cX + f          if X <  d
	case 0x0004:
		g := t.Params[0].Float64()
		a := t.Params[1].Float64()
		b := t.Params[2].Float64()
		c := t.Params[3].Float64()
		d := t.Params[4].Float64()
		e := t.Params[5].Float64()
		f := t.Params[6].Float64()

		if x < d {
			return c*x + f
		}
		y := math.Pow(a*x+b, g) + e
		return max(0, min(1, y)) // clip to [0.0, 1.0]
	}
	return x
}

func (t *TagContentParametricCurve) DecodeTone(y float64) float64 {
	y = max(0, min(1, y)) // clip to [0.0, 1.0]
	switch t.FunctionType {
	// Y = X^g
	case 0x0000:
		g := t.Params[0].Float64()
		return math.Pow(y, 1/g)

	// Y = (aX + b)^g   if X >= -b/a
	// Y = 0            if X <  -b/a
	case 0x0001:
		g := t.Params[0].Float64()
		a := t.Params[1].Float64()
		b := t.Params[2].Float64()

		x := (math.Pow(y, 1/g) - b) / a
		return x

	// Y = (aX + b)^g + c  if X >= -b/a
	// Y = c               if X <  -b/a
	case 0x0002:
		g := t.Params[0].Float64()
		a := t.Params[1].Float64()
		b := t.Params[2].Float64()
		c := t.Params[3].Float64()

		if y < c {
			return -b / a
		}
		y = max(0, y-c)
		x := (math.Pow(y, 1/g) - b) / a
		return max(0, min(1, x)) // clip to [0.0, 1.0]

	// Y = (aX + b)^g     if X >= d
	// Y = cX             if X <  d
	case 0x0003:
		g := t.Params[0].Float64()
		a := t.Params[1].Float64()
		b := t.Params[2].Float64()
		c := t.Params[3].Float64()
		d := t.Params[4].Float64()

		x := (math.Pow(y, 1/g) - b) / a
		if x < d {
			return y / c
		}
		return max(0, min(1, x)) // clip to [0.0, 1.0]

	// Y = (aX + b)^g + e  if X >= d
	// Y = cX + f          if X <  d
	case 0x0004:
		g := t.Params[0].Float64()
		a := t.Params[1].Float64()
		b := t.Params[2].Float64()
		c := t.Params[3].Float64()
		d := t.Params[4].Float64()
		e := t.Params[5].Float64()
		f := t.Params[6].Float64()

		x := (y - f) / c
		if x < d {
			return max(0, min(1, x)) // clip to [0.0, 1.0]
		}
		y = max(0, y-e)
		x = (math.Pow(y, 1/g) - b) / a
		return max(0, min(1, x)) // clip to [0.0, 1.0]
	}
	return y
}
