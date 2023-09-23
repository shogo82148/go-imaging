package icc

import "fmt"

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
