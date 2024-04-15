package resize

import (
	"testing"

	"github.com/shogo82148/go-imaging/resize/internal/golden"
)

func TestAreaAverage(t *testing.T) {
	src := golden.InputPattern()
	dst := golden.NewDst()
	AreaAverage(dst, src)
	golden.Assert(t, "area_average", dst)
}

func BenchmarkAreaAverage(b *testing.B) {
	src := golden.InputPattern()
	dst := golden.NewDst()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		AreaAverage(dst, src)
	}
}
