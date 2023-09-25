package icc

import (
	"bytes"
	"encoding/binary"
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"os"
	"slices"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func readPNG(filename string) (image.Image, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return png.Decode(f)
}

func decodeProfile(filename string) (*Profile, error) {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	return Decode(f)
}

func compareImage(t *testing.T, got, want image.Image) {
	if got.Bounds() != want.Bounds() {
		t.Errorf("bounds mismatch: got %v, want %v", got.Bounds(), want.Bounds())
		return
	}

	for y := got.Bounds().Min.Y; y < got.Bounds().Max.Y; y++ {
		for x := got.Bounds().Min.X; x < got.Bounds().Max.X; x++ {
			c0 := color.NRGBA64Model.Convert(got.At(x, y)).(color.NRGBA64)
			c1 := color.NRGBA64Model.Convert(want.At(x, y)).(color.NRGBA64)
			if c0 != c1 {
				t.Errorf("color mismatch at (%d, %d): got %v, want %v", x, y, c0, c1)
				return
			}
		}
	}
}

func TestProfile_DecodeTone(t *testing.T) {
	input, err := readPNG("../testdata/senkakuwan.png")
	if err != nil {
		t.Fatal(err)
	}

	want, err := readPNG("./testdata/senkakuwan.golden.png")
	if err != nil {
		t.Fatal(err)
	}

	profile, err := decodeProfile("testdata/iPhone12Pro.icc")
	if err != nil {
		t.Fatal(err)
	}
	got := profile.DecodeTone(input)

	compareImage(t, got, want)
}

func TestTagContentParametricCurve(t *testing.T) {
	check := func(t *testing.T, curve Curve, min, max float64) {
		t.Helper()
		for i := 0; i <= 1000; i++ {
			x0 := float64(i) / 1000

			y0 := curve.DecodeTone(x0)
			if y0 < 0 || y0 > 1 || math.IsNaN(y0) {
				t.Errorf("encode(%f) is out of range [0, 1]: %f", x0, y0)
			}

			if min <= x0 && x0 <= max {
				x1 := curve.EncodeTone(y0)
				if got, want := x1, x0; !roughEqual(got, want) {
					t.Errorf("got %v, want %v", got, want)
				}
			}
		}

		for i := 0; i <= 1000; i++ {
			y0 := float64(i) / 1000
			x0 := curve.EncodeTone(y0)
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
		binary.BigEndian.PutUint32(data, uint32(iccHeaderSize-1))

		_, err = Decode(bytes.NewReader(data))
		if err == nil {
			t.Fatal("expected an error")
		}
	})

	t.Run("large size", func(t *testing.T) {
		data, err := os.ReadFile("testdata/iPhone12Pro.icc")
		if err != nil {
			t.Fatal(err)
		}

		// overwrite the profile size
		binary.BigEndian.PutUint32(data, uint32(0xffffffff))

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

	buf := new(bytes.Buffer)
	if err := p0.Encode(buf); err != nil {
		t.Fatalf("failed to encode: %v", err)
	}
	encoded0 := slices.Clone(buf.Bytes())

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

	buf.Reset()
	if err := p1.Encode(buf); err != nil {
		t.Fatalf("failed to encode: %v", err)
	}
	encoded1 := slices.Clone(buf.Bytes())

	if diff := cmp.Diff(encoded0, encoded1); diff != "" {
		t.Errorf("mismatch (-want +got):\n%s", diff)
	}
}
