// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package jpeg

import (
	"bytes"
	"image"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func FuzzDecode(f *testing.F) {
	if testing.Short() {
		f.Skip("Skipping in short mode")
	}

	testdata, err := os.ReadDir("../testdata")
	if err != nil {
		f.Fatalf("failed to read testdata directory: %s", err)
	}
	for _, de := range testdata {
		if de.IsDir() || !strings.HasSuffix(de.Name(), ".jpeg") {
			continue
		}
		b, err := os.ReadFile(filepath.Join("../testdata", de.Name()))
		if err != nil {
			f.Fatalf("failed to read testdata: %s", err)
		}
		f.Add(b)
	}

	f.Fuzz(func(t *testing.T, b []byte) {
		cfg, _, err := image.DecodeConfig(bytes.NewReader(b))
		if err != nil {
			return
		}
		if cfg.Width*cfg.Height > 1e6 {
			return
		}
		img, typ, err := image.Decode(bytes.NewReader(b))
		if err != nil || typ != "jpeg" {
			return
		}
		for q := 1; q <= 100; q++ {
			var w bytes.Buffer
			err := Encode(&w, img, &Options{Quality: q})
			if err != nil {
				t.Errorf("failed to encode valid image: %s", err)
				continue
			}
			img1, err := Decode(&w)
			if err != nil {
				t.Errorf("failed to decode roundtripped image: %s", err)
				continue
			}
			got := img1.Bounds()
			want := img.Bounds()
			if !got.Eq(want) {
				t.Errorf("roundtripped image bounds have changed, got: %s, want: %s", got, want)
			}
		}
	})
}

func FuzzDecodeWithMeta(f *testing.F) {
	if testing.Short() {
		f.Skip("Skipping in short mode")
	}

	testdata, err := os.ReadDir("../testdata")
	if err != nil {
		f.Fatalf("failed to read testdata directory: %s", err)
	}
	for _, de := range testdata {
		if de.IsDir() || !strings.HasSuffix(de.Name(), ".jpeg") {
			continue
		}
		b, err := os.ReadFile(filepath.Join("../testdata", de.Name()))
		if err != nil {
			f.Fatalf("failed to read testdata: %s", err)
		}
		f.Add(b)
	}

	f.Fuzz(func(t *testing.T, b []byte) {
		cfg, _, err := image.DecodeConfig(bytes.NewReader(b))
		if err != nil {
			return
		}
		if cfg.Width*cfg.Height > 1e6 {
			return
		}

		img0, err := DecodeWithMeta(bytes.NewReader(b))
		if err != nil {
			return
		}

		w := new(bytes.Buffer)
		err = EncodeWithMeta(w, img0, nil)
		if err != nil {
			t.Fatalf("failed to encode valid image: %s", err)
		}

		img1, err := DecodeWithMeta(w)
		if err != nil {
			t.Fatalf("failed to decode roundtripped image: %s", err)
		}

		if (img0.ICCProfile != nil) != (img1.ICCProfile != nil) {
			t.Errorf("ICC profile mismatch")
		}
	})
}
