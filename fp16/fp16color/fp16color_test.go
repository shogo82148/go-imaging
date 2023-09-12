package fp16color

import (
	"math"
	"testing"

	"github.com/shogo82148/float16"
)

func TestRGBAh_RGBA(t *testing.T) {
	tests := []struct {
		fr, fg, fb, fa float64
		r, g, b, a     uint32
	}{
		{
			0.0, 0.0, 0.0, 0.0,
			0, 0, 0, 0,
		},
		{
			1.0, 1.0, 1.0, 1.0,
			0xffff, 0xffff, 0xffff, 0xffff,
		},
		{
			0.5, 0.5, 0.5, 0.5,
			0x7fff, 0x7fff, 0x7fff, 0x7fff,
		},
	}
	for _, tt := range tests {
		c := RGBAh{toFloat16(t, tt.fr), toFloat16(t, tt.fg), toFloat16(t, tt.fb), toFloat16(t, tt.fa)}
		r, g, b, a := c.RGBA()
		if r != tt.r || g != tt.g || b != tt.b || a != tt.a {
			t.Errorf("RGBAh(%v).RGBA() = %v, %v, %v, %v, want %v, %v, %v, %v", c, r, g, b, a, tt.r, tt.g, tt.b, tt.a)
		}
	}
}

func TestNRGBAh_RGBA(t *testing.T) {
	tests := []struct {
		fr, fg, fb, fa float64
		r, g, b, a     uint32
	}{
		{
			0.0, 0.0, 0.0, 0.0,
			0, 0, 0, 0,
		},
		{
			1.0, 1.0, 1.0, 1.0,
			0xffff, 0xffff, 0xffff, 0xffff,
		},
		{
			0.5, 0.5, 0.5, 1.0,
			0x7fff, 0x7fff, 0x7fff, 0xffff,
		},
	}
	for _, tt := range tests {
		c := NRGBAh{toFloat16(t, tt.fr), toFloat16(t, tt.fg), toFloat16(t, tt.fb), toFloat16(t, tt.fa)}
		r, g, b, a := c.RGBA()
		if r != tt.r || g != tt.g || b != tt.b || a != tt.a {
			t.Errorf("RGBAh(%v).RGBA() = %v, %v, %v, %v, want %v, %v, %v, %v", c, r, g, b, a, tt.r, tt.g, tt.b, tt.a)
		}
	}
}

func toFloat16(t *testing.T, f float64) float16.Float16 {
	t.Helper()
	if f < 0 || f > 1 {
		t.Errorf("%x is out of range", f)
	}
	if math.IsNaN(f) {
		t.Errorf("%x is NaN", f)
	}
	f16 := float16.FromFloat64(f)
	if f16.Float64() != f {
		t.Errorf("invalid test case: converting %x to float16 is lossy", f)
	}
	return f16
}
