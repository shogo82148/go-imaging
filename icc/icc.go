package icc

import (
	"bytes"
	"encoding"
	"encoding/binary"
	"log"
	"slices"
)

type dateTimeNumber struct {
	Year   uint16
	Month  uint16
	Day    uint16
	Hour   uint16
	Minute uint16
	Second uint16
}

type positionNumber struct {
	Offset uint32
	Size   uint32
}

type response16Number struct {
	Device   uint16
	Reserved uint16
	Attr     s15Fixed16Number
}

// s15Fixed16Number is a 16-bit signed integer with a 16-bit fraction.
type s15Fixed16Number int32

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
	X s15Fixed16Number
	Y s15Fixed16Number
	Z s15Fixed16Number
}

type Profile struct {
	Tags []TagEntry
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
		log.Printf("%x", tagType)
		switch tagType {
		case TagTypeCurve:
			var tag TagContentCurve
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
		Tags: tags,
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
	Version            uint32
	Class              uint32
	ColorSpace         uint32
	PCS                uint32
	Date               dateTimeNumber
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

type Tag uint32

type tagTable struct {
	Signature Tag
	Offset    uint32
	Size      uint32
}

type TagType uint32

const (
	TagTypeRaw TagType = 0xffffffff

	TagTypeColorantOrder             TagType = 0x636c726f // 'clro'
	TagTypeCurve                     TagType = 0x63757276 // 'curv'
	TagTypeDataType                  TagType = 0x64617461 // 'data'
	TagTypeDateTime                  TagType = 0x6474696d // 'dtim'
	TagTypeDict                      TagType = 0x64637420 // 'dict'
	TagTypeEmbeddedHeightImage       TagType = 0x6568696d // 'ehim'
	TagTypeEmbeddedNormalImage       TagType = 0x656e696d // 'enim'
	TagTypeFloat16Array              TagType = 0x666c3136 // 'fl16'
	TagTypeFloat32Array              TagType = 0x666c3234 // 'fl32'
	TagTypeFloat64Array              TagType = 0x666c3634 // 'fl64'
	TagTypeLutAtoB                   TagType = 0x6d414220 // 'mAB '
	TagTypeLutBtoA                   TagType = 0x6d424120 // 'mBA '
	TagTypeMeasurement               TagType = 0x6d656173 // 'meas'
	TagTypeMultiLocalizedUnicode     TagType = 0x6d6c7563 // 'mluc'
	TagTypeMultiProcessElements      TagType = 0x6d706574 // 'mpet'
	TagTypeParametricCurve           TagType = 0x70617261 // 'para'
	TagTypeS15Fixed16Array           TagType = 0x73663332 // 'sf32'
	TagTypeSignature                 TagType = 0x73696720 // 'sig '
	TagTypeSparseMatrixArray         TagType = 0x736d6174 // 'smat'
	TagTypeSpectralViewingConditions TagType = 0x7376636e // 'svcn'
	TagTypeTagArrayType              TagType = 0x74617279 // 'tary'
	TagTypeTagStruct                 TagType = 0x74737472 // 'tstr'
	TagTypeU16Fixed16Array           TagType = 0x75663332 // 'uf32'
	TagTypeUint16Array               TagType = 0x75693136 // 'ui16'
	TagTypeUint32Array               TagType = 0x75693332 // 'ui32'
	TagTypeUint64Array               TagType = 0x75693634 // 'ui64'
	TagTypeUint8Array                TagType = 0x75693038 // 'ui08'
	TagTypeUTF16                     TagType = 0x75743136 // 'ut16'
	TagTypeUTF8                      TagType = 0x75746638 // 'utf8'
	TagTypeUTF8Zip                   TagType = 0x7a757438 // 'zut8'
	TagTypeXYZ                       TagType = 0x58595a20 // 'XYZ '
	TagTypeZipXML                    TagType = 0x7a786d6c // 'zxml'
)

type TagContent interface {
	encoding.BinaryMarshaler
	encoding.BinaryUnmarshaler
	TagType() TagType
}

var _ TagContent = (*TagContentRaw)(nil)

type TagContentRaw struct {
	Data []byte
}

func (t *TagContentRaw) TagType() TagType { return TagTypeRaw }

func (t *TagContentRaw) MarshalBinary() ([]byte, error) {
	return t.Data, nil
}

func (t *TagContentRaw) UnmarshalBinary(data []byte) error {
	t.Data = slices.Clone(data)
	return nil
}

var _ TagContent = (*TagContentCurve)(nil)

type TagContentCurve struct {
	Count uint32
	Data  []uint16
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
		Count:   t.Count,
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
	t.Count = curve.Count
	t.Data = make([]uint16, t.Count)
	if err := binary.Read(r, binary.BigEndian, &t.Data); err != nil {
		return err
	}
	return nil
}

func (t *tagContentCurve) Covert(x float64) float64 {
	return 0
}
