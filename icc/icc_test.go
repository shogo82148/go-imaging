package icc

import (
	"bytes"
	"encoding/binary"
	"math"
	"os"
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
			{0.1, 0.003981127655368823},
			{0.5, 0.18946537237087296},
			{0.9, 0.7765730269567972},
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
			if got, want := trc.EncodeTone(tt.in), tt.out; !roughEqual(got, want) {
				t.Errorf("encode: got %v, want %v", got, want)
			}
			if got, want := trc.DecodeTone(tt.out), tt.in; !roughEqual(got, want) {
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

func TestTagContentParametricCurve(t *testing.T) {
	check := func(t *testing.T, curve Curve, min, max float64) {
		t.Helper()
		for i := 0; i <= 1000; i++ {
			x0 := float64(i) / 1000

			y0 := curve.EncodeTone(x0)
			if y0 < 0 || y0 > 1 || math.IsNaN(y0) {
				t.Errorf("encode(%f) is out of range [0, 1]: %f", x0, y0)
			}

			if min <= x0 && x0 <= max {
				x1 := curve.DecodeTone(y0)
				if got, want := x1, x0; !roughEqual(got, want) {
					t.Errorf("got %v, want %v", got, want)
				}
			}
		}

		for i := 0; i <= 1000; i++ {
			y0 := float64(i) / 1000
			x0 := curve.DecodeTone(y0)
			if x0 < 0 || x0 > 1 || math.IsNaN(x0) {
				t.Errorf("decode(%f) is out of range [0, 1]: %f", y0, x0)
			}
		}
	}

	t.Run("Gamma", func(t *testing.T) {
		// parameters come from https://github.com/LuaDist/lcms2/blob/6b9d6c0e27ab4e5922046ba7ce375b1297ff48f1/testbed/testcms2.c#L2826-L2830
		curve := new(TagContentParametricCurve)
		curve.FunctionType = 0
		curve.Params[0] = S15Fixed16NumberFromFloat64(2.2)
		check(t, curve, 0, 1)
	})

	t.Run("CIE122-1966", func(t *testing.T) {
		// parameters come from https://github.com/LuaDist/lcms2/blob/6b9d6c0e27ab4e5922046ba7ce375b1297ff48f1/testbed/testcms2.c#L2832-L2840
		curve := new(TagContentParametricCurve)
		curve.FunctionType = 1
		curve.Params[0] = S15Fixed16NumberFromFloat64(2.2)
		curve.Params[1] = S15Fixed16NumberFromFloat64(1.5)
		curve.Params[2] = S15Fixed16NumberFromFloat64(-0.5)
		check(t, curve, 1/3.0, 1)
	})

	t.Run("IEC 61966-3", func(t *testing.T) {
		// parameters come from https://github.com/LuaDist/lcms2/blob/6b9d6c0e27ab4e5922046ba7ce375b1297ff48f1/testbed/testcms2.c#L2842-L2852
		curve := new(TagContentParametricCurve)
		curve.FunctionType = 2
		curve.Params[0] = S15Fixed16NumberFromFloat64(2.2)
		curve.Params[1] = S15Fixed16NumberFromFloat64(1.5)
		curve.Params[2] = S15Fixed16NumberFromFloat64(-0.5)
		curve.Params[3] = S15Fixed16NumberFromFloat64(0.3)
		check(t, curve, 1/3.0, 0.9)
	})

	t.Run("IEC 61966-2.1 (sRGB)", func(t *testing.T) {
		// parameters come from https://github.com/LuaDist/lcms2/blob/6b9d6c0e27ab4e5922046ba7ce375b1297ff48f1/testbed/testcms2.c#L2854-L2864
		curve := new(TagContentParametricCurve)
		curve.FunctionType = 3
		curve.Params[0] = S15Fixed16NumberFromFloat64(2.4)
		curve.Params[1] = S15Fixed16NumberFromFloat64(1 / 1.055)
		curve.Params[2] = S15Fixed16NumberFromFloat64(0.055 / 1.055)
		curve.Params[3] = S15Fixed16NumberFromFloat64(1 / 12.92)
		curve.Params[4] = S15Fixed16NumberFromFloat64(0.04045)
		check(t, curve, 0, 1)
	})

	t.Run("param_5", func(t *testing.T) {
		// parameters come from https://github.com/LuaDist/lcms2/blob/6b9d6c0e27ab4e5922046ba7ce375b1297ff48f1/testbed/testcms2.c#L2867-L2878
		curve := new(TagContentParametricCurve)
		curve.FunctionType = 4
		curve.Params[0] = S15Fixed16NumberFromFloat64(2.2)
		curve.Params[1] = S15Fixed16NumberFromFloat64(0.7)
		curve.Params[2] = S15Fixed16NumberFromFloat64(0.2)
		curve.Params[3] = S15Fixed16NumberFromFloat64(0.3)
		curve.Params[4] = S15Fixed16NumberFromFloat64(0.1)
		curve.Params[5] = S15Fixed16NumberFromFloat64(0.5)
		curve.Params[6] = S15Fixed16NumberFromFloat64(0.2)
		check(t, curve, 0, 0.756)
	})
}

func TestDecode(t *testing.T) {
	t.Run("invalid magic", func(t *testing.T) {
		_, err := Decode(bytes.NewReader(make([]byte, 128)))
		if err == nil {
			t.Fatal("expected an error")
		}
	})

	t.Run("invalid size", func(t *testing.T) {
		data, err := os.ReadFile("testdata/iPhone12Pro.icc")
		if err != nil {
			t.Fatal(err)
		}

		// overwrite the profile size
		binary.BigEndian.PutUint32(data, uint32(len(data)+1))

		_, err = Decode(bytes.NewReader(data))
		if err == nil {
			t.Fatal("expected an error")
		}
	})

	t.Run("invalid tag offset", func(t *testing.T) {
		data, err := os.ReadFile("testdata/iPhone12Pro.icc")
		if err != nil {
			t.Fatal(err)
		}

		// overwrite the tag size
		binary.BigEndian.PutUint32(data[0x88:], uint32(len(data)-0x30+1))

		_, err = Decode(bytes.NewReader(data))
		if err == nil {
			t.Fatal("expected an error")
		}
	})

	t.Run("tag offset overflow", func(t *testing.T) {
		data, err := os.ReadFile("testdata/iPhone12Pro.icc")
		if err != nil {
			t.Fatal(err)
		}

		// overwrite the tag size
		binary.BigEndian.PutUint32(data[0x88:], uint32(0x100000000-0x30))

		_, err = Decode(bytes.NewReader(data))
		if err == nil {
			t.Fatal("expected an error")
		}
	})
}

func TestEncode(t *testing.T) {
	data, err := os.ReadFile("testdata/iPhone12Pro.icc")
	if err != nil {
		t.Fatal(err)
	}

	p0, err := Decode(bytes.NewReader(data))
	if err != nil {
		t.Fatal(err)
	}

	encoded0, err := Encode(p0)
	if err != nil {
		t.Fatalf("failed to encode: %v", err)
	}

	p1, err := Decode(bytes.NewReader(encoded0))
	if err != nil {
		t.Fatalf("failed to decode: %v", err)
	}

	// ignore differences in the profile size
	p0.Size = 0
	p1.Size = 0

	if diff := cmp.Diff(p0, p1); diff != "" {
		t.Errorf("mismatch (-want +got):\n%s", diff)
	}

	encoded1, err := Encode(p1)
	if err != nil {
		t.Fatalf("failed to encode: %v", err)
	}
	if diff := cmp.Diff(encoded0, encoded1); diff != "" {
		t.Errorf("mismatch (-want +got):\n%s", diff)
	}
}
