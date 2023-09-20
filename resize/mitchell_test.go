package resize

import (
	"testing"

	"github.com/shogo82148/go-imaging/resize/internal/golden"
)

func TestMitchell(t *testing.T) {
	src := golden.InputPattern()
	dst := golden.NewDst()
	Mitchell(dst, src)
	golden.Assert(t, "mitchell", dst)
}

func BenchmarkMitchell(b *testing.B) {
	src := golden.InputPattern()
	dst := golden.NewDst()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Mitchell(dst, src)
	}
}
