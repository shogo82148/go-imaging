package pnm

import (
	"strings"
	"testing"
)

func TestDecode(t *testing.T) {
	// example from Wikipedia https://en.wikipedia.org/wiki/Netpbm#PBM_example
	data := `P1
# This is an example bitmap of the letter "J"
6 10
0 0 0 0 1 0
0 0 0 0 1 0
0 0 0 0 1 0
0 0 0 0 1 0
0 0 0 0 1 0
0 0 0 0 1 0
1 0 0 0 1 0
0 1 1 1 0 0
0 0 0 0 0 0
0 0 0 0 0 0
`
	r := strings.NewReader(data)
	img, err := Decode(r)
	if err != nil {
		t.Error(err)
	}
	img.At(0, 0)
}

func TestDecodeConfig(t *testing.T) {
	// example from Wikipedia https://en.wikipedia.org/wiki/Netpbm#PBM_example
	data := `P1
# This is an example bitmap of the letter "J"
6 10
0 0 0 0 1 0
0 0 0 0 1 0
0 0 0 0 1 0
0 0 0 0 1 0
0 0 0 0 1 0
0 0 0 0 1 0
1 0 0 0 1 0
0 1 1 1 0 0
0 0 0 0 0 0
0 0 0 0 0 0
`
	r := strings.NewReader(data)
	cfg, err := DecodeConfig(r)
	if err != nil {
		t.Error(err)
	}
	if cfg.Width != 6 {
		t.Errorf("expected width 6, got %d", cfg.Width)
	}
	if cfg.Height != 10 {
		t.Errorf("expected height 10, got %d", cfg.Height)
	}
}
