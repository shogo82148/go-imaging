package resize

import (
	"testing"

	"github.com/shogo82148/go-imaging/resize/internal/golden"
)

func TestCatmullRom(t *testing.T) {
	src := golden.InputPattern()
	dst := golden.NewDst()
	CatmullRom(dst, src)
	golden.Assert(t, "catmull_rom", dst)
}

func BenchmarkCatmullRom(b *testing.B) {
	src := golden.InputPattern()
	dst := golden.NewDst()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		CatmullRom(dst, src)
	}
}
