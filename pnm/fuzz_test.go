package pnm

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

var testFiles = []string{
	"feep-ascii.pbm",
	"feep-ascii.pgm",
	"feep-binary.pgm",
	"maze.pbm",
	"testgrid.pbm",
	"testimg.ppm",
	"wikipedia_example_j.pbm",
	"wikipedia_example_j2.pbm",
	"wikipedia_example_ppm1.ppm",
	"wikipedia_example_ppm2.ppm",
	"wikipedia_example_ppm3.ppm",
}

func FuzzDecode(f *testing.F) {
	for _, name := range testFiles {
		data, err := os.ReadFile(filepath.Join("testdata", name))
		if err != nil {
			panic(err)
		}
		f.Add(data)
	}
	f.Fuzz(func(t *testing.T, data []byte) {
		r := bytes.NewReader(data)
		Decode(r)
	})
}

func FuzzDecodeConfig(f *testing.F) {
	for _, name := range testFiles {
		data, err := os.ReadFile(filepath.Join("testdata", name))
		if err != nil {
			panic(err)
		}
		f.Add(data)
	}
	f.Fuzz(func(t *testing.T, data []byte) {
		r := bytes.NewReader(data)
		DecodeConfig(r)
	})
}
