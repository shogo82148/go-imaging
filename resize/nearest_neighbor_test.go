package resize

import (
	"image"
	"testing"

	"github.com/shogo82148/go-imaging/internal/testutils"
)

func TestNearestNeighbor(t *testing.T) {
	dst := image.NewNRGBA64(image.Rect(0, 0, 4, 4))
	src := testutils.LoadPNM(`P3
2 2
65535
0     0 0  65535 65535 65535
32768 0 0      0 32768     0
`)
	NearestNeighbor(dst, src)

	want := testutils.LoadPNM(`P3
4 4
65535
0     0 0  0     0 0  65535 65535 65535  65535 65535 65535
0     0 0  0     0 0  65535 65535 65535  65535 65535 65535
32768 0 0  32768 0 0      0 32768     0      0 32768     0
32768 0 0  32768 0 0      0 32768     0      0 32768     0
`)
	testutils.Compare(t, dst, want)
}