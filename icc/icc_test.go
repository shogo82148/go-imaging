package icc

import (
	"bytes"
	"math"
	"os"
	"slices"
	"testing"

	"github.com/google/go-cmp/cmp"
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
			{0.1, 0.010023912650605525},
			{0.5, 0.2140451926881055},
			{0.9, 0.7874141418532061},
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

	// profile id
	profileID := [16]byte{0xec, 0xfd, 0xa3, 0x8e, 0x38, 0x85, 0x47, 0xc3, 0x6d, 0xb4, 0xbd, 0x4f, 0x7a, 0xda, 0x18, 0x2f}
	if profile.ProfileID != profileID {
		t.Errorf("unexpected profile id: want %016x, got %016x", profileID, profile.ProfileID)
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

func TestEncode(t *testing.T) {
	t.Run("iPhone12Pro", func(t *testing.T) {
		data, err := os.ReadFile("testdata/iPhone12Pro.icc")
		if err != nil {
			t.Fatal(err)
		}

		p0, err := Decode(bytes.NewReader(data))
		if err != nil {
			t.Fatal(err)
		}

		buf := new(bytes.Buffer)
		if err := p0.Encode(buf); err != nil {
			t.Fatalf("failed to encode: %v", err)
		}
		encoded0 := slices.Clone(buf.Bytes())

		p1, err := Decode(bytes.NewReader(encoded0))
		if err != nil {
			t.Fatalf("failed to decode: %v", err)
		}

		if diff := cmp.Diff(p0, p1); diff != "" {
			t.Errorf("mismatch (-want +got):\n%s", diff)
		}

		buf.Reset()
		if err := p1.Encode(buf); err != nil {
			t.Fatalf("failed to encode: %v", err)
		}
		encoded1 := slices.Clone(buf.Bytes())

		if diff := cmp.Diff(data, encoded0); diff != "" {
			t.Errorf("mismatch (-want +got):\n%s", diff)
		}

		if diff := cmp.Diff(encoded0, encoded1); diff != "" {
			t.Errorf("mismatch (-want +got):\n%s", diff)
		}

		if len(encoded0)%4 != 0 {
			t.Errorf("encoded0 is not a multiple of 4 bytes")
		}
	})

	t.Run("gimp-linear", func(t *testing.T) {
		data, err := os.ReadFile("testdata/gimp-linear.icc")
		if err != nil {
			t.Fatal(err)
		}

		p0, err := Decode(bytes.NewReader(data))
		if err != nil {
			t.Fatal(err)
		}

		buf := new(bytes.Buffer)
		if err := p0.Encode(buf); err != nil {
			t.Fatalf("failed to encode: %v", err)
		}
		encoded0 := slices.Clone(buf.Bytes())

		p1, err := Decode(bytes.NewReader(encoded0))
		if err != nil {
			t.Fatalf("failed to decode: %v", err)
		}

		// ignore differences in the profile size and the profile id
		p0.Size = 0
		p1.Size = 0
		clear(p0.ProfileID[:])
		clear(p1.ProfileID[:])

		if diff := cmp.Diff(p0, p1); diff != "" {
			t.Errorf("mismatch (-want +got):\n%s", diff)
		}

		buf.Reset()
		if err := p1.Encode(buf); err != nil {
			t.Fatalf("failed to encode: %v", err)
		}
		encoded1 := slices.Clone(buf.Bytes())

		if diff := cmp.Diff(encoded0, encoded1); diff != "" {
			t.Errorf("mismatch (-want +got):\n%s", diff)
		}

		if len(encoded0)%4 != 0 {
			t.Errorf("encoded0 is not a multiple of 4 bytes")
		}
	})
}
