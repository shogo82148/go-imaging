package icc

import (
	"bytes"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func FuzzDecode(f *testing.F) {
	testdata, err := os.ReadDir("testdata")
	if err != nil {
		f.Fatalf("failed to read testdata directory: %s", err)
	}
	for _, de := range testdata {
		if de.IsDir() || !strings.HasSuffix(de.Name(), ".icc") {
			continue
		}
		b, err := os.ReadFile(filepath.Join("testdata", de.Name()))
		if err != nil {
			f.Fatalf("failed to read testdata: %s", err)
		}
		f.Add(b)
	}
	f.Fuzz(func(t *testing.T, data []byte) {
		if len(data) > 1024*1024 {
			return
		}
		p0, err := Decode(bytes.NewReader(data))
		if err != nil {
			return
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

	})
}
