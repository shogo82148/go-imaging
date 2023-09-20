package resize

import (
	"testing"

	"github.com/shogo82148/go-imaging/resize/internal/golden"
)

func TestGeneral(t *testing.T) {
	src := golden.InputPattern()
	dst := golden.NewDst()
	General(dst, src)
	golden.Assert(t, "general", dst)
}

func BenchmarkGeneral(b *testing.B) {
	src := golden.InputPattern()
	dst := golden.NewDst()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		General(dst, src)
	}
}
