package resize

import (
	"testing"

	"github.com/shogo82148/go-imaging/resize/internal/golden"
)

func TestLanczos2(t *testing.T) {
	src := golden.InputPattern()
	dst := golden.NewDst()
	Lanczos2(dst, src)
	golden.Assert(t, "lanczos2", dst)
}

func BenchmarkLanczos2(b *testing.B) {
	src := golden.InputPattern()
	dst := golden.NewDst()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Lanczos2(dst, src)
	}
}

func TestLanczos3(t *testing.T) {
	src := golden.InputPattern()
	dst := golden.NewDst()
	Lanczos3(dst, src)
	golden.Assert(t, "lanczos3", dst)
}

func BenchmarkLanczos3(b *testing.B) {
	src := golden.InputPattern()
	dst := golden.NewDst()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Lanczos3(dst, src)
	}
}

func TestLanczos4(t *testing.T) {
	src := golden.InputPattern()
	dst := golden.NewDst()
	Lanczos4(dst, src)
	golden.Assert(t, "lanczos4", dst)
}

func BenchmarkLanczos4(b *testing.B) {
	src := golden.InputPattern()
	dst := golden.NewDst()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Lanczos4(dst, src)
	}
}
