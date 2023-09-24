package icc

import (
	"bytes"
	"slices"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func FuzzDecode(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
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

	})
}
