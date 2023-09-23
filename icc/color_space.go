package icc

import "fmt"

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
			"Unknown Color Space(%08xh '%c%c%c%c')",
			uint32(cs),
			printable(byte(cs>>24)),
			printable(byte(cs>>16)),
			printable(byte(cs>>8)),
			printable(byte(cs)),
		)
	}
}
