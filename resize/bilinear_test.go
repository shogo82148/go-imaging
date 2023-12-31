package resize

import (
	"testing"

	"github.com/shogo82148/go-imaging/resize/internal/golden"
)

func TestBiLinear(t *testing.T) {
	src := golden.InputPattern()
	dst := golden.NewDst()
	BiLinear(dst, src)
	golden.Assert(t, "bilinear", dst)
}

func BenchmarkBiLinear(b *testing.B) {
	src := golden.InputPattern()
	dst := golden.NewDst()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		BiLinear(dst, src)
	}
}
