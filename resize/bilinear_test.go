package resize

import (
	"image"
	"testing"

	"github.com/shogo82148/go-imaging/internal/testutils"
)

func TestBiLinear(t *testing.T) {
	dst := image.NewNRGBA64(image.Rect(0, 0, 4, 4))
	src := testutils.LoadPNM(`P3
2 2
65535
0     0 0  65535 65535 65535
32768 0 0      0 32768     0
`)
	BiLinear(dst, src)

	want := testutils.LoadPNM(`P3
4 4
65535
0     0 0  21845 21845 21845  43690 43690 43690  65535 65535 65535
10923 0 0  21845 18204 14563  32768 36408 29127  43690 54613 43690
21845 0 0  21845 14563  7282  21845 29127 14563  21845 43690 21845
32768 0 0  21845 10923     0  10923 21845     0      0 32768     0
`)
	testutils.Compare(t, dst, want)
}

func BenchmarkBiLinear(b *testing.B) {
	dst := image.NewNRGBA64(image.Rect(0, 0, 512, 512))
	src := image.NewNRGBA64(image.Rect(0, 0, 256, 256))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		BiLinear(dst, src)
	}
}
