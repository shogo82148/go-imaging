package icc

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func FuzzDecode(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		p0, err := Decode(data)
		if err != nil {
			return
		}

		encoded0, err := Encode(p0)
		if err != nil {
			t.Fatalf("failed to encode: %v", err)
		}

		p1, err := Decode(encoded0)
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

	})
}
