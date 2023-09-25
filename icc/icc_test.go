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

func TestVersion(t *testing.T) {
	if got, want := Version(0x05000000).String(), "5.0.0.0"; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	if got, want := Version(0x04400000).String(), "4.4.0.0"; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	if got, want := Version(0x02400000).String(), "2.4.0.0"; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}

func Test_D65_XYZ(t *testing.T) {
	f, err := os.Open("testdata/D65_XYZ.icc")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	profile, err := Decode(f)
	if err != nil {
		t.Fatal(err)
	}

	// check the profile class
	if got, want := profile.Class, ClassDisplay; got != want {
		t.Errorf("got %v, want %v", got, want)
	}

	// check the color space
	if got, want := profile.ColorSpace, ColorSpaceRGB; got != want {
		t.Errorf("got %v, want %v", got, want)
	}

	for _, tag := range []Tag{TagRedTRC, TagGreenTRC, TagBlueTRC} {
		trc, ok := profile.Get(tag).(Curve)
		if !ok {
			t.Fatal("tag is not a tone reproduction curve")
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
			if got, want := trc.EncodeTone(tt.in), tt.out; !roughEqual(got, want) {
				t.Errorf("encode: got %v, want %v", got, want)
			}
			if got, want := trc.DecodeTone(tt.out), tt.in; !roughEqual(got, want) {
				t.Errorf("decode: got %v, want %v", got, want)
			}
		}
	}
}

func Test_ILFORD_CANpro(t *testing.T) {
	f, err := os.Open("testdata/ILFORD_CANpro-4000_GPGFG_ProPlatin.icc")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	profile, err := Decode(f)
	if err != nil {
		t.Fatal(err)
	}

	// check the profile class
	if got, want := profile.Class, ClassOutput; got != want {
		t.Errorf("got %v, want %v", got, want)
	}

	// check the color space
	if got, want := profile.ColorSpace, ColorSpaceRGB; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}

func Test_iPhone12Pro(t *testing.T) {
	f, err := os.Open("testdata/iPhone12Pro.icc")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	profile, err := Decode(f)
	if err != nil {
		t.Fatal(err)
	}

	// check the profile class
	if got, want := profile.Class, ClassDisplay; got != want {
		t.Errorf("got %v, want %v", got, want)
	}

	// check the color space
	if got, want := profile.ColorSpace, ColorSpaceRGB; got != want {
		t.Errorf("got %v, want %v", got, want)
	}

	for _, tag := range []Tag{TagRedTRC, TagGreenTRC, TagBlueTRC} {
		trc, ok := profile.Get(tag).(Curve)
		if !ok {
			t.Fatal("tag is not a tone reproduction curve")
		}
		tests := []struct {
			in, out float64
		}{
			{0.0, 0.0},
			{0.1, 0.003981127655368823},
			{0.5, 0.18946537237087296},
			{0.9, 0.7765730269567972},
			{1.0, 1.0},
		}
		for _, tt := range tests {
			if got, want := trc.DecodeTone(tt.in), tt.out; !roughEqual(got, want) {
				t.Errorf("encode: got %v, want %v", got, want)
			}
			if got, want := trc.EncodeTone(tt.out), tt.in; !roughEqual(got, want) {
				t.Errorf("decode: got %v, want %v", got, want)
			}
		}
	}
}

func Test_Probev_ICCv2(t *testing.T) {
	f, err := os.Open("testdata/Probev1_ICCv2.icc")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	profile, err := Decode(f)
	if err != nil {
		t.Fatal(err)
	}

	// check the profile class
	if got, want := profile.Class, ClassOutput; got != want {
		t.Errorf("got %v, want %v", got, want)
	}

	// check the color space
	if got, want := profile.ColorSpace, ColorSpaceCMYK; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}

func Test_Probev_ICCv4(t *testing.T) {
	f, err := os.Open("testdata/Probev1_ICCv4.icc")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	profile, err := Decode(f)
	if err != nil {
		t.Fatal(err)
	}

	// check the profile class
	if got, want := profile.Class, ClassOutput; got != want {
		t.Errorf("got %v, want %v", got, want)
	}

	// check the color space
	if got, want := profile.ColorSpace, ColorSpaceCMYK; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}

func Test_sRGB_ICC_v4_Appearance(t *testing.T) {
	f, err := os.Open("testdata/sRGB_ICC_v4_Appearance.icc")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	profile, err := Decode(f)
	if err != nil {
		t.Fatal(err)
	}

	// check the profile class
	if got, want := profile.Class, ClassDisplay; got != want {
		t.Errorf("got %v, want %v", got, want)
	}

	// check the color space
	if got, want := profile.ColorSpace, ColorSpaceRGB; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}

func Test_sRGB_IEC61966(t *testing.T) {
	f, err := os.Open("testdata/sRGB_IEC61966-2-1_black_scaled.icc")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	profile, err := Decode(f)
	if err != nil {
		t.Fatal(err)
	}

	// check the profile class
	if got, want := profile.Class, ClassDisplay; got != want {
		t.Errorf("got %v, want %v", got, want)
	}

	// check the color space
	if got, want := profile.ColorSpace, ColorSpaceRGB; got != want {
		t.Errorf("got %v, want %v", got, want)
	}

	for _, tag := range []Tag{TagRedTRC, TagGreenTRC, TagBlueTRC} {
		trc, ok := profile.Get(tag).(Curve)
		if !ok {
			t.Fatal("tag is not a tone reproduction curve")
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
			if got, want := trc.DecodeTone(tt.in), tt.out; !roughEqual(got, want) {
				t.Errorf("encode: got %v, want %v", got, want)
			}
			if got, want := trc.EncodeTone(tt.out), tt.in; !roughEqual(got, want) {
				t.Errorf("decode: got %v, want %v", got, want)
			}
		}
	}
}

func Test_USWebCoatedSWOP(t *testing.T) {
	f, err := os.Open("testdata/USWebCoatedSWOP.icc")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	profile, err := Decode(f)
	if err != nil {
		t.Fatal(err)
	}

	// check the profile class
	if got, want := profile.Class, ClassOutput; got != want {
		t.Errorf("got %v, want %v", got, want)
	}

	// check the color space
	if got, want := profile.ColorSpace, ColorSpaceCMYK; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}
