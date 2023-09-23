package icc

import "fmt"

type Platform Signature

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
			"Unknown Platform(%08xh '%s')",
			uint32(p),
			Signature(p),
		)
	}
}
