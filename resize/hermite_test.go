package resize

import (
	"testing"

	"github.com/shogo82148/go-imaging/resize/internal/golden"
)

func TestHermite(t *testing.T) {
	src := golden.InputPattern()
	dst := golden.NewDst()
	Hermite(dst, src)
	golden.Assert(t, "hermite", dst)
}

func BenchmarkHermite(b *testing.B) {
	src := golden.InputPattern()
	dst := golden.NewDst()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Hermite(dst, src)
	}
}
