package icc

import (
	"bytes"
	"encoding"
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"slices"
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
			"Unknown(%08xh '%c%c%c%c')",
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

type Tag uint32

// ICC.2 tags
// https://www.color.org/specification/ICC.2-2019.pdf
const (
	TagAToB0                          Tag = 0x41324230 // 'A2B0'
	TagAToB1                          Tag = 0x41324231 // 'A2B1'
	TagAToB2                          Tag = 0x41324232 // 'A2B2'
	TagAToB3                          Tag = 0x41324233 // 'A2B3'
	TagAToM0                          Tag = 0x41324d30 // 'A2M0'
	TagBRDFColormetricParameter0      Tag = 0x62637030 // 'bcp0'
	TagBRDFColormetricParameter1      Tag = 0x62637031 // 'bcp1'
	TagBRDFColormetricParameter2      Tag = 0x62637032 // 'bcp2'
	TagBRDFColormetricParameter3      Tag = 0x62637033 // 'bcp3'
	TagBRDFSpectralParameter0         Tag = 0x62737030 // 'bsp0'
	TagBRDFSpectralParameter1         Tag = 0x62737031 // 'bsp1'
	TagBRDFSpectralParameter2         Tag = 0x62737032 // 'bsp2'
	TagBRDFSpectralParameter3         Tag = 0x62737033 // 'bsp3'
	TagBRDFAToB0                      Tag = 0x62414230 // 'bAB0'
	TagBRDFAToB1                      Tag = 0x62414231 // 'bAB1'
	TagBRDFAToB2                      Tag = 0x62414232 // 'bAB2'
	TagBRDFAToB3                      Tag = 0x62414233 // 'bAB3'
	TagBRDFBToA0                      Tag = 0x62424130 // 'bBA0'
	TagBRDFBToA1                      Tag = 0x62424131 // 'bBA1'
	TagBRDFBToA2                      Tag = 0x62424132 // 'bBA2'
	TagBRDFBToA3                      Tag = 0x62424133 // 'bBA3'
	TagBRDFBToD0                      Tag = 0x62424430 // 'bBD0'
	TagBRDFBToD1                      Tag = 0x62424431 // 'bBD1'
	TagBRDFBToD2                      Tag = 0x62424432 // 'bBD2'
	TagBRDFBToD3                      Tag = 0x62424433 // 'bBD3'
	TagBRDFDToB0                      Tag = 0x62444230 // 'bDB0'
	TagBRDFDToB1                      Tag = 0x62444231 // 'bDB1'
	TagBRDFDToB2                      Tag = 0x62444232 // 'bDB2'
	TagBRDFDToB3                      Tag = 0x62444233 // 'bDB3'
	TagBRDFMToB0                      Tag = 0x624d4230 // 'bMB0'
	TagBRDFMToB1                      Tag = 0x624d4231 // 'bMB1'
	TagBRDFMToB2                      Tag = 0x624d4232 // 'bMB2'
	TagBRDFMToB3                      Tag = 0x624d4233 // 'bMB3'
	TagBRDFMToS0                      Tag = 0x624d5330 // 'bMS0'
	TagBRDFMToS1                      Tag = 0x624d5331 // 'bMS1'
	TagBRDFMToS2                      Tag = 0x624d5332 // 'bMS2'
	TagBRDFMToS3                      Tag = 0x624d5333 // 'bMS3'
	TagBToA0                          Tag = 0x42324130 // 'B2A0'
	TagBToA1                          Tag = 0x42324131 // 'B2A1'
	TagBToA2                          Tag = 0x42324132 // 'B2A2'
	TagBToA3                          Tag = 0x42324133 // 'B2A3'
	TagBToD0                          Tag = 0x42324430 // 'B2D0'
	TagBToD1                          Tag = 0x42324431 // 'B2D1'
	TagBToD2                          Tag = 0x42324432 // 'B2D2'
	TagBToD3                          Tag = 0x42324433 // 'B2D3'
	TagCalibrationDateTime            Tag = 0x63616c74 // 'calt'
	TagCharTarget                     Tag = 0x74617267 // 'targ'
	TagColorEncodingParams            Tag = 0x63657074 // 'cept'
	TagColorSpaceName                 Tag = 0x63736e6d // 'csnm'
	TagColorantOrder                  Tag = 0x636c726f // 'clro'
	TagColorantOrderOut               Tag = 0x636c6f6f // 'cloo'
	TagColorantInfo                   Tag = 0x636c696e // 'clin'
	TagColorantInfoOut                Tag = 0x636c696f // 'clio'
	TagColorimetricIntentImageState   Tag = 0x63696973 // 'ciis'
	TagCopyright                      Tag = 0x63707274 // 'cprt'
	TagCustomToStandardPcc            Tag = 0x63327370 // 'c2sp'
	TagCXF                            Tag = 0x43784620 // 'CxF '
	TagDeviceMfgDesc                  Tag = 0x646d6e64 // 'dmnd'
	TagDeviceModelDesc                Tag = 0x646d6464 // 'dmdd'
	TagDirectionalAToB0               Tag = 0x64414230 // 'dAB0'
	TagDirectionalAToB1               Tag = 0x64414231 // 'dAB1'
	TagDirectionalAToB2               Tag = 0x64414232 // 'dAB2'
	TagDirectionalAToB3               Tag = 0x64414233 // 'dAB3'
	TagDirectionalBToA0               Tag = 0x64424130 // 'dBA0'
	TagDirectionalBToA1               Tag = 0x64424131 // 'dBA1'
	TagDirectionalBToA2               Tag = 0x64424132 // 'dBA2'
	TagDirectionalBToA3               Tag = 0x64424133 // 'dBA3'
	TagDirectionalBToD0               Tag = 0x64424430 // 'dBD0'
	TagDirectionalBToD1               Tag = 0x64424431 // 'dBD1'
	TagDirectionalBToD2               Tag = 0x64424432 // 'dBD2'
	TagDirectionalBToD3               Tag = 0x64424433 // 'dBD3'
	TagDirectionalDToB0               Tag = 0x64444230 // 'dDB0'
	TagDirectionalDToB1               Tag = 0x64444231 // 'dDB1'
	TagDirectionalDToB2               Tag = 0x64444232 // 'dDB2'
	TagDirectionalDToB3               Tag = 0x64444233 // 'dDB3'
	TagDToB0                          Tag = 0x44324230 // 'D2B0'
	TagDToB1                          Tag = 0x44324231 // 'D2B1'
	TagDToB2                          Tag = 0x44324232 // 'D2B2'
	TagDToB3                          Tag = 0x44324233 // 'D2B3'
	TagGamutBoundaryDescription0      Tag = 0x67626430 // 'gbd0'
	TagGamutBoundaryDescription1      Tag = 0x67626431 // 'gbd1'
	TagGamutBoundaryDescription2      Tag = 0x67626432 // 'gbd2'
	TagGamutBoundaryDescription3      Tag = 0x67626433 // 'gbd3'
	TagMultiplexDefaultValues         Tag = 0x6d647620 // 'mdv '
	TagMultiplexTypeArray             Tag = 0x6d637461 // 'mcta'
	TagMeasurementInfo                Tag = 0x6d696e66 // 'minf'
	TagMeasurementInputInfo           Tag = 0x6d69696e // 'miin'
	TagMediaWhitePoint                Tag = 0x77747074 // 'wtpt'
	TagMetadata                       Tag = 0x6d657461 // 'meta'
	TagMToA0                          Tag = 0x4d546130 // 'M2A0'
	TagMToB0                          Tag = 0x4d546230 // 'M2B0'
	TagMToB1                          Tag = 0x4d324231 // 'M2B1'
	TagMToB2                          Tag = 0x4d324232 // 'M2B2'
	TagMToB3                          Tag = 0x4d324233 // 'M2B3'
	TagMToS0                          Tag = 0x4d546130 // 'M2S0'
	TagMToS1                          Tag = 0x4d546131 // 'M2S1'
	TagMToS2                          Tag = 0x4d546132 // 'M2S2'
	TagMToS3                          Tag = 0x4d546133 // 'M2S3'
	TagNamedColor                     Tag = 0x6e6d636c // 'nmcl'
	TagPerceptualRenderingIntentGamut Tag = 0x72696730 // 'rig0'
	TagProfileDescription             Tag = 0x64657363 // 'desc'
	TagProfileSequenceInformation     Tag = 0x7073696e // 'psin'
	TagReferenceName                  Tag = 0x72666e6d // 'rfnm'
	TagSaturationRenderingIntentGamut Tag = 0x72696732 // 'rig2'
	TagSpectralViewingConditions      Tag = 0x7376636e // 'svcn'
	TagSpectralWhitePoint             Tag = 0x73777074 // 'swpt'
	TagStandardToCustomPcc            Tag = 0x73326370 // 's2cp'
	TagSurfaceMap                     Tag = 0x736d6170 // 'smap'
	TagTechnology                     Tag = 0x74656368 // 'tech'
)

// ICC.2 colorEncodingParamsStructure element sub-tags
const (
	TagCeptBluePrimaryXYZMbr                Tag = 0x6258595a // 'bXYZ'
	TagCeptGreenPrimaryXYZMbr               Tag = 0x6758595a // 'gXYZ'
	TagCeptRedPrimaryXYZMbr                 Tag = 0x7258595a // 'rXYZ'
	TagCeptTransferFunctionMbr              Tag = 0x66756e63 // 'func'
	TagCeptLumaChromaMatrixMbr              Tag = 0x6c6d6174 // 'lmat'
	TagCeptWhitePointLuminanceMbr           Tag = 0x776c756d // 'wlum'
	TagCeptWhitePointChromaticityMbr        Tag = 0x7758595a // 'wXYZ'
	TagCeptEncodingRangeMbr                 Tag = 0x65526e67 // 'eRng'
	TagCeptBitDepthMbr                      Tag = 0x62697473 // 'bits'
	TagCeptImageStateMbr                    Tag = 0x696d7374 // 'imst'
	TagCeptImageBackgroundMbr               Tag = 0x69626b67 // 'ibkg'
	TagCeptViewingSurroundMbr               Tag = 0x73726e64 // 'srnd'
	TagCeptAmbientIlluminanceMbr            Tag = 0x61696c6d // 'ailm'
	TagCeptAmbientWhitePointLuminanceMbr    Tag = 0x61776c6d // 'awlm'
	TagCeptAmbientWhitePointChromaticityMbr Tag = 0x61777063 // 'awpc'
	TagCeptViewingFlareMbr                  Tag = 0x666c6172 // 'flar'
	TagCeptValidRelativeLuminanceRangeMbr   Tag = 0x6c726e67 // 'lrng'
	TagCeptMediumWhitePointLuminanceMbr     Tag = 0x6d77706c // 'mwpl'
	TagCeptMediumWhitePointChromaticityMbr  Tag = 0x6d777063 // 'mwpc'
	TagCeptMediumBlackPointLuminanceMbr     Tag = 0x6d62706c // 'mbpl'
	TagCeptMediumBlackPointChromaticityMbr  Tag = 0x6d627063 // 'mbpc'
)

// ICC.1 tags
// https://www.color.org/specification/ICC.1-2022-05.pdf
const (
	// TagAToB0               Tag = 0x41324230 // 'A2B0'
	// TagAToB1               Tag = 0x41324231 // 'A2B1'
	// TagAToB2               Tag = 0x41324232 // 'A2B2'
	TagBlueMatrixColumn Tag = 0x6258595a // 'bXYZ'
	TagBlueTRC          Tag = 0x62545243 // 'bTRC'
	// TagBToA0            Tag = 0x42324130 // 'B2A0'
	// TagBToA1               Tag = 0x42324131 // 'B2A1'
	// TagBToA2               Tag = 0x42324132 // 'B2A2'
	// TagBToD0               Tag = 0x42324430 // 'B2D0'
	// TagBToD1               Tag = 0x42324431 // 'B2D1'
	// TagBToD2               Tag = 0x42324432 // 'B2D2'
	// TagBToD3               Tag = 0x42324433 // 'B2D3'
	// TagCalibrationDateTime Tag = 0x63616c74 // 'calt'
	// TagCharTarget          Tag = 0x74617267 // 'targ'
	TagChromaticAdaptation Tag = 0x63686164 // 'chad'
	TagCICP                Tag = 0x63696370 // 'cicp'
	// TagColorantOrder                  Tag = 0x636c726f // 'clro'
	TagColorantTable    Tag = 0x636c7274 // 'clrt'
	TagColorantTableOut Tag = 0x636c6f74 // 'clot'
	// TagColorimetricIntentImageState   Tag = 0x63696973 // 'ciis'
	// TagCopyright                      Tag = 0x63707274 // 'cprt'
	// TagDeviceMfgDesc                  Tag = 0x646d6e64 // 'dmnd'
	// TagDeviceModelDesc                Tag = 0x646d6464 // 'dmdd'
	// TagDToB0                          Tag = 0x44324230 // 'D2B0'
	// TagDToB1                          Tag = 0x44324231 // 'D2B1'
	// TagDToB2                          Tag = 0x44324232 // 'D2B2'
	// TagDToB3                          Tag = 0x44324233 // 'D2B3'
	TagGamut             Tag = 0x67616d74 // 'gamt'
	TagGrayTRC           Tag = 0x6b545243 // 'kTRC'
	TagGreenMatrixColumn Tag = 0x6758595a // 'gXYZ'
	TagGreenTRC          Tag = 0x67545243 // 'gTRC'
	TagLuminance         Tag = 0x6c756d69 // 'lumi'
	TagMeasurement       Tag = 0x6d656173 // 'meas'
	// TagMetadata                       Tag = 0x6d657461 // 'meta'
	// TagMediaWhitePoint                Tag = 0x77747074 // 'wtpt'
	TagNamedColor2    Tag = 0x6e636c32 // 'ncl2'
	TagOutputResponse Tag = 0x72657370 // 'resp'
	// TagPerceptualRenderingIntentGamut Tag = 0x72696730 // 'rig0'
	TagPreview0 Tag = 0x70726530 // 'pre0'
	TagPreview1 Tag = 0x70726531 // 'pre1'
	TagPreview2 Tag = 0x70726532 // 'pre2'
	// TagProfileDescription             Tag = 0x64657363 // 'desc'
	TagProfileSequenceDesc       Tag = 0x70736571 // 'pseq'
	TagProfileSequenceIdentifier Tag = 0x70736964 // 'psid'
	TagRedMatrixColumn           Tag = 0x7258595a // 'rXYZ'
	TagRedTRC                    Tag = 0x72545243 // 'rTRC'
	// TagSaturationRenderingIntentGamut Tag = 0x72696732 // 'rig2'
	// TagTechnology                     Tag = 0x74656368 // 'tech'
	TagViewingCondDesc   Tag = 0x76756564 // 'vued'
	TagViewingConditions Tag = 0x76696577 // 'view'
)

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

func (t *TagContentCurve) EncodeTone(x float64) float64 {
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
