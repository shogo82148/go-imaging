package icc

import (
	"math"
	"os"
	"testing"
)

func roughEqual(a, b float64) bool {
	d := math.Abs(a - b)
	return d < 1.0/0xffff
}

func Test_D65_XYZ(t *testing.T) {
	data, err := os.ReadFile("testdata/D65_XYZ.icc")
	if err != nil {
		t.Fatal(err)
	}

	profile, err := Decode(data)
	if err != nil {
		t.Fatal(err)
	}

	redTRC, ok := profile.Get(TagRedTRC).(Curve)
	if !ok {
		t.Fatal("TagRedTRC is not a Curve")
	}
	tests := []struct {
		in, out float64
	}{
		{0.0, 0.0},
		{0.1, 0.1},
		{0.5, 0.5},
		{0.9, 0.9},
		{1.0, 1.0},
	}
	for _, tt := range tests {
		if got, want := redTRC.EncodeTone(tt.in), tt.out; !roughEqual(got, want) {
			t.Errorf("encode: got %v, want %v", got, want)
		}
		if got, want := redTRC.DecodeTone(tt.out), tt.in; !roughEqual(got, want) {
			t.Errorf("decode: got %v, want %v", got, want)
		}
	}
}

func Test_sRGB_IEC61966(t *testing.T) {
	data, err := os.ReadFile("testdata/sRGB_IEC61966-2-1_black_scaled.icc")
	if err != nil {
		t.Fatal(err)
	}

	profile, err := Decode(data)
	if err != nil {
		t.Fatal(err)
	}

	redTRC, ok := profile.Get(TagRedTRC).(Curve)
	if !ok {
		t.Fatal("TagRedTRC is not a Curve")
	}
	tests := []struct {
		in, out float64
	}{
		{0.0, 0.0},
		{0.1, 0.01002517738612955},
		{0.5, 0.21404592965590905},
		{0.9, 0.7874097810330358},
		{1.0, 1.0},
	}
	for _, tt := range tests {
		if got, want := redTRC.EncodeTone(tt.in), tt.out; !roughEqual(got, want) {
			t.Errorf("encode: got %v, want %v", got, want)
		}
		if got, want := redTRC.DecodeTone(tt.out), tt.in; !roughEqual(got, want) {
			t.Errorf("decode: got %v, want %v", got, want)
		}
	}
}
