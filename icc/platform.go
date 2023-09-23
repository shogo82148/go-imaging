package icc

import "fmt"

type Platform uint32

const (
	PlatformApple     Platform = 0x4150504c // 'APPL'
	PlatformMicrosoft Platform = 0x4d534654 // 'MSFT'
	PlatformSGI       Platform = 0x53474920 // 'SGI '
	PlatformSun       Platform = 0x53554e57 // 'SUNW'
)

func (p Platform) String() string {
	switch p {
	case PlatformApple:
		return "Apple Computer, Inc."
	case PlatformMicrosoft:
		return "Microsoft Corporation"
	case PlatformSGI:
		return "Silicon Graphics, Inc."
	case PlatformSun:
		return "Sun Microsystems, Inc."
	default:
		return fmt.Sprintf(
			"Unknown Platform(%08xh '%c%c%c%c')",
			uint32(p),
			printable(byte(p>>24)),
			printable(byte(p>>16)),
			printable(byte(p>>8)),
			printable(byte(p)),
		)
	}
}
